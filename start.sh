#!/bin/bash

git pull
if [ "nowp" != "$1" ]
then
    echo "Building webpack..."
    webpack
fi
export GOPATH=$GOPATH:$(pwd)
go build github.com/enderian/confessions-go/
./confessions-go