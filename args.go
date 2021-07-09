package main

import "flag"

type ArgsAuthor []string

func (i *ArgsAuthor) String() string {
	return "my string representation"
}

func (i *ArgsAuthor) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Config struct {
	Branch, Path, OutputFile string

	Depth int

	IsDebug, IsShowProgress, IsSearchOnlyTargetAuthor, IsSkipParse bool

	sortCommitCount, sortPath bool

	Authors ArgsAuthor
}

func ParseArgs() Config {
	var branch, path, outputFile string
	var depth int
	var isDebug, isShowProgress, isSearchOnlyTargetAuthor, isSkipParse bool
	var sortCommitCount, sortPath bool
	var authors ArgsAuthor

	flag.StringVar(&branch, "branch", "master", "git branch")
	flag.StringVar(&path, "path", "./", "repository path")
	flag.StringVar(&outputFile, "parse-file", "", "result file")
	flag.IntVar(&depth, "depth", -1, "commit depth")
	flag.Var(&authors, "author", "author")
	flag.BoolVar(&isDebug, "debug", false, "output debug message")
	flag.BoolVar(&isSkipParse, "skip-parse", false, "skip parse")
	flag.BoolVar(&isShowProgress, "show-progress", false, "show progress")
	flag.BoolVar(&isSearchOnlyTargetAuthor, "search-only-target-author", false, "Search for files with only the target author.")
	flag.BoolVar(&sortCommitCount, "sort-commit-count", true, "sort commit count")
	flag.BoolVar(&sortPath, "sort-path", true, "sort path")
	flag.Parse()

	return Config{
		Branch:                   branch,
		Path:                     path,
		OutputFile:               outputFile,
		Depth:                    depth,
		Authors:                  authors,
		IsDebug:                  isDebug,
		IsShowProgress:           isShowProgress,
		IsSearchOnlyTargetAuthor: isSearchOnlyTargetAuthor,
		IsSkipParse:              isSkipParse,
		sortCommitCount:          sortCommitCount,
		sortPath:                 sortPath,
	}
}
