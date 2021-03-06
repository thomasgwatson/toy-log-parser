package main

import (
	"reflect"
	"testing"
)

type Case struct {
	Line        string
	FirstMatch  string
	SecondMatch string
	ThirdMatch  string
	FourthMatch string
	FifthMatch  string
	SixthMatch  string
}

var cases = []Case{
	{`34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`, `34.149.47.34`, `28/Sep/2008:23:15:00 +0000`, `GET`, `/product/catalog`, `1531`, `Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0`},
	{`51.232.15.21 - markp [14/Jul/2009:03:35:00 +0000] "GET /product/catalog?item=fe23acd HTTP/1.1" 200 649 "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9a1) Gecko/20070308 Minefield/3.0a1"`, `51.232.15.21`, `14/Jul/2009:03:35:00 +0000`, `GET`, `/product/catalog`, `649`, `Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9a1) Gecko/20070308 Minefield/3.0a1`},
	{`92.177.30.4 - - [09/Dec/2010:15:43:00 +0000] "POST /product/cart HTTP/1.1" 200 1198 "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"`, `92.177.30.4`, `09/Dec/2010:15:43:00 +0000`, `POST`, `/product/cart`, `1198`, `Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36`},
	{`112.21.100.55 - root [22/May/2012:07:44:38 +0000] "GET /admin.php HTTP/1.1" 404 50 "python-requests/2.21.0"`, `112.21.100.55`, `22/May/2012:07:44:38 +0000`, `GET`, `/admin.php`, `50`, `python-requests/2.21.0`},
	{`112.21.100.55 - - [doesntmatter] "HEAD /kind/of/valid/url?anything somethinghere" 999 1234 "literallyany text you like"`, `112.21.100.55`, `doesntmatter`, `HEAD`, `/kind/of/valid/url`, `1234`, `literallyany text you like`},
	{`3.112.21.21 - root [06/Jun/2020:02:24:12 +0000] "HEAD %5C%22%3Bdrop%20table%20users%3B%20select%5C%22 HTTP/1.1" 404 5503 "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)"`, `3.112.21.21`, `06/Jun/2020:02:24:12 +0000`, `HEAD`, `%5C%22%3Bdrop%20table%20users%3B%20select%5C%22`, `5503`, `Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)`},
}

var badCases = []string{
	`34.149.47.34 - - 28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
	`34.149.47.34 [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
	`34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
	// `<remote_addr> - <remote_user> [<date>] "<http_verb> <http_path> <http_version>" <http_response_code> <http_response_time_milliseconds> "<user_agent_string>"
}

var topClientIps = TopClientIps{
	`3.112.21.21`:   1,
	`34.149.47.34`:  1,
	`51.232.15.21`:  1,
	`92.177.30.4`:   1,
	`112.21.100.55`: 2,
}
var topPathAvgSeconds = TopPathAvgSeconds{
	`%5C%22%3Bdrop%20table%20users%3B%20select%5C%22`: TopPathRecord{1, 5503},
	`/product/catalog`:   TopPathRecord{2, 2180},
	`/product/cart`:      TopPathRecord{1, 1198},
	`/admin.php`:         TopPathRecord{1, 50},
	`/kind/of/valid/url`: TopPathRecord{1, 1234},
}
var goodResults = Results{6, 6, 0, topClientIps, topPathAvgSeconds}

func CaseStructToSlice(cases []Case) (result [][]string) {
	for _, c := range cases {
		slice := []string{
			c.Line,
			c.FirstMatch,
			c.SecondMatch,
			c.ThirdMatch,
			c.FourthMatch,
			c.FifthMatch,
			c.SixthMatch,
		}
		result = append(result, slice)
	}
	return
}

func TestParser(t *testing.T) {

	for _, testcase := range cases {
		t.Run(`return first match from test-cases`, func(t *testing.T) {
			got := ParseLine(testcase.Line)[1]
			want := testcase.FirstMatch
			assertStringEquality(t, got, want)
		})

		t.Run("return second match from test-cases", func(t *testing.T) {
			got := ParseLine(testcase.Line)[4]
			want := testcase.FourthMatch
			assertStringEquality(t, got, want)
		})

		t.Run("return third match from test-cases", func(t *testing.T) {
			got := ParseLine(testcase.Line)[6]
			want := testcase.SixthMatch
			assertStringEquality(t, got, want)
		})
	}

	t.Run("returns nothing on deformed input", func(t *testing.T) {
		for _, testcase := range badCases {
			got := ParseLine(testcase)

			if got != nil {
				t.Errorf("Expected nil, instead got %v", got)
			}
		}
	})
}

func TestCompileResults(t *testing.T) {
	input := CaseStructToSlice(cases)
	got := CompileResults(input)

	if !reflect.DeepEqual(got, goodResults) {
		t.Errorf("got %v want %v", got, goodResults)
	}
}

// Obviously would want to complete this test in real-life circumstances

// func TestTransformResults(t *testing.T){
// 	maxClientIpsFlag := 2
// 	maxPathsFlag := 2
// 	got := TransformResults(goodResults, maxClientIpsFlag, maxPathsFlag)

// 	if !reflect.DeepEqual(got, goodTransform) {
// 		t.Errorf("got %v want %v", got, goodTransform)
// 	}
// }

func TestIntFlagValidation(t *testing.T) {
	t.Run("valid flag is valid", func(t *testing.T) {
		got := IntFlagIsValid(2)
		if !got {
			t.Errorf("Expected flag to be valid and true, got false")
		}
	})

	t.Run("negative out of range int flag returns as invalid", func(t *testing.T) {
		got := IntFlagIsValid(-234)
		if got {
			t.Errorf("Expected flag to be invalid and false, got true")
		}
	})

	t.Run("positive out of range flag returns as invalid", func(t *testing.T) {
		got := IntFlagIsValid(23400)
		if got {
			t.Errorf("Expected flag to be invalid and false, got true")
		}
	})
}

func TestResults(t *testing.T) {
	topClientIps := TopClientIps{}
	topPathAvgSeconds := TopPathAvgSeconds{}
	results := Results{
		0, 0, 0, topClientIps, topPathAvgSeconds,
	}

	t.Run("Update", func(t *testing.T) {
		expectedTopClientIps := TopClientIps{}
		expectedTopPathAvgSeconds := TopPathAvgSeconds{}
		expected := Results{
			1, 1, 0, expectedTopClientIps, expectedTopPathAvgSeconds,
		}

		results.TotalNumberOfLinesOk += 1
		results.TotalNumberOfLinesProcessed += 1

		if !reflect.DeepEqual(expected, results) {
			t.Errorf("got %v want %v", expected, results)
		}

	})

}

func TestTopPathRecord(t *testing.T) {
	t.Run("Simplest use-case for TopPathRecord", func(t *testing.T) {
		record := TopPathRecord{}
		record.Update(100)
		got := record.AverageResponseTime()
		want := 100.0

		if got != want {
			t.Errorf("expected '%f' but got '%f'", want, got)
		}
	})

	t.Run("Test multiple Updates", func(t *testing.T) {
		record := TopPathRecord{}
		record.Update(100)
		record.Update(200)
		got := record.AverageResponseTime()
		want := 150.0

		if got != want {
			t.Errorf("expected '%f' but got '%f'", want, got)
		}
	})
}

func assertStringEquality(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
