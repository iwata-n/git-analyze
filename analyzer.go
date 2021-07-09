package main

func searchOnlyTargetAuthor(parseResult ParseResult, authors ArgsAuthor) ([]CommitFile, error) {
	dlog.Print("searchOnlyTargetAuthor")
	var err error
	files := []CommitFile{}

	for _, v := range parseResult {
		dlog.Printf("%+v", v)
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

	return files, err
}
