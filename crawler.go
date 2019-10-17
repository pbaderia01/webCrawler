package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var insertCounter int64
var visitedCounter int64
var inserted sync.Map

func usage() {
	fmt.Fprintf(os.Stderr, "Please Use the module as : go run <filename>.go <URI>\n")
	flag.PrintDefaults()
}

func main() {
	var args []string
	if os.Getenv("CRAWL_URL")==""{
		flag.Usage = usage
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			usage()
			fmt.Println("Please specify start page")
			os.Exit(1)
		}} else {
			args = append(args,os.Getenv("CRAWL_URL"))
		}
	args[0] = checkValidBaseURL(args[0])
	hostURL, err := url.Parse(args[0])
	if err!=nil{
		fmt.Println("Invalid URL provided. Please find the error below:")
		fmt.Println(err)
		os.Exit(1)
	}
	writeOnDisk := checkWriteOnDisk()
	inserted.Store(args[0],true)
	inserted.Store(args[0]+"/",true)
	hostBaseURL := hostURL.Hostname()
	queue := make(chan string)
	go func() { queue <- args[0]}()
	atomic.AddInt64(&insertCounter,1)
	done := make(chan bool)
	rootPath := os.Getenv("ROOT_PATH")
	if writeOnDisk && !checkValidDiskPath(rootPath){
		os.Exit(1)
	}
	uriOutputflag := os.Getenv("DISPLAY_URI")
	uriOutput,_ := strconv.ParseBool(uriOutputflag)
	createConcurrentThreads(done,queue,hostBaseURL, writeOnDisk, rootPath, uriOutput)
}

/* The function takes the following arguments and starts execution with a defa
Arguments:
  	Done : Bool channel to synchronize thread execution
  	Queue: String Channel to put URIs in
  	hostBaseURL: String with the base url of the URI with which program was run
  	inserted: SyncMap to maintain a list of all the URIs found
  	writeOnDisk: bool variable to check if a value was supplied to STORE_ON_DISK env variable
  	insertCounter: An int counter to maintain number of URIs parsed
  	visitedCounter: An int counter to maintain number of URIs fetched*/

func createConcurrentThreads(done chan bool, queue chan string, hostBaseURL string, writeOnDisk bool, rootPath string, uriOutput bool) {
	var threads int64
	var err error
	if os.Getenv("THREAD_COUNT") != ""{
		threads, err = strconv.ParseInt(os.Getenv("THREAD_COUNT"),10,32)
		if err!=nil{
			fmt.Println("Invalid value for THREAD_COUNT env variable")
			os.Exit(1)
		}
	} else {
		threads = int64(5)
	}
	for i := 0; i < int(threads); i++ {
		go func() {
			for uri := range queue {
				enqueueAndFetch(uri, queue, hostBaseURL, writeOnDisk, rootPath, uriOutput)
			}
			done <- true
		}()
	}
	<-done
}

/*Checks if the base URL is valid and is of the format <protocol>://<baseURL>.<top-level-domain>
  Arguments:
	URL : a string argument which is passed in the args for baseURL
  Returns:
	URL: Returns a string with a valid baseURL
 */

func checkValidBaseURL (URL string) string {
	if !strings.HasPrefix(URL,"https") && !strings.HasPrefix(URL,"http") {
		URL = "https://"+URL
	}
	if !strings.Contains(URL,".") {
		URL = URL+".com"
	}
	return URL
}

/*Checks the value of STORE_ON_DISK env variable and updates it to the bool variable
  Returns:
	A boolean value after checking if STORE_ON_DISK is set
*/

func checkWriteOnDisk() bool {
	writeOnDisk, diskError := strconv.ParseBool(os.Getenv("STORE_ON_DISK"))
	if diskError != nil {
		fmt.Println("Warning: Invalid Value supplied for STORE_ON_DISK env variable")
		fmt.Println("No response will be saved to the disk")
	}
	if writeOnDisk{
		return true
	} else {
		return false
	}
}

/* This function exits once the inserted Counter is equal to the visited counter
Arguments:
	insertCounter: a int64 counter to keep track of number of URIs inserted in the channel
	visitedCounter: a int64 counter to keep track of number of URIs visited
 */

func checkCounters()  {
	if insertCounter ==  visitedCounter && insertCounter>1 {
		fmt.Println("Total Visited URIs: "+strconv.FormatInt(visitedCounter,10))
		os.Exit(1)
	}
}

/* This function stores the http responses for different URIs on disk
   Arguments:
		httpBody: The resp.body is passed as an ioReader to the function
  Returns:
		A string with the response body
 */

func storeOnDisk(httpBody io.Reader, rootPath string) string {
		buf := new(bytes.Buffer)
		_, bufReadFromerr := buf.ReadFrom(httpBody)
		if bufReadFromerr!= nil{
			fmt.Println("Error reading from buffer")
			fmt.Println(bufReadFromerr)
		}
		s := buf.String()
		strCounter := strconv.FormatInt(visitedCounter,10)
		out, _ := os.Create(rootPath+strCounter+".html")
		defer out.Close()
		_, outWriteStringErr := out.WriteString(s)
		if outWriteStringErr!= nil{
			fmt.Println("Error writing string to file")
			fmt.Println(outWriteStringErr)
		}
		return s
}

