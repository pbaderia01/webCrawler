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
	"time"
)

var insertCounter int64 //A counter to keep a track of the number of URIs inserted into the channel
var visitedCounter int64 //A counter to keep a track of the number of URIs visited
var inserted sync.Map //A syncMap to keep a track of the URIs parsed by the HTML

func main() {
	crawlURI := getCrawlURI()
	crawlURI = checkValidBaseURL(crawlURI)
	_,crawlErr := strconv.ParseBool(crawlURI)
	if crawlErr==nil{
		fmt.Println("REACHED")
		fmt.Println(crawlURI)
	}
	hostBaseURL := getBaseHostname(crawlURI)
	writeOnDisk := checkWriteOnDisk()
	queue := make(chan string)
	done := make(chan bool)
	insertInitialURI(crawlURI,queue)
	rootPath := getRootPath(&writeOnDisk)
	uriOutput := checkDisplay()
	threads := getThreadCount()
	createConcurrentThreads(done,queue,hostBaseURL, writeOnDisk, rootPath, uriOutput, threads, crawlURI)
	close(done)
}

//Specifies the usage instructions for the source code to be run

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Please Use the module as : go run <filename>.go <URI>\n")
}

/* The function prompts the user with a message on how to run the source files in case
	user doesn't provide an initial input
	Returns:
		Checks first for the value of the env var CRAWL_URL and if it is set returns the value set in CRAWL_URL
		Else checks the flags provided by user with the source code to run and if not found exits with an error message
 */

func getCrawlURI() string{
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
	return args[0]
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
		fmt.Println("Invalid URI provided. No top level domain provided in the URI")
		return "false"
	}
	return URL
}

/*Gets the value of the crawlURI's hostname
	Arguments:
		crawlURI: A string with the value of the initial URI provided by the user.
	Returns:
		A string with the hostname of the crawlURI
*/

func getBaseHostname(crawlURI string) string{
	hostURL, err := url.Parse(crawlURI)
	if err!=nil{
		fmt.Println("Invalid URL provided. Please find the error below:")
		fmt.Println(err)
		os.Exit(1)
	}
	return hostURL.Hostname()
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
	return writeOnDisk
}

/* Inserts initial crawlURI into the queue to start processing
   Arguments:
		crawlURI: A string that specifies the initial URI
 */

func insertInitialURI(crawlURI string, queue chan string){
	inserted.Store(crawlURI,true)
	inserted.Store(crawlURI+"/",true)
	go func() { queue <- crawlURI}()
	atomic.AddInt64(&insertCounter,1)
}

/*
	The function returns the value of the rootPath if the value specified in the ROOT_PATH env variable is valid
	Else sets the writeOnDisk boolean flag to false and returns an empty string
	Arguments:
		A pointer to writeOnDisk boolean flag
	Returns:
		A valid root path as a string or sets the writeOnDisk flag to false
 */

func getRootPath(writeOnDisk *bool) string{
	var rootPath string
	if *writeOnDisk{
		rootPath = os.Getenv("ROOT_PATH")
		if !checkValidDiskPath(rootPath){
			fmt.Println("Not saving any files to disk")
			*writeOnDisk = false
			rootPath = ""
		}
	}
	return rootPath
}

/*
	The function checks the value of the DISPLAY_URI env variable and returns false if an invalid value is provided
	or returns the value of the DISPLAY_URI env variable
 */

func checkDisplay() bool{
	uriOutputFlag := os.Getenv("DISPLAY_URI")
	uriOutput, err:= strconv.ParseBool(uriOutputFlag)
	if err != nil{
		fmt.Println("Invalid value specified for display URI")
		fmt.Println("No URI will be printed")
		return false
	}
	return uriOutput
}

/*
	The function checks the value of the env variable THREAD_COUNT if the value specified in the env variable
	is invalid an error message is generated and the program exits.
	Defaults to 5
	Returns:
		An int64 with the number of threads.
 */

func getThreadCount() int64{
	var threads int64
	if os.Getenv("THREAD_COUNT") != ""{
		var err error
		threads, err = strconv.ParseInt(os.Getenv("THREAD_COUNT"),10,64)
		if err != nil{
			fmt.Println("Invalid value for THREAD_COUNT env variable")
			os.Exit(1)
		}
	} else {
		threads = int64(5)
	}
	return threads
}


/* The function takes the following arguments and starts execution with a default value of 5
Arguments:
  	Done : Bool channel to synchronize thread execution
  	Queue: String Channel to put URIs in
  	hostBaseURL: String with the base url of the URI with which program was run
  	writeOnDisk: bool variable to check if a value was supplied to STORE_ON_DISK env variable
  	rootPath: String with the root path on the disk to which response is to be saved
	uriOutput: Bool fetched from DISPLAY_URI to determine if the visited URI should be printed
	threads: An int64 with the number of threads fetched from THREAD_COUNT defaults to 5
	crawlURI: A string variable with the initial URI provided that must be crawled
 */

func createConcurrentThreads(done chan bool, queue chan string, hostBaseURL string, writeOnDisk bool, rootPath string, uriOutput bool, threads int64, crawlURI string) {
	for i := 0; i < int(threads); i++ {
		go func() {
			for uri := range queue {
				httpBody := fetchURI(uri)
				httpBodyReader := uriOutputStore(httpBody,uriOutput,uri,writeOnDisk,rootPath)
				links := getAllLinksHTML(httpBodyReader)
				filterAndEnqueue(links, queue, hostBaseURL, crawlURI)
				checkCounters(queue,done)
			}
			done <- true
		}()
	}
	<-done
}

