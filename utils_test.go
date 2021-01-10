package main

import (
	"testing"
)

var cases = []struct {
	Line        string
	FirstMatch  string
	SecondMatch string
	ThirdMatch  string
}{
	{`34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`, `34.149.47.34`, `/product/catalog`, `1531`},
	{`51.232.15.21 - markp [14/Jul/2009:03:35:00 +0000] "GET /product/catalog?item=fe23acd HTTP/1.1" 200 649 "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9a1) Gecko/20070308 Minefield/3.0a1"`, `51.232.15.21`, `/product/catalog`, `649`},
	{`92.177.30.4 - - [09/Dec/2010:15:43:00 +0000] "POST /product/cart HTTP/1.1" 200 1198 "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"`, `92.177.30.4`, `/product/cart`, `1198`},
	{`112.21.100.55 - root [22/May/2012:07:44:38 +0000] "GET /admin.php HTTP/1.1" 404 50 "python-requests/2.21.0"`, `112.21.100.55`, `/admin.php`, `50`},
}

var badCases = []string{
	`34.149.47.34 - - 28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
	`34.149.47.34 [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
	`34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`,
}

func TestParser(t *testing.T) {

	for _, testcase := range cases {
		t.Run(`return first match from test-cases`, func(t *testing.T) {
			got := ParseLine(testcase.Line)[1]
			want := testcase.FirstMatch
			assertStringEquality(t, got, want)
		})

		t.Run("return second match from test-cases", func(t *testing.T) {
			got := ParseLine(testcase.Line)[2]
			want := testcase.SecondMatch
			assertStringEquality(t, got, want)
		})

		t.Run("return third match from test-cases", func(t *testing.T) {
			got := ParseLine(testcase.Line)[3]
			want := testcase.ThirdMatch
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

func TestIntFlagValidation(t *testing.T) {
	t.Run("valid flag is valid", func(t *testing.T){
		got := IntFlagIsValid(2)
		if !got {
			t.Errorf("Expected flag to be valid and true, got false")
		}
	})

	t.Run("negative out of range int flag returns as invalid", func(t *testing.T){
		got := IntFlagIsValid(-234)
		if got {
			t.Errorf("Expected flag to be invalid and false, got true")
		}
	})

	t.Run("positive out of range flag returns as invalid", func(t *testing.T){
		got := IntFlagIsValid(23400)
		if got {
			t.Errorf("Expected flag to be invalid and false, got true")
		}
	})
}


func assertStringEquality(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
