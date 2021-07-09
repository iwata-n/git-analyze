# git-analyze

git commit log analyze tool

## install
```
brew tap iwata-n/git-analyze
brew install iwata-n/git-analyze/git-analyze
```

## usage
Output the parse results to result.json

```
git analyze -parse-file=result.json
```

Check for files that have been committed only by the specified Author.

```
git analyze -search-only-target-author -author=<AUTHOR1> -author=<AUTHOR2> -author=<more...>
```