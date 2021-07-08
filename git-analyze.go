package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

var ilog *log.Logger
var dlog *log.Logger

type FilePath string

type CommitFile struct {
	// file path
	Path FilePath

	// Author list
	Authors []string

	CommitHash []string

	// Create
	CreateBy string
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

func parseCommitLog(cIter object.CommitIter, depth int) (map[FilePath]CommitFile, error) {
	files := make(map[FilePath]CommitFile)
	count := 0
	err := cIter.ForEach(func(c *object.Commit) error {
		// ignore marge commit
		if len(c.ParentHashes) > 1 {
			return nil
		}
		count++
		if count > depth && depth > 0 {
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

			path := FilePath(v.From.Name)
			if val, ok := files[path]; ok {
				if !contains(val.Authors, c.Author.Name) {
					val.Authors = append(val.Authors, c.Author.Name)
				}
				val.CommitHash = append(val.CommitHash, c.Hash.String())
				files[path] = val
			} else {
				action, _ := v.Action()
				var createBy string = ""
				if action == merkletrie.Insert {
					createBy = c.Author.Name
				}
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

func main() {
	var branch, path string
	var depth int
	var isDebug, isShowProgress bool
	flag.StringVar(&branch, "branch", "master", "branch")
	flag.StringVar(&path, "path", "./", "path")
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

	// cIter, err := r.CommitObjects()
	checkIfError(err)

	files, err := parseCommitLog(cIter, depth)
	checkIfError(err)

	for _, file := range files {
		fmt.Printf("path=%s, authors=[%s]\n", file.Path, strings.Join(file.Authors, ","))
	}
}
