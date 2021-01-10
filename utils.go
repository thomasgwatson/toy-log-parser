package main

import "regexp"

func ParseLine(line string) (match []string) {
	var logRegex = regexp.MustCompile(`^(\S+)[\s-]+\[([^][]*)]\s+"GET\s+([^"]+)"[^"](\d{3})[^"](\d+)`)
	match = logRegex.FindStringSubmatch(line)
	return
}
