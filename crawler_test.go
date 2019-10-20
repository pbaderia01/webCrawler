package main

import (
	"fmt"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func TestGetAllLinksHTML1 (t *testing.T) {
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	testLinks := getAllLinksHTML(testReader)
	if len(testLinks) != 2 {
		fmt.Print("The GetAllLinksHTML function failed with wrong number of links")
		t.Fail()
	} else {
		fmt.Println("Get All Links function passed test 1")
	}
}


func TestGetAllLinksHTML2 (t *testing.T) {
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	testLinks := getAllLinksHTML(testReader)
	if testLinks[0] != "http://testlink1.com" {
		fmt.Println("The GetAllLinksHTML function returned the first link in the test html snippet wrong.")
		t.Fail()
	} else {
		fmt.Println("Get All Links function passed test 2")
	}
}


func TestGetAllLinksHTML3 (t *testing.T) {
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	testLinks := getAllLinksHTML(testReader)
	if testLinks[1] != "http://testlink2.com" {
		fmt.Println("The GetAllLinksHTML function returned the second link in the test html snippet wrong.")
		t.Fail()
	} else {
		fmt.Println("Get All Links function passed test 3")
	}
}

func TestCheckHrefURLs1(t *testing.T){
	testHrefURls := [] string {"a","b","c"}
	testUris := [] string {"a","d"}
	checkHrefURL(&testHrefURls,testUris)
	if len(testHrefURls) == 5 {
		fmt.Println("CheckHrefUrl function added all the Urls into the slice and checkURI function failed")
		t.Fail()
	} else if len(testHrefURls) == 4 {
		fmt.Println("CheckHrefUrl passed test 1")
	}
}

func TestCheckHrefURLs2(t *testing.T){
	testHrefURls := [] string {"a","b","c"}
	testUris := [] string {"a","d"}
	checkHrefURL(&testHrefURls,testUris)
	if len(testHrefURls) == 3 {
		fmt.Println("CheckHrefUrl did not add any new elements to the slice")
		t.Fail()
	} else if len(testHrefURls) == 4 {
		fmt.Println("CheckHrefUrl passed test 2 with flying colours")
	}
}

func TestCheckWriteOnDisk1(t *testing.T){
	err := os.Setenv("STORE_ON_DISK","true")
	if err!= nil{
		fmt.Println("Failure in setting the env variable for STORE_ON_DISK")
		fmt.Println(err)
		t.Fail()
	} else {
		testStoreOnDisk := checkWriteOnDisk()
		if !testStoreOnDisk{
			fmt.Println("STORE_ON_DISK value incorrect")
			t.Fail()
		} else {
			fmt.Println("Check Write On Disk passed first test")
		}
	}
}

func TestCheckWriteOnDisk2(t *testing.T){
	err := os.Setenv("STORE_ON_DISK","false")
	if err!= nil{
		fmt.Println("Failure in setting the env variable for STORE_ON_DISK")
		fmt.Println(err)
		t.Fail()
	} else {
		testStoreOnDisk := checkWriteOnDisk()
		if testStoreOnDisk{
			fmt.Println("STORE_ON_DISK value incorrect")
			t.Fail()
		} else {
			fmt.Println("Check Write On Disk passed second test")
		}
	}
}

func TestCheckWriteOnDisk3(t *testing.T){
	err := os.Setenv("STORE_ON_DISK","invalid_value")
	if err!= nil{
		fmt.Println("Failure in setting the env variable for STORE_ON_DISK")
		fmt.Println(err)
		t.Fail()
	} else {
		testStoreOnDisk := checkWriteOnDisk()
		if testStoreOnDisk{
			fmt.Println("STORE_ON_DISK returned true for invalid value")
			t.Fail()
		} else {
			fmt.Println("Check Write On Disk passed third test")
		}
	}
}

func TestFixUrl1(t *testing.T){
	testURI := "/test1"
	testBaseURI := "https://test1.com/test2"
	testResult := fixURL(testURI,testBaseURI)
	if testResult!= "https://test1.com/test1" {
		fmt.Println("FixURL test failed with an output"+testResult)
		fmt.Println("The expected Output is https://test1.com/test1")
		t.Fail()
	} else {
		fmt.Println("FixURL test case passed")
	}
}

func TestRemovePound1(t *testing.T){
	testUri := "testuri"
	testResult := removePound(testUri)
	if testResult != testUri{
		fmt.Println("Remove Pound function returned invalid value without a # ")
		fmt.Println("Remove Pound Function test case 1 failed with the output: "+testResult)
		t.Fail()
	} else {
		fmt.Println("Remove Pound Test case 1 passed")
	}
}


func TestRemovePound2(t *testing.T){
	testUri := "testuri#test"
	testResult := removePound(testUri)
	if testResult != "testuri"{
		fmt.Println("Remove Pound function returned invalid value without a # ")
		fmt.Println("Remove Pound Function test case 2 failed with the output: "+testResult)
		t.Fail()
	} else {
		fmt.Println("Remove Pound Test case 2 passed")
	}
}

func TestStoreOnDisk1(t *testing.T){
	testReader := strings.NewReader("Test text")
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testDirRoot = testDirRoot+"/testDir/"
	mkDirErr := os.Mkdir(testDirRoot,0755)
	if mkDirErr!=nil{
		fmt.Println("Creation of directory failed in Store On Disk test 2 with the error message")
		fmt.Print(mkDirErr)
		t.Fail()
	}
	testStr := storeOnDisk(testReader,testDirRoot)
	if testStr != "Test text" {
		fmt.Println("Store On Disk function returned invalid value with the output"+testStr)
		t.Fail()
	} else {
		fmt.Println("Test 1 for store on disk passed.")
	}
	_ = os.RemoveAll(testDirRoot)
}

func TestStoreOnDisk2(t *testing.T){
	testReader := strings.NewReader("Test text")
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testDirRoot = testDirRoot+"/testDir/"
	mkDirEerr := os.Mkdir(testDirRoot,0755)
	if mkDirEerr!=nil{
		fmt.Println("Creation of directory failed in Store On Disk test 2 with the error message")
		fmt.Print(mkDirEerr)
		t.Fail()
	}
	storeOnDisk(testReader,testDirRoot)
	testFiles, _ := ioutil.ReadDir(testDirRoot)
	if len(testFiles) != 1 {
		fmt.Println("Store On Disk test returned wrong number of files")
		fmt.Println("It returned "+strconv.FormatInt(int64(len(testFiles)),10)+" file/files")
		t.Fail()
	} else {
		fmt.Println("Test 2 for store on disk passed.")
	}
	_ = os.RemoveAll(testDirRoot)
}

func TestCheckValidDiskPath1(t *testing.T){
	setEnvErr := os.Setenv("ROOT_PATH","")
	if setEnvErr != nil{
		fmt.Println("Error while setting env variable in Check Valid Disk Path Test 1")
		fmt.Println(setEnvErr)
	}
	testRootPath := os.Getenv("ROOT_PATH")
	if !checkValidDiskPath(testRootPath){
		if testRootPath != "" {
			fmt.Println("ROOT_PATH env variable not set to empty string")
			fmt.Println("Test 1 for check valid disk path failed with value"+testRootPath)
			t.Fail()
		} else {
			fmt.Println("Test 1 for check valid disk path passed")
		}
	}
}

func TestCheckValidDiskPath2(t *testing.T){
	setEnvErr := os.Setenv("ROOT_PATH","foo-bar")
	if setEnvErr != nil{
		fmt.Println("Error while setting env variable in Check Valid Disk Path Test 2")
		fmt.Println(setEnvErr)
	}
	testRootPath := os.Getenv("ROOT_PATH")
	if !checkValidDiskPath(testRootPath){
		if testRootPath != "foo-bar" {
			fmt.Println("ROOT_PATH env variable value is invalid")
			fmt.Println("Test 2 for check valid disk path failed with value"+testRootPath)
			t.Fail()
		} else {
			fmt.Println("Test 2 for check valid disk path passed")
		}
	}
}

func TestCheckValidDiskPath3(t *testing.T){
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	setEnvErr := os.Setenv("ROOT_PATH",testDirRoot)
	if setEnvErr != nil{
		fmt.Println("Error while setting env variable in Check Valid Disk Path Test 2")
		fmt.Println(setEnvErr)
	}
	testRootPath := os.Getenv("ROOT_PATH")
	if !checkValidDiskPath(testRootPath) {
		fmt.Println("Invalid value set for ROOT_PATH env variable")
		fmt.Println("Test 3 for check valid disk path failed with the output"+testRootPath)
		t.Fail()
	}else {
		fmt.Println("Test 3 for check valid disk path passed")
	}
}

func TestCheckValidBaseURL1(t *testing.T){
	uri := "test"
	testUri := checkValidBaseURL(uri)
	if testUri != "false" {
		fmt.Println("The CheckValidBaseURL function returned an invalid value of"+testUri)
		fmt.Println("The expected value is https://test.com")
		t.Fail()
	}
	fmt.Println("Check Valid Base URL Test 1 passed")
}

func TestCheckValidBaseURL2(t *testing.T) {
	uri := "http://test"
	testUri := checkValidBaseURL(uri)
	if testUri != "false" {
		fmt.Println("The CheckValidBaseURL function returned an invalid value of"+testUri)
		fmt.Println("The expected value is http//test.com")
		t.Fail()
	}
	fmt.Println("Check Valid Base URL Test 2 passed")
}

func TestCheckValidBaseURL3(t *testing.T){
	uri := "test.in"
	testUri := checkValidBaseURL(uri)
	if testUri != "https://test.in" {
		fmt.Println("The CheckValidBaseURL function returned an invalid value of"+testUri)
		fmt.Println("The expected value is https//test.in")
		t.Fail()
	}
	fmt.Println("Check Valid Base URL Test 3 passed")
}

func TestCheck1 (t *testing.T) {
	testStringSlice := [] string {"one","two","three"}
	testString1 := "four"
	testResult1 := checkURI(testStringSlice,testString1)
	if testResult1 {
		fmt.Println("Check function returned true for a string not present in slice")
		t.Fail()
	} else {
		fmt.Println("Check function test 1 passed")
	}
}


func TestCheck2 (t *testing.T) {
	testStringSlice := [] string {"one","two","three"}
	testString2 := "one"
	testResult2 := checkURI(testStringSlice,testString2)
	if !testResult2 {
		fmt.Println("Check function returned false for a string present in slice")
		t.Fail()
	} else {
		fmt.Println("Check function test 2 passed")
	}
}

func TestCheckCounters1(t *testing.T) {
	testQueue := make(chan string)
	testDone := make(chan bool)
	atomic.AddInt64(&insertCounter,1)
	testInsertValue := "testValue"
	go func() {
		testQueue <- testInsertValue
	}()
	checkCounters(testQueue,testDone)
	go func() {
		for {
			testValue, signal := <- testQueue
			if !signal{
				fmt.Println("Channel was closed by checkCounters. Test 1 for checkCounters failed")
				fmt.Println("The value that should be read from testQueue is"+testValue)
				t.Fail()
			} else {
				fmt.Println("Test 1 for checkCounters passed!")
			}
		}
	}()
	atomic.SwapInt64(&insertCounter,0)
}

func TestCheckDisplay1(t *testing.T){
	setEnvErr := os.Setenv("DISPLAY_URI","true")
	if setEnvErr!=nil{
		fmt.Println("Error occurred while setting env variable")
		fmt.Println(setEnvErr)
	}
	testCheckDisplay := checkDisplay()
	if !testCheckDisplay{
		fmt.Println("Check Display function returned an invalid value")
		t.Fail()
	} else {
		fmt.Println("Test 1 for check display passed!")
	}
	_ = os.Setenv("DISPLAY_URI", "")
}

func TestCheckDisplay2(t *testing.T){
	setEnvErr := os.Setenv("DISPLAY_URI","false")
	if setEnvErr!=nil{
		fmt.Println("Error occurred while setting env variable")
		fmt.Println(setEnvErr)
	}
	testCheckDisplay := checkDisplay()
	if testCheckDisplay{
		fmt.Println("Check Display function returned an invalid value")
		t.Fail()
	} else {
		fmt.Println("Test 2 for check display passed!")
	}
	_ = os.Setenv("DISPLAY_URI", "")
}

func TestCheckDisplay3(t *testing.T){
	setEnvErr := os.Setenv("DISPLAY_URI","invalid value")
	if setEnvErr!=nil{
		fmt.Println("Error occurred during setting env variable")
		fmt.Println(setEnvErr)
	}
	testCheckDisplay := checkDisplay()
	if testCheckDisplay{
		fmt.Println("Check Display function returned an invalid value")
		t.Fail()
	} else {
		fmt.Println("Test 3 for check display passed!")
	}
	_ = os.Setenv("DISPLAY_URI", "")
}

func TestGetBaseHostName1(t *testing.T){
	testURI := "https://testuri.com"
	testResult := getBaseHostname(testURI)
	if testResult != "testuri.com"{
		fmt.Println("The getBaseHostName function returned an invalid value")
		fmt.Println("getBaseHostName function failed with output"+testResult)
	} else {
		fmt.Println("GetBaseHostName function passed test.")
	}
}

func TestGetRootPath1(t *testing.T){
	writeOnDisk := true
	rootPathErr := os.Setenv("ROOT_PATH","testValue")
	if rootPathErr!=nil{
		fmt.Println("Error while setting ROOT_PATH env value in test 1")
		fmt.Println(rootPathErr)
	}
	testResult := getRootPath(&writeOnDisk)
	if writeOnDisk {
		fmt.Println("Value of write on disk was not modified")
		t.Fail()
	} else if testResult != "" {
		fmt.Println("getRootPath returned an invalid value")
		fmt.Println("getRootPath failed with the output"+testResult)
		t.Fail()
	} else {
		fmt.Println("getRootPath function passed test case 1!")
	}
	_ = os.Setenv("ROOT_PATH", "")
}

func TestGetRootPath2(t *testing.T){
	writeOnDisk := false
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	rootPathErr := os.Setenv("ROOT_PATH",testDirRoot)
	if rootPathErr!=nil{
		fmt.Println("Error while setting ROOT_PATH env value in test 1")
		fmt.Println(rootPathErr)
	}
	testResult := getRootPath(&writeOnDisk)
	if writeOnDisk {
		fmt.Println("Value of write on disk was modified")
		t.Fail()
	} else if testResult != "" {
		fmt.Println("getRootPath returned an invalid value")
		fmt.Println("getRootPath failed with the output"+testResult)
		t.Fail()
	} else if !writeOnDisk && testResult==""{
		fmt.Println("getRootPath function passed test case 2!")
	}
	_ = os.Setenv("ROOT_PATH", "")
}


func TestGetRootPath3(t *testing.T){
	writeOnDisk := true
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	rootPathErr := os.Setenv("ROOT_PATH",testDirRoot)
	if rootPathErr!=nil{
		fmt.Println("Error while setting ROOT_PATH env value in test 1")
		fmt.Println(rootPathErr)
	}
	testResult := getRootPath(&writeOnDisk)
	if !writeOnDisk {
		fmt.Println("Value of write on disk was modified")
		t.Fail()
	} else if testResult != testDirRoot {
		fmt.Println("getRootPath returned an invalid value")
		fmt.Println("getRootPath failed with the output"+testResult)
		t.Fail()
	} else if writeOnDisk && testResult==testDirRoot{
		fmt.Println("getRootPath function passed test case 3!")
	}
	_ = os.Setenv("ROOT_PATH", "")
}

func TestGetCrawlURI1(t *testing.T){
	presentOsArgs := os.Args
	os.Args = [] string{"cmd","testValue"}
	setEnvErr := os.Setenv("CRAWL_URL","")
	if setEnvErr != nil{
		fmt.Println("Error while setting crawl uri env var")
		fmt.Println(setEnvErr)
	}
	testResult := getCrawlURI()
	if testResult != "testValue"{
		fmt.Println("getCrawlURI returned an invalid value")
		fmt.Println("Test 1 for getCrawlURI failed with the output"+testResult)
		t.Fail()
	} else {
		fmt.Println("Test 1 for getCrawlURI function passed.")
	}
	os.Args = presentOsArgs
}

func TestGetCrawlURI2(t *testing.T){
	presentOsArgs := os.Args
	os.Args = [] string{"cmd","testValue"}
	setEnvErr := os.Setenv("CRAWL_URL","testValue2")
	if setEnvErr != nil{
		fmt.Println("Error while setting crawl uri env var")
		fmt.Println(setEnvErr)
	}
	testResult := getCrawlURI()
	if testResult == "testValue"{
		fmt.Println("getCrawlURI set an invalid value")
		fmt.Println("ENV variable has higher precedence which was broken by the function")
		t.Fail()
	} else if testResult!="testValue2" {
		fmt.Println("getCrawlURI returned an invalid value")
		fmt.Println("Test 2 for getCrawlURI failed with the output"+testResult)
		t.Fail()
	}else {
		fmt.Println("Test 2 for getCrawlURI function passed.")
	}
	os.Args = presentOsArgs
}


func TestGetCrawlURI3(t *testing.T){
	setEnvErr := os.Setenv("CRAWL_URL","testValue3")
	if setEnvErr != nil{
		fmt.Println("Error while setting crawl uri env var")
		fmt.Println(setEnvErr)
	}
	testResult := getCrawlURI()
	if testResult!="testValue3" {
		fmt.Println("getCrawlURI returned an invalid value")
		fmt.Println("Test 3 for getCrawlURI failed with the output"+testResult)
		t.Fail()
	}else {
		fmt.Println("Test 3 for getCrawlURI function passed.")
	}
}

func TestGetStringFromReader(t *testing.T){
	testReader := strings.NewReader("Test Text")
	testResult := getStringFromReader(testReader)
	if testResult != "Test Text" {
		fmt.Println("getStringFromReader returned an invalid value")
		fmt.Println("getStringFromReader test failed with an output"+testResult)
		t.Fail()
	} else if testResult == "Test Text" && reflect.TypeOf(testResult)!= reflect.TypeOf("string") {
		fmt.Println("getStringFromReader returned a correct value but an incorrect type")
		fmt.Println(reflect.TypeOf(testResult))
		t.Fail()
	} else if testResult == "Test Text" && reflect.TypeOf(testResult)== reflect.TypeOf("string"){
		fmt.Println("getStringFromReader passed the test case")
	}
}

func TestGetThreadCount1(t *testing.T){
	setEnvErr := os.Setenv("THREAD_COUNT","4")
	if setEnvErr != nil{
		fmt.Println("Error while setting env var in get thread count variable")
		fmt.Println(setEnvErr)
	}
	testResult := getThreadCount()
	if testResult==5 {
		fmt.Println("getThreadCount not setting the value of the thread using the env var")
		t.Fail()
	} else if testResult != 5 && testResult != 4{
		fmt.Println("getThreadCount failed with an invalid value")
		fmt.Print(testResult)
		t.Fail()
	} else {
		fmt.Println("Test Case 1 for getThreadCount Passed")
	}
	setEnvErr = os.Setenv("THREAD_COUNT","")
	if setEnvErr != nil{
		fmt.Println("Error while setting env var in get thread count variable")
		fmt.Println(setEnvErr)
	}
}

func TestGetThreadCount2(t *testing.T){
	setEnvErr := os.Setenv("THREAD_COUNT","")
	if setEnvErr != nil{
		fmt.Println("Error while setting env var in get thread count variable")
		fmt.Println(setEnvErr)
	}
	testResult := getThreadCount()
	 if testResult != 5 {
		fmt.Println("getThreadCount failed with an invalid value")
		fmt.Print(testResult)
		t.Fail()
	} else {
		fmt.Println("Test Case 2 for getThreadCount Passed")
	}
	setEnvErr = os.Setenv("THREAD_COUNT","")
	if setEnvErr != nil{
		fmt.Println("Error while setting env var in get thread count variable")
		fmt.Println(setEnvErr)
	}
}

func TestUriOutputStore1(t *testing.T){
	testHttpBodyReader := strings.NewReader("Test Text")
	testUri := "test URI"
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testResultReader := uriOutputStore(testHttpBodyReader, false,testUri, false,testDirRoot)
	testStringResult := getStringFromReader(testResultReader)
	if testStringResult!="Test Text"{
		fmt.Println("uriOutputStore returned an invalid value"+testStringResult)
		t.Fail()
	} else if testStringResult == "Test Text" && reflect.TypeOf(testResultReader) != reflect.TypeOf(strings.NewReader("string")){
		fmt.Println("uriOutputStore returned correct value but incorrect type")
		fmt.Print(reflect.TypeOf(testResultReader))
		t.Fail()
	} else if testStringResult == "Test Text" && reflect.TypeOf(testResultReader) == reflect.TypeOf(strings.NewReader("string")){
		fmt.Println("Test 1 for uriOutputStore passed!")
	}
}

func TestUriOutputStore2(t *testing.T){
	testHttpBodyReader := strings.NewReader("Test Text")
	testUri := "test URI"
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testDirRoot = testDirRoot +"/uriOutputStore2/"
	_ = os.Mkdir(testDirRoot, 0755)
	testResultReader := uriOutputStore(testHttpBodyReader, false,testUri, true,testDirRoot)
	testStringResult := getStringFromReader(testResultReader)
	testFiles, _ := ioutil.ReadDir(testDirRoot)
	if testStringResult!="Test Text"{
		fmt.Println("uriOutputStore returned an invalid value"+testStringResult)
		t.Fail()
	} else if testStringResult == "Test Text" && reflect.TypeOf(testResultReader) != reflect.TypeOf(strings.NewReader("string")){
		fmt.Println("uriOutputStore returned correct value but incorrect type")
		fmt.Print(reflect.TypeOf(testResultReader))
		t.Fail()
	} else if testStringResult == "Test Text" && reflect.TypeOf(testResultReader) == reflect.TypeOf(strings.NewReader("string")) && len(testFiles)!=1 {
		fmt.Println("uriOutputStore returned correct value and a correct type but didn't store the output on disk")
		fmt.Println("Number of files on disk are")
		fmt.Print(len(testFiles))
		t.Fail()
	} else if testStringResult == "Test Text" && reflect.TypeOf(testResultReader) == reflect.TypeOf(strings.NewReader("string")) && len(testFiles)==1{
		fmt.Println("Test 2 for uriOutputStore passed!")
	}
}

func TestFetchURI(t *testing.T){
	defer gock.Off()
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	initialTestString := getStringFromReader(testReader)
	gock.New("http://testfetchuri.com").Get("/test").Reply(200).BodyString(initialTestString)
	testResultReader := fetchURI("http://testfetchuri.com/test")
	testResultString := getStringFromReader(testResultReader)
	if initialTestString != testResultString {
		fmt.Println("Invalid value returned by fetchURI")
		fmt.Println(testResultString)
		t.Fail()
	} else if initialTestString == testResultString && reflect.TypeOf(testResultReader) != reflect.TypeOf(testReader){
		fmt.Println("Correct Value returned but incorrect type returned by fetchURI")
		fmt.Println(reflect.TypeOf(testResultReader))
		t.Fail()
	} else if initialTestString == testResultString && reflect.TypeOf(testResultReader) == reflect.TypeOf(testReader){
		fmt.Println("Test 1 for Fetch URI passed")
	}
}

func TestFetchURI2(t *testing.T){
	defer gock.Off()
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	initialTestString := getStringFromReader(testReader)
	gock.New("http://testfetchuri.com").Get("/test").Reply(500).BodyString(initialTestString)
	testResultReader := fetchURI("http://testfetchuri.com/test")
	test5XXString := "URI returned a 5xx status code"
	testResultString := getStringFromReader(testResultReader)
	if testResultString != test5XXString {
		fmt.Println("Invalid value returned by fetchURI")
		fmt.Println(testResultString)
		t.Fail()
	} else if test5XXString == testResultString && reflect.TypeOf(testResultReader) != reflect.TypeOf(testReader){
		fmt.Println("Correct Value returned but incorrect type returned by fetchURI")
		fmt.Println(reflect.TypeOf(testResultReader))
		t.Fail()
	} else if test5XXString == testResultString && reflect.TypeOf(testResultReader) == reflect.TypeOf(testReader){
		fmt.Println("Test 2 for Fetch URI passed")
	}
}

func TestFetchURI3(t *testing.T){
	defer gock.Off()
	testReader := strings.NewReader(
		` <p>
  				<a href="http://testlink1.com">1</a>
				<a style=\"\" href=http://testlink2.com>3</a>
 					 http://negativetestlink.com
			</p>`)
	initialTestString := getStringFromReader(testReader)
	gock.New("http://testfetchuri.com").Get("/test").Reply(450).BodyString(initialTestString)
	testResultReader := fetchURI("http://testfetchuri.com/test")
	test4XXString := "URI returned a 4xx status code"
	testResultString := getStringFromReader(testResultReader)
	if testResultString != test4XXString {
		fmt.Println("Invalid value returned by fetchURI")
		fmt.Println(testResultString)
		t.Fail()
	} else if test4XXString == testResultString && reflect.TypeOf(testResultReader) != reflect.TypeOf(testReader){
		fmt.Println("Correct Value returned but incorrect type returned by fetchURI")
		fmt.Println(reflect.TypeOf(testResultReader))
		t.Fail()
	} else if test4XXString == testResultString && reflect.TypeOf(testResultReader) == reflect.TypeOf(testReader){
		fmt.Println("Test 3 for Fetch URI passed")
	}
}

func TestInsertInitialURI(t *testing.T){
	testQueue := make(chan string)
	uri := "http://test.com"
	insertInitialURI(uri,testQueue)
	test := <- testQueue
	_,testBool1 := inserted.Load(uri)
	_,testBool := inserted.Load(uri+"/")
	if test != uri{
		fmt.Println("Invalid value inserted in testQueue"+test)
		t.Fail()
	} else if test == uri && !testBool1 && !testBool {
		fmt.Println("The channel returned the correct value but failed to store the uri in the sync Map")
		t.Fail()
	} else if test == uri && testBool1 && testBool{
		fmt.Println("Test 1 for initalURI passed")
	}
	atomic.SwapInt64(&insertCounter,0)
	atomic.SwapInt64(&visitedCounter,0)
	close(testQueue)
}

func TestFilterAndEnqueue1(t *testing.T){
	testLinks := []string{"/test1","https://test.com/test2","/test1"}
	testQueue := make(chan string,3)
	testDone := make(chan bool)
	testHostBaseURL := "test.com"
	testCrawlURI := "https://test.com"
	testCounter := 0
	filterAndEnqueue(testLinks,testQueue,testHostBaseURL,testCrawlURI)
	go func() {
		for {
			testValue,signal := <-testQueue
			//fmt.Println(signal)
			//fmt.Println(testValue)
			if signal{
				testCounter+=1
			} else {
				testDone <- true
				return
			}
			 if testCounter >= 3 {
					fmt.Println("The test counter reached 3 or more with an invalid value ",testValue)
					t.Fail()
				}
			}
	}()
}

func TestFilterAndEnqueue2(t *testing.T){
	testLinks := []string{"/test1","https://test.com/test2","/test1"}
	testQueue := make(chan string,3)
	testDone := make(chan bool)
	testHostBaseURL := "testing.com"
	testCrawlURI := "https://test.com"
	testCounter := 0
	filterAndEnqueue(testLinks,testQueue,testHostBaseURL,testCrawlURI)
	go func() {
		for {
			testValue,signal := <-testQueue
			//fmt.Println(signal)
			//fmt.Println(testValue)
			if signal{
				testCounter+=1
			} else {
				testDone <- true
				return
			}
			fmt.Println(testCounter)
			if testCounter == 1 {
				fmt.Println("The counter should not reach here!")
				fmt.Println("An invalid value was entered ", testValue)
				t.Fail()
			}
		}
	}()
}

func TestFilterAndEnqueue3(t *testing.T){
	testLinks := []string{"/test1","https://test.com/test2","/test1"}
	testQueue := make(chan string,3)
	testDone := make(chan bool)
	testHostBaseURL := "test.com"
	testCrawlURI := "https://test1.com"
	testCounter := 0
	filterAndEnqueue(testLinks,testQueue,testHostBaseURL,testCrawlURI)
	go func() {
		for {
			testValue,signal := <-testQueue
			if signal{
				testCounter+=1
			} else {
				testDone <- true
				return
			}
			fmt.Println(testCounter)
			if testCounter >= 1 {
				fmt.Println("The counter should not reach here!")
				fmt.Println("An invalid value was entered ", testValue)
				t.Fail()
			}
		}
	}()
}