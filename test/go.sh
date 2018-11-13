#!/bin/bash -eu

trap "relax keychain reset" EXIT

test/setup.sh

echo "Run go tests"

>/dev/null pushd go
export GOPATH="$PWD"
go get ./src/relparser/relfile/...
go test -v ./src/relparser/relfile/...
go get ./src/lspp/...
go test -v ./src/lspp/...
>/dev/null popd

test/teardown.sh
