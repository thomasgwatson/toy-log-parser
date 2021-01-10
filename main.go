package main

import (
	"flag"
	"fmt"
)

var inFlag string
var outFlag string
var maxClientIpsFlag int
var maxPathsFlag int

func init() {
	flag.StringVar(&inFlag, "in", "log.txt", "input text file name")
	flag.StringVar(&outFlag, "out", "results.json", "output json file name")
	flag.IntVar(&maxClientIpsFlag, "max-client-ips", 10, "integer that defines the maximum number of results to output in the top_client_ips field")
	flag.IntVar(&maxPathsFlag, "max-paths", 10, "integer that defines the max number of results to output on the top_path_avg_seconds field.")
}

func main() {
	flag.Parse()
	//setup output struct
	//read flags

	//read in file, line by line bufio.Scanner

	// -in string
	// -out string
	// -max-client-ips int
	// -max-paths int
	// Read given file or from STDIN
	// var logReader io.Reader
	// var err error
	// logReader = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1" 200`)
	// logReader = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`)
	// logReader = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	// logReader = strings.NewReader(`89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`)
	// Use nginx config file to extract format by the name

	// Read from STDIN and use log_format to parse log records

	fmt.Println(inFlag)
	fmt.Println(outFlag)
}
