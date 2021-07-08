package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

// Logger
var ilog *log.Logger
var dlog *log.Logger

// file path
type FilePath string

type CommitFile struct {
	// file path
	Path FilePath

	// Author list
	Authors []string

	CommitHash []string

	CreateBy string
}

func initLog(isDebug bool, isShowProgress bool) {
	ilog = log.New(io.Discard, "", log.Ltime|log.Lmicroseconds|log.LUTC)
	dlog = log.New(io.Discard, "[DEBUG]", log.Ltime|log.Lmicroseconds|log.LUTC)
	if isShowProgress {
		ilog.SetOutput(os.Stdout)
	}

	if isDebug {
		dlog.SetOutput(os.Stdout)
		ilog.SetOutput(os.Stdout)
		ilog.SetPrefix("[INFO ]")
	}
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	log.Fatal(fmt.Sprintf("error: %s", err))
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if e == a {
			return true
		}
	}
	return false
}

var emptyChange object.ChangeEntry

func name(c *object.Change) string {
	if c.From != emptyChange {
		return c.From.Name
	}
	return c.To.Name
}

func parseCommitLog(cIter object.CommitIter, depth int) (map[FilePath]CommitFile, error) {
	files := make(map[FilePath]CommitFile)
	count := 0
	err := cIter.ForEach(func(c *object.Commit) error {
		// ignore marge commit
		if len(c.ParentHashes) > 1 {
			dlog.Println("ignore marge commit")
			return nil
		}
		count++
		if count > depth && depth > 0 {
			dlog.Printf("over depth count=%d, depth=%d", count, depth)
			return nil
		}

		ilog.Printf("commit=%d hash:%40s author:%-30s\r", count, c.Hash.String(), c.Author.Name)

		fromTree, err := c.Tree()
		checkIfError(err)

		toTree := &object.Tree{}
		if c.NumParents() != 0 {
			firstParent, err := c.Parents().Next()
			if err != nil {
				return nil
			}

			toTree, err = firstParent.Tree()
			if err != nil {
				return nil
			}
		}

		// very slow...
		diff, err := toTree.Diff(fromTree)
		checkIfError(err)

		for _, v := range diff {
			dlog.Println(v)

			path := FilePath(name(v))
			action, _ := v.Action()
			var createBy string = ""
			if action == merkletrie.Insert {
				createBy = c.Author.Name
			}

			if val, ok := files[path]; ok {
				dlog.Printf("exist %s", path)
				if !contains(val.Authors, c.Author.Name) {
					val.Authors = append(val.Authors, c.Author.Name)
				}
				val.CommitHash = append(val.CommitHash, c.Hash.String())
				if createBy != "" {
					val.CreateBy = createBy
				}
				files[path] = val
			} else {
				dlog.Printf("new file %s", path)
				files[path] = CommitFile{
					Path:       path,
					Authors:    []string{c.Author.Name},
					CommitHash: []string{c.Hash.String()},
					CreateBy:   createBy,
				}
			}
		}
		return nil
	})
	dlog.Printf("\r\ntotal commit=%d\n", count)

	return files, err
}

func main() {
	var branch, path, outputFile string
	var depth int
	var isDebug, isShowProgress bool
	flag.StringVar(&branch, "branch", "master", "branch")
	flag.StringVar(&path, "path", "./", "path")
	flag.StringVar(&outputFile, "parse-file", "", "result file")
	flag.IntVar(&depth, "depth", -1, "depth")
	flag.BoolVar(&isDebug, "debug", false, "debug")
	flag.BoolVar(&isShowProgress, "show-progress", false, "show progress")
	flag.Parse()
	initLog(isDebug, isShowProgress)
	dlog.Printf("path=%s branch=%s depth=%d", path, branch, depth)

	r, err := git.PlainOpen(path)
	checkIfError(err)

	head, err := r.Head()
	checkIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	checkIfError(err)

	files, err := parseCommitLog(cIter, depth)
	checkIfError(err)

	j, err := json.Marshal(files)
	checkIfError(err)

	var buf bytes.Buffer
	json.Indent(&buf, j, "", "  ")

	fmt.Printf("%s\n", buf.String())
	if outputFile != "" {
		err := ioutil.WriteFile(outputFile, buf.Bytes(), 0666)
		checkIfError(err)
	}
}
