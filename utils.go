package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
)

func ParseLine(line string) (match []string) {
	var logRegex = regexp.MustCompile(`^(\S+)[\s-]-\s\S+\s\[[^][]*]\s"\S+\s([^\?\s]*)\?*\S*\s[^"]+"\s\d{3}\s(\d+)\s"[^"]+"\s{0,1}.*`)
	// regex workspace https://regex101.com/r/I7EPUI/3
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

type TopPathRecord struct {
	requests                int
	culminativeResponseTime int
}

func (t *TopPathRecord) Update(responseTime int) {
	t.requests += 1
	t.culminativeResponseTime += responseTime
}

func (t *TopPathRecord) AverageResponseTime() float64 {
	if t.culminativeResponseTime == 0 {
		return 0
	}
	if t.requests == 0 {
		return 0
	}
	return float64(t.culminativeResponseTime) / float64(t.requests)
}

type TopClientIps map[string]int
type TopPathAvgSeconds map[string]TopPathRecord
type Matches [][]string

type Results struct {
	TotalNumberOfLinesProcessed int
	TotalNumberOfLinesOk        int
	TotalNumberOfLinesFailed    int
	TopClientIps                TopClientIps
	TopPathAvgSeconds           TopPathAvgSeconds
}

func CompileResults(matches Matches) Results {
	topClientIps := TopClientIps{}
	topPathAvgSeconds := TopPathAvgSeconds{}
	results := Results{0, 0, 0, topClientIps, topPathAvgSeconds}

	for _, matchSeries := range matches {
		results.TotalNumberOfLinesProcessed += 1
		if len(matchSeries) == 0 {
			results.TotalNumberOfLinesFailed += 1
			continue
		}

		ipAddress := matchSeries[1]
		path := matchSeries[2]
		responseTime, _ := strconv.Atoi(matchSeries[3])
		results.TotalNumberOfLinesOk += 1
		updateTopClientIps(results.TopClientIps, ipAddress)
		updateTopPathAvgSeconds(results.TopPathAvgSeconds, path, responseTime)
	}

	return results
}

func updateTopClientIps(topClientIps TopClientIps, ipAddress string) {
	_, ok := topClientIps[ipAddress]
	if !ok {
		topClientIps[ipAddress] = 1
		return
	}
	topClientIps[ipAddress] += 1
}

func updateTopPathAvgSeconds(topPathAvgSeconds TopPathAvgSeconds, path string, responseTime int) {
	topPathRecord, ok := topPathAvgSeconds[path]
	if !ok {
		topPathAvgSeconds[path] = TopPathRecord{1, responseTime}
	} else {
		topPathRecord.Update(responseTime)
		topPathAvgSeconds[path] = topPathRecord
	}
}

type TopPathAvgReponseOutput map[string]int

type TransformedResults struct {
	TotalNumberOfLinesProcessed int                `json:"total_number_of_lines_processed"`
	TotalNumberOfLinesOk        int                `json:"total_number_of_lines_ok"`
	TotalNumberOfLinesFailed    int                `json:"total_number_of_lines_failed"`
	TopClientIps                map[string]int     `json:"top_client_ips"`
	TopPathAvgSeconds           map[string]float64 `json:"top_path_avg_seconds"`
}

func TransformResults(results Results, maxClientIpsFlag, maxPathsFlag int) TransformedResults {
	topClientIps := transformTopClientIps(results.TopClientIps, maxClientIpsFlag)
	topPathAvgSeconds := transformTopPathAvgSeconds(results.TopPathAvgSeconds, maxPathsFlag)
	transformedResults := TransformedResults{
		results.TotalNumberOfLinesProcessed,
		results.TotalNumberOfLinesOk,
		results.TotalNumberOfLinesFailed,
		topClientIps,
		topPathAvgSeconds,
	}
	return transformedResults
}

type PairInt struct {
	Key   string
	Value int
}

type PairListInt []PairInt

func (p PairListInt) Len() int           { return len(p) }
func (p PairListInt) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairListInt) Less(i, j int) bool { return p[i].Value > p[j].Value }

type PairFloat struct {
	Key   string
	Value float64
}

type PairListFloat []PairFloat

func (p PairListFloat) Len() int           { return len(p) }
func (p PairListFloat) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairListFloat) Less(i, j int) bool { return p[i].Value > p[j].Value }

func transformTopPathAvgSeconds(input TopPathAvgSeconds, maxPathsFlag int) map[string]float64 {
	pairList := make(PairListFloat, len(input))

	i := 0
	for k, v := range input {
		pairList[i] = PairFloat{k, toFixed(v.AverageResponseTime(), 2)}
		i++
	}

	sort.Sort(pairList)
	
	fmt.Print(pairList)
	if len(pairList) > maxPathsFlag {
		pairList = pairList[:maxPathsFlag]
	}
	fmt.Print(pairList)
	backToMap := map[string]float64{}
	for _, pair := range pairList {
		backToMap[pair.Key] = pair.Value
	}

	return backToMap
}

func transformTopClientIps(input TopClientIps, maxClientIpsFlag int) map[string]int {
	pairList := make(PairListInt, len(input))

	i := 0
	for k, v := range input {
		pairList[i] = PairInt{k, v}
		i++
	}

	sort.Sort(pairList)
	if len(pairList) > maxClientIpsFlag {
		pairList = pairList[:maxClientIpsFlag]
	}

	backToMap := map[string]int{}
	for _, pair := range pairList {
		backToMap[pair.Key] = pair.Value
	}

	return backToMap
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
