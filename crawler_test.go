package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func TestCheck1 (t *testing.T) {
	testStringSlice := [] string {"one","two","three"}
	testString1 := "four"
	testResult1 := check(testStringSlice,testString1)
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
	testResult2 := check(testStringSlice,testString2)
	if !testResult2 {
		fmt.Println("Check function returned false for a string present in slice")
		t.Fail()
	} else {
		fmt.Println("Check function test 2 passed")
	}
}

func TestCheckHrefURLs1(t *testing.T){
	testHrefURls := [] string {"a","b","c"}
	testUris := [] string {"a","d"}
	checkHrefURL(&testHrefURls,testUris)
	if len(testHrefURls) == 5 {
		fmt.Println("CheckHrefUrl function added all the Urls into the slice and check function failed")
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

func TestCheckValidBaseURL1(t *testing.T){
	uri := "test"
	testUri := checkValidBaseURL(uri)
	if testUri != "https://test.com" {
		fmt.Println("The CheckValidBaseURL function returned an invalid value of"+testUri)
		fmt.Println("The expected value is https://test.com")
		t.Fail()
	}
	fmt.Println("Check Valid Base URL Test 1 passed")
}

func TestCheckValidBaseURL2(t *testing.T) {
	uri := "http://test"
	testUri := checkValidBaseURL(uri)
	if testUri != "http://test.com" {
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

func TestFixUrl(t *testing.T){
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

func TestStoreOnDis1(t *testing.T){
	testReader := strings.NewReader("Test text")
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testDirRoot = testDirRoot+"/testDir/"
	mkdirerr := os.Mkdir(testDirRoot,0755)
	if mkdirerr!=nil{
		fmt.Println("Creation of directory failed in Store On Disk test 2 with the error message")
		fmt.Print(mkdirerr)
		t.Fail()
	}
	testStr := storeOnDisk(testReader,testDirRoot)
	if testStr != "Test text" {
		fmt.Println("Store On Disk function returned invalid value with the output"+testStr)
		t.Fail()
	} else {
		fmt.Println("Test 1 for store on disk passed.")
	}
	os.RemoveAll(testDirRoot)
}

func TestStoreOnDisk2(t *testing.T){
	testReader := strings.NewReader("Test text")
	testDirRoot, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	testDirRoot = testDirRoot+"/testDir/"
	mkdirerr := os.Mkdir(testDirRoot,0755)
	if mkdirerr!=nil{
		fmt.Println("Creation of directory failed in Store On Disk test 2 with the error message")
		fmt.Print(mkdirerr)
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
	os.RemoveAll(testDirRoot)
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