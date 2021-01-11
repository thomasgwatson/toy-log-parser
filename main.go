package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var inFlag string
var outFlag string
var maxClientIpsFlag int
var maxPathsFlag int

func init() {
	flag.StringVar(&inFlag, "in", "log.txt", "input text file name")
	flag.StringVar(&outFlag, "out", "results.json", "output json file name")
	flag.IntVar(&maxClientIpsFlag, "max-client-ips", 2, "integer that defines the maximum number of results to output in the top_client_ips field")
	flag.IntVar(&maxPathsFlag, "max-paths", 2, "integer that defines the max number of results to output on the top_path_avg_seconds field.")
}

func main() {
	flag.Parse()

	if !IntFlagIsValid(maxPathsFlag) || !IntFlagIsValid(maxClientIpsFlag) {
		fmt.Println("Flags invalid")
		os.Exit(1)
	}

	// setup output struct
	// read flags

	// read in file, line by line bufio.Scanner

	// -in string
	// -out string
	// -max-client-ips int
	// -max-paths int
	// Read given file or from STDIN
	// var logReader io.Reader
	// var err error

	// Read from STDIN and use log_format to parse log records

	file, errFile := os.Open(inFlag)

	if errFile != nil {
		fmt.Println("file wouldn't open")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var matches Matches

	for scanner.Scan() {
		match := ParseLine(scanner.Text())
		if match == nil {
			failed := []string{}
			matches = append(matches, failed)
		} else {
			matches = append(matches, match)
		}
	}

	results := CompileResults(matches)
	// sort/trim results and port to JSON
	// write to file
	transformedResults := TransformResults(results, maxClientIpsFlag, maxPathsFlag)

	// Looking pretty good to here
	// - Need to figure out an efficient way to port the Pair structs into JSONable format
	// - Then I need to marshall into JSON and file to file
	
	json_representation, errJson := json.Marshal(transformedResults)

	fmt.Println(string(json_representation))

	if errJson != nil {

		log.Fatal(errJson)
}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(inFlag)
	// fmt.Println(outFlag)
}
