package main

import "sort"

func Sort(result ParseResult, config Config) []CommitFile {
	var list []CommitFile
	for _, v := range result {
		list = append(list, v)
	}

	if config.sortPath {
		dlog.Println("sortPath")
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Path < list[j].Path
		})
	}

	if config.sortCommitCount {
		dlog.Println("sortCommitCount")
		sort.SliceStable(list, func(i, j int) bool {
			return len(list[i].CommitHash) > len(list[j].CommitHash)
		})
	}

	return list
}
