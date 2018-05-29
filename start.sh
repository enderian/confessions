#!/bin/bash
set -e
git pull
if [ "--skip-webpack" != "$1" ]
then
    echo "Building webpack..."
    npm run build
fi
export GOPATH=$GOPATH:$(pwd)
go build github.com/enderian/confessions
./confessions