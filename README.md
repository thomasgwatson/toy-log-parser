# Log Parsing

This is a toy log parser, for nginx access logs, written in Go.

- Run the tests? `go test`
- Run the program? `go run main.go utils.go` with the appropriate arguments and input files (example available in the exampleData folder)
- You want to run it on docker? `docker build -t="log-parser" .` and then (for example) `docker run -t -i log-parser --in=my-special-log-file.txt --out results-2021-01-15.json`

You will be providing technical assistance to `Chunkfred`, a well positioned
Software as a Service company that provides insights around their clients'
applications activity and usage.

They have built `chunk-o-lytics`: an analytics engine that generates structured
data, that later on can be consumed by clients. Some market research showed the
need to build a new plugin for `chunk-o-lytics` that is able to parse `nginx` logs
and extract some specific information.

## What do I need to do?

Design and implement the plugin in question which will be shipped to `Chunkfred`s
production and bombarded with real life input data from their nginx servers, so it
can analyze them and output the important data for making decisions about company future.

Therefore, it is extremely important to design the plugin according to specifications below.

### The plugin requirements

1. The plugin must support the following command line flags:
   - **-in**: string, specifies the input file,
              and defaults to `log.txt` if not provided.
   - **-out**: string, specifies the output `JSON` file,
              and defaults to `results.json` if not provided.
   - **-max-client-ips**: integer, defines the maximum number of results
              to output in the `top_client_ips` field. Defaults to `10` if not provided.
   - **-max-paths**: integer, defined the maximum number of results
              to output on the `top_path_avg_seconds` field. Defaults to `10` if not provided.

   - both `-max-client-ips` and `-max-paths` flags can take a value from the following range: `0 to 10000`
   - when any of the command line flags are invalid, the plugin must exit with return code of `1`.

2. The plugin must create the output `JSON` file according to these rules:
   - All fields are mandatory and they need to be present in the output file
   - **total_number_of_lines_processed**: must be an integer value, containing
     total number of lines processed by the plugin.
   - **total_number_of_lines_ok**: must be an integer value, containing total
     number of lines which were parsed successfully and were accounted for in
     the output.
    - **total_number_of_lines_failed**: must be an integer value, containing total
     number of lines which failed to be parsed successfully due to format error,
     and hence didn't contribute to the end result counters.
   - **top_client_ips**: must be an array with length not more than `max-client-ips`
     command line parameter, and must contain string representation of IP addresses,
     which appeared the most in the input data.
   - **top_path_avg_seconds**: must be an array with length not more than `max-paths`
     command line parameter, and must contain string representations of the HTTP _paths_,
     which had the _slowest average response times_.
     Times should be floating point numbers with a precision of two decimals.

### Anatomy of the input log file

All log lines are generated on `Chunkfred` nginx servers, and are expected to have the following format:
```
<remote_addr> - <remote_user> [<date>] "<http_verb> <http_path> <http_version>" <http_response_code> <http_response_time_milliseconds> "<user_agent_string>"
```

However, since the plugin will consume raw unfiltered data, some discrepancies are to be expected in the input, so
the plugin should handle them properly. At this stage in development, `Chunkfred` are only interested
in deep processing of the valid log lines, and simple counting of the invalid ones. Valid log lines
are lines where every element -- `remote_addr`, `date`, `http_response_code`, etc -- are valid.
`Chunkfred` doesn't use custom extensions on their nginx servers, therefore all the values that conform
with [RFC 2616](https://tools.ietf.org/html/rfc2616) are to be considered valid, and any values that don't
match either RFC or [NginX](https://nginx.org/en/docs/http/ngx_http_core_module.html#variables) variables
specification render the whole log line invalid. Date field in valid lines is expected to be compliant
with date from [Common Log Format](https://en.wikipedia.org/wiki/Common_Log_Format) definition.

For example, the log file can contain this data:

```
34.149.47.34 - - [28/Sep/2008:23:15:00 +0000] "GET /product/catalog HTTP/1.1" 200 1531 "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:24.0) Gecko/20100101 Firefox/24.0"
51.232.15.21 - markp [14/Jul/2009:03:35:00 +0000] "GET /product/catalog?item=fe23acd HTTP/1.1" 200 649 "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9a1) Gecko/20070308 Minefield/3.0a1"
92.177.30.4 - - [09/Dec/2010:15:43:00 +0000] "POST /product/cart HTTP/1.1" 200 1198 "Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36"
112.21.100.55 - root [22/May/2012:07:44:38 +0000] "GET /admin.php HTTP/1.1" 404 50 "python-requests/2.21.0"
```

### Anatomy of the output json file

The challenge result must be a valid JSON file and have the following format:

```json
{
    "total_number_of_lines_processed": <integer>,
    "total_number_of_lines_ok": <integer>,
    "total_number_of_lines_failed": <integer>,
    "top_client_ips": {
        "ip.ad.dr.es": <integer>,
        "ip.ad.dr.es": <integer>,
        ...
    },
    "top_path_avg_seconds": {
        "/path/without/query": <float>,
        "/another/path/w/o/query": <float>,
        ...
    }
}
```

For example, correct output for the lines presented above, when called
with default command line parameters, would be this json:

```json
{
  "total_number_of_lines_processed": 4,
  "total_number_of_lines_ok": 4,
  "total_number_of_lines_failed": 0,
  "top_client_ips": {
    "92.177.30.4": 1,
    "51.232.15.21": 1,
    "34.149.47.34": 1,
    "112.21.100.55": 1
  },
  "top_path_avg_seconds": {
    "/admin.php": 0.05,
    "/product/catalog": 1.09,
    "/product/cart": 1.2
  }
}
```
