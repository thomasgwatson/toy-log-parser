package main

// Thank you for reading my code. As it stands, 'in production' there are possible 'out of memory' concerns
// And possible invalid lines slipping into the results set.
// Of course, if this was real-life, I would anticipate that production logging like splunk or honeycomb.io
// could be used to zero in on these issues. Beyond that, there are numerous possible test-cases one could add
// to expand the rigor of line validation.
// I hope you enjoy inspecting this code as much as I did writing it; it is my first program using Go and I am really digging it

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	flag.IntVar(&maxClientIpsFlag, "max-client-ips", 10, "integer that defines the maximum number of results to output in the top_client_ips field")
	flag.IntVar(&maxPathsFlag, "max-paths", 10, "integer that defines the max number of results to output on the top_path_avg_seconds field.")
}

func main() {
	flag.Parse()
	// validate intFlags
	if !IntFlagIsValid(maxPathsFlag) || !IntFlagIsValid(maxClientIpsFlag) {
		fmt.Println("Flags invalid")
		os.Exit(1)
	}

	// Consider validation of inFlag and outFlag

	// open input file
	file, errFile := os.Open(inFlag)

	if errFile != nil {
		fmt.Println("file wouldn't open")
		os.Exit(1)
	}
	defer file.Close()

	// scan file line by line, parse lines and append to matches
	scanner := bufio.NewScanner(file)
	var matches Matches
	for scanner.Scan() {
		match := ParseLine(scanner.Text())
		if match == nil {
			failed := []string{}
			matches = append(matches, failed) // done for type consistency
		} else {
			matches = append(matches, match)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// compile matches from file into intermediate format
	results := CompileResults(matches)
	// sort/trim results and port to JSON
	transformedResults := TransformResults(results, maxClientIpsFlag, maxPathsFlag)
	json_representation, errJson := json.MarshalIndent(transformedResults, "", " ")

	if errJson != nil {
		log.Fatal(errJson)
	}

	fmt.Println(string(json_representation))
	// write to output file
	errWrite := ioutil.WriteFile(outFlag, json_representation, 0644)
	if errWrite != nil {
		log.Fatal(errWrite)
	}
}
