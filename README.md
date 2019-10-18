# GO-CRAWLER
This project is an attempt to build a web crawler using golang and implement concurrency

The project provides multiple options which are configurable using environment variables as follows:

| Option | Environment Variable | Values Accepted | Default Value | Description | Required |
| --- | --- | --- | --- | --- | --- |
| Concurrency | THREAD_COUNT | Integer | 5 | This option lets you control the concurrency at which the crawler runs defaulting to 5 | False |
| URL to crawl | CRAWL_URL | String | - | This lets you configure the URL which you want to crawl and should be provided | True |
| Root Path | ROOT_PATH | String | - | This lets you configure the root path in which responses should be saved if you want to save responses to the disk. Needs to be set to a valid directory path if STORE_ON_DISK is set to True | False |
| Output Control | DISPLAY_URI | Boolean | false | This lets you configure if you want to view the URIs that are being visited by the crawler | False |
| Store On Disk | STORE_ON_DISK | Boolean | false | This lets you configure if you want to save the responses fetched on the local disk | False |

## Usage
Prerequisites: 
- Before running the go script locally please install go version 1.13.1
- Install all dependencies using `go mod download`
```
The source code can be run with defaults as:
You can run the source code as: go run crawler.go <URL>
You can also use the env variable to specify URL: CRAWL_URL=<URL> go run crawler.go
```
All other options can be configured as env variables by either setting them as a env variable or supplying the env variable with go command as:
```
ENV_VAR_1=value go run crawler.go
```

##Examples
To run with a concurrency of 3:
```
THREAD_COUNT=3 go run crawler.go <URL> 
```
To store responses on disk:
```
STORE_ON_DISK=true ROOT_PATH=/Users/piyushbaderia/response/ go run crawler.go <URL>
```
To Display URIs that are being crawled:
```
DISPLAY_URI=true go run crawler.go <URL>
```

## Features
The crawler performs the following tasks:
- Crawls a single subdomain i.e. the base domain of the URI passed to crawl on
- Option to view URIs that are being crawled
- Option to store the responses on local
- Provides control over concurrency
- The requests timeout after 30 sec

## Enhancements
The crawler can be enhanced on the following points:
- More tests : The crawler currently does not have tests for the functions that need to fetch data over the internet
- BenchMark Tests: Benchmark tests need to be added to the crawler to benchmark performance for every change
- Support for Robots.txt: Crawler currently does not support robots.txt restrictions
- Transport Configuration: More http client options should be added to support for TLS 