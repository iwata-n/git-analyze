package main

import (
	"log"
)

// Info
var ilog *log.Logger

// Debug
var dlog *log.Logger

// Output
var output *log.Logger

func main() {

	config := ParseArgs()

	dlog, ilog, output = InitLog(config.IsDebug, config.IsShowProgress)
	dlog.Printf("config=%+v", config)

	var result ParseResult
	if config.IsSkipParse {
		result = open_result(config)
	} else {
		result = parse(config)
	}

	list := Sort(result, config)

	if config.IsSearchOnlyTargetAuthor {
		searchOnlyTargetAuthor(list, config)
	}
}
