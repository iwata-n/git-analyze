package main

import (
	"fmt"
	"log"
	"strings"
)

var ilog *log.Logger
var dlog *log.Logger

func main() {

	config := ParseArgs()

	dlog, ilog = InitLog(config.IsDebug, config.IsShowProgress)
	dlog.Printf("config=%+v", config)

	var result ParseResult
	if config.IsSkipParse {
		result = open_result(config)
	} else {
		result = parse(config)
	}

	if config.IsSearchOnlyTargetAuthor {
		r, err := searchOnlyTargetAuthor(result, config.Authors)
		checkIfError(err)

		// fmt.Printf("%+v\n", JsonString(f))
		for _, f := range r {
			fmt.Printf("%s, %s\n", f.Path, strings.Join(f.Authors, ", "))
		}
	}
}
