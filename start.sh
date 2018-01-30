#!/bin/bash

git pull
if [ "--skip-webpack" != "$1" ]
then
    echo "Building webpack..."
    npm run build
fi
export GOPATH=$GOPATH:$(pwd)
go build ender.gr/confessions
./confessions