package main

import "regexp"

func ParseLine(line string) (match []string) {
	var logRegex = regexp.MustCompile(`^(\S+)[\s-]-\s\S+\s\[[^][]*]\s"\S+\s(/[^\?\s]*)\?*\S*\s[^"]+"\s\d{3}\s(\d+)\s"([^"]+)"`)
	//regex workspace https://regex101.com/r/I7EPUI/3
	match = logRegex.FindStringSubmatch(line)
	return
}

func IntFlagIsValid(intFlag int) bool {
	if intFlag > 10000 {
		return false
	}
	if intFlag < 0 {
		return false
	}
	return true
}
