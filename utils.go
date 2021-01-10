package main

import "regexp"

func ParseLine(line string) (match []string) {
	var logRegex = regexp.MustCompile(`^(\S+)[\s-]-\s\S+\s\[[^][]*]\s"\S+\s(/[^\?\s]*)\?*\S*\s[^"]+"\s\d{3}\s(\d+)\s"([^"]+)"`)
	//regex workspace https://regex101.com/r/I7EPUI/3
	match = logRegex.FindStringSubmatch(line)
	return
}