/*
The function checks for valid root directory and returns false if invalid directory is provided or no directory is provided
	Arguments:
		rootPath: A String with the value of ROOT_PATH env variable
	Returns:
		Returns a true value if directory is valid
		Returns a false value if env variable is not set or dir is invalid
 */

func checkValidDiskPath(rootPath string) bool{
	if rootPath == ""{
		fmt.Println("STORE_ON_DISK was set to True but no valid root path to create files was provided")
		return false
	} else {
		_, err := os.Open(rootPath)
		if err!=nil{
			fmt.Println("Invalid root directory")
			fmt.Println(err)
			return false
		}
	}
	return true
}

/* This function inserts URI into the channel, fetches the response from the URI
	And adds all other URIs that are to be fetched and found in response body
   Arguments:
		uri: A string that is passed to the function to fetch the response from
		insertCounter: A int64 arg to check the exit condition in checkCounters function and increment when a new URI is inserted
		visitedCounter: A int64 arg to check the exit condition in checkCounters function and increment when a new URI is visited
		writeOnDisk: A bool arg to check if the responses should be stored on disk
		inserted: A SyncMap to maintain a map of all the URIs inserted into the channel to be visited
		hostBaseURL: The base URI hostname to check for domain
		queue: The channel to insert the URIs into
*/

func enqueueAndFetch(uri string, queue chan string, hostBaseURL string, writeOnDisk bool, rootPath string, uriOutput bool) {
	links:= fetchURI(uri, writeOnDisk, rootPath, uriOutput)
	for _, link := range links {
		absolute := fixURL(link, uri)
			absoluteURL, er := url.Parse(absolute)
			if er!=nil{
				return
			}
			//Will only insert it into the channel if the hostname is same as the hostName of the URL supplied in args
			if absoluteURL.Hostname() == hostBaseURL{
				_, er := inserted.Load(absolute)
				if !er {
					atomic.AddInt64(&insertCounter,1)
					inserted.Store(absolute,true)
					inserted.Store(absolute+"/",true)
					go func() {
						queue <- absolute }()
				}
			}
	}
}

/* This function returns the string array of all relative and absolute URIs for the given uri
   Arguments:
		uri: A string that is passed to the function to fetch the response from
		insertCounter: A int64 arg to check the exit condition in checkCounters function
		visitedCounter: A int64 arg to check the exit condition in checkCounters function
		writeOnDisk: A bool arg to check if the responses should be stored on disk
		uriOutput: A bool arg to specify if the URIs should be printed after fetching them
	Returns:
		An array of strings with relative and absolute URIs present in the response body html
*/

func fetchURI(uri string, writeOnDisk bool, rootPath string, uriOutput bool) []string{
	atomic.AddInt64(&visitedCounter,1)
	resp, err := http.Get(uri)
	checkCounters()
	if uriOutput{
		fmt.Println(uri)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var links []string
	if writeOnDisk {
		s  := storeOnDisk(resp.Body, rootPath)
		links = getAllLinksHTML(strings.NewReader(s))
	} else {
		links = getAllLinksHTML(resp.Body)
	}
	defer resp.Body.Close()
	return links
}

/* The function returns the absolute uri from the relative uri as a string
   Arguments:
		href: String with the relative URI
		base: String with the base URI
   Returns:
		An absolute URI for the given relative URI as a string
 */

func fixURL(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

/*This function takes a reader object and returns a array of string slices for all the anchor links
in the html response
	Arguments:
		httpBody: Response body is passed as a reader object
	Returns :
		An array of uris as string slices found in the anchor elements
 */

func getAllLinksHTML(httpBody io.Reader) []string {
	links := []string{}
	hrefLinks := []string{}
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					poundRemovedURI := removePound(attr.Val)
					hrefLinks = append(hrefLinks, poundRemovedURI)
					checkHrefURL(&links, hrefLinks)
				}
			}
		}
	}
}

/* removePound slices the url if a # is present in the URL
   Arguments:
		uri: Takes a string argument for the uri
   Returns:
		Returns a slice of the string in case URI has a #
		In case there is no # in the string it returns the original URI

 */

func removePound(uri string) string {
	if strings.Contains(uri, "#") {
		var index int
		for n, str := range uri {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		return uri[:index]
	}
	return uri
}

/* This function checks if a string is present in the string array
   Arguments:
	uris: Takes a string array with all the uris that need to be checked
	checkString: A string to be checked for presence in the string
   Returns:
	A true value if the string is present in the array else returns false
 */

func check(uris []string, checkString string) bool {
	var check bool
	for _, str := range uris {
		if str == checkString {
			check = true
			break
		}
	}
	return check
}

/*This function appends a string to the string array for the hrefURL array
  Arguments:
	hrefURLs: Receives a pointer to the hrefURL string array
    stringSlice: Receives a string slice to be appended to the hrefURLs array
 */

func checkHrefURL(hrefURLs *[]string, stringSlice []string) {
	for _, str := range stringSlice {
		if !check(*hrefURLs, str) {
			*hrefURLs = append(*hrefURLs, str)
		}
	}
}

