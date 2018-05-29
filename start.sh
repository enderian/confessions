#!/bin/bash
set -e
git pull
if [ "--skip-webpack" != "$1" ]
then
    echo "Building webpack..."
    npm install
    npm run build
fi
go build github.com/enderian/confessions
./confessions