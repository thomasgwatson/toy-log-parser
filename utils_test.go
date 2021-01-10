package main

import (
	"testing"
)

var logExample = []string{
	
}

var cases =[]struct {
	Line string
	FirstMatch string
	SecondMatch string
	ThirdMatch string
	FourthMatch string
}{
	{`57.157.87.86 - - [06/Feb/2020:00:11:04 +0100] "GET /?parammore=1&customer_id=1&version=1.56&param=meaningful HTTP/1.1" 204 0 "https://www.woooo.com/more/woooo/" "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0) Gecko/20100101 Firefox/72.0"`, "57.157.87.86", "/", },
	{`34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"`},
	{`51.232.15.21 - markp [14/Jul/2009:03:35:00 +0000] "GET /product/catalog?item=fe23acd HTTP/1.1" 200 649 "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9a1) Gecko/20070308 Minefield/3.0a1"`},
	{`92.177.30.4 - - [09/Dec/2010:15:43:00 +0000] "POST /product/cart HTTP/1.1" 200 1198 "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"`},
	{`112.21.100.55 - root [22/May/2012:07:44:38 +0000] "GET /admin.php HTTP/1.1" 404 50 "python-requests/2.21.0"`},
}

func TestParser(t *testing.T) {
	t.Run("return first match from first test-case", func(t *testing.T) {
		testcase := logExample[0]
		got := ParseLine(testcase)[1]
		want := "57.157.87.86"
		assertStringEquality(t, got, want)
	})

	t.Run("return second match from first test-case", func(t *testing.T) {
		testcase := logExample[0]
		got := ParseLine(testcase)[2]
		want := "06/Feb/2020:00:11:04 +0100"
		assertStringEquality(t, got, want)
	})

	t.Run("return third match from first test-case", func(t *testing.T) {
		testcase := logExample[0]
		got := ParseLine(testcase)[3]
		want := "/?parammore=1&customer_id=1&version=1.56&param=meaningful&customer_name=somewebsite.com&some_id=4&cachebuster=1580944263903 HTTP/1.1"
		assertStringEquality(t, got, want)
	})

}

func TableTestParser(t *testing.T){

}

func assertStringEquality(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
