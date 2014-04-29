#/usr/bin/env/sh

set -o nounset
set -o errexit

# go fmt check
readonly GO_FMT_COUNT=`go fmt ./... | wc -l`
if [ $GO_FMT_COUNT -ne 0 ]
then
    echo $GO_FMT_COUNT files need to be reformatted with go fmt.
fi
