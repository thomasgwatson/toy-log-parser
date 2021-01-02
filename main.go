package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello world")

	// -in string
	// -out string
	// -max-client-ips int
	// -max-paths int

	args := os.Args[1:]
	fmt.Println(args)

	wordPtr := flag.String("in", "log.txt", "input file")
	flag.Parse()
	fmt.Println(*wordPtr)
}
