#! /bin/sh
set -e

 count=`git ls-files | grep '.go$' | xargs gofmt -l -s | wc -l`
 if [ $count -gt 0 ]; then
     echo "Files not formatted correctly\n"
     exit 1
 fi
 go vet .
 go test -race -v .
 go install -race -v .