/*
	The function checks the values of insertCounter and visitedCounter and checks for equality
	On equality the function closes the queue channel and throws the flow back to the original function
	Arguments:
		queue: A string channel that is used to maintain a list of the URIs
		done: A channel to synchronize the execution of multiple threads
	Return:
		If equal prints the number of URIs visited and closes queue channel and the flow is thrown back
 */

func checkCounters(queue chan string, done chan bool)  {
	if insertCounter ==  visitedCounter {
		fmt.Println("Total Visited URIs: "+strconv.FormatInt(visitedCounter,10))
		close(queue)
		<- done
	}
}

/* The function stores the string in the reader object and stores at the rootPath 
	Arguments:
		httpBody: A io reader containing the response body
		rootPath: A string containing the path of the root directory
	Returns:
		A string with the response body
 */

func storeOnDisk(httpBody io.Reader, rootPath string) string {
	buf := new(bytes.Buffer)
	_, bufReadFromErr := buf.ReadFrom(httpBody)
	if bufReadFromErr!= nil{
		fmt.Println("Error reading from buffer")
		fmt.Println(bufReadFromErr)
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
	The function checks the value of the rootPath variable and returns a boolean 
	Whether rootPath is a valid disk path 
 */

func checkValidDiskPath(rootPath string) bool{
	if rootPath == ""{
		fmt.Println("STORE_ON_DISK was set to True but no valid root path to create files was provided")
		return false
	} else {
		_, err := os.Open(rootPath)
		if err!=nil {
			fmt.Println("Invalid root directory")
			fmt.Println(err)
			return false
		}
	}
	return true
}

/*
	The function takes an array of strings containing all the links in the HTML response, converts the relative URIs
	to absolute URIs, filters them based on the hostname and then inserts them into the queue channel to be processed
	Arguments:
		links: An array containing all the relative and absolute URIs
		queue: A string channel to insert the URIs in
		hostBaseURL: The hostname of the crawlURI (the initial URI provided by user)
		crawlURI: The initial URI provided by the user to resolve references using
 */

func filterAndEnqueue(links []string, queue chan string, hostBaseURL string, crawlURI string) {
	for _, link := range links {
		absolute := fixURL(link, crawlURI)
		absoluteURL, er := url.Parse(absolute)
		if er!=nil{
			return
		}
		//Will only insert it into the channel if the hostname is same as the hostName of the URL supplied in args
		if absoluteURL.Hostname() == hostBaseURL{
			_, _ = inserted.LoadOrStore(absolute+"/",true)
			_, er := inserted.LoadOrStore(absolute,true)
			if er!=true{
				atomic.AddInt64(&insertCounter,1)
				go func() {
					queue <- absolute }()
			}
		}
	}
}

/*
	The function creates an httpClient with a request timeout of 30 seconds.
	It closes the response body io.ReadCloser object and returns a new io.Reader object with the response body
	Arguments:
		uri: A string with the value of the uri from which the response is to be fetched
	Returns:
		An io.Reader with the response body
 */

func fetchURI(uri string) io.Reader{
	var httpClient = &http.Client{
		Timeout: 30*time.Second,
	}

	resp, reqErr := httpClient.Get(uri)
	if reqErr!=nil{
		fmt.Println("Error while fetching response")
		fmt.Println(reqErr)
	}

	atomic.AddInt64(&visitedCounter,1)
	response := getStringFromReader(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		response = "URI returned a 4xx status code"
	} else if resp.StatusCode >= 500 && resp.StatusCode <= 599{
		response = "URI returned a 5xx status code"
	}
	return strings.NewReader(response)
}

/*
	The function prints on the terminal window if DISPLAY_URI=true and stores on disk if STORE_ON_DISK=true
	this also returns an io.Reader object with the httpResponse
	Arguments:
		httBody: An io.Reader object with the response body
		uriOutput: A bool argument with the value of DISPLAY_URI if the value set is valid bool
		uri: The uri that is being visited and should be printed
		writeOnDisk: A bool argument with the value of STORE_ON_DISK if the value set is valid bool
		rootPath: A string argument with the root path of the disk (if the value set in ROOT_PATH is a valid path)
	Returns:
		An io reader object with the response body
 */

func uriOutputStore(httpBody io.Reader, uriOutput bool, uri string, writeOnDisk bool, rootPath string) io.Reader{
	if uriOutput{
		fmt.Println(uri)
	}
	if writeOnDisk {
		s  := storeOnDisk(httpBody, rootPath)
		return strings.NewReader(s)
	} else {
		return httpBody
	}
}

/* The function returns a string from a io.Reader object
   Arguments:
       httpBody: An io.reader from which to read and convert to string
   Returns:
       A string that was buffered in the io.Reader object
 */

func getStringFromReader(httpBody io.Reader) string{
	buf := new(bytes.Buffer)
	_, bufReadFromErr := buf.ReadFrom(httpBody)
	if bufReadFromErr!= nil{
		fmt.Println("Error reading from buffer")
		fmt.Println(bufReadFromErr)
	}
	s := buf.String()
	return s
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
	var links []string
	var hrefLinks []string
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

func checkURI(uris []string, checkString string) bool {
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
		if !checkURI(*hrefURLs, str) {
			*hrefURLs = append(*hrefURLs, str)
		}
	}
}