package main

import (
	"strings"
)

func searchOnlyTargetAuthor(parseResult []CommitFile, config Config) {
	output.Print("searchOnlyTargetAuthor")
	files := []CommitFile{}
	authors := config.Authors

	for _, v := range parseResult {
		isContain := true
		for _, a := range v.Authors {
			if contains(authors, a) {
				continue
			}
			isContain = false
		}

		if isContain {
			files = append(files, v)
		}
	}

	output.Println("commit count, path, authors")
	for _, f := range files {
		output.Printf("%d, %s, [%s]\n", len(f.CommitHash), f.Path, strings.Join(f.Authors, " "))
	}
}
