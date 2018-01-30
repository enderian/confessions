#!/usr/bin/env bash
set -e
npm install
go get github.com/valyala/fasthttp
go get github.com/tyler-sommer/stick
go get github.com/google/uuid
go get gopkg.in/mgo.v2
go get github.com/buaazp/fasthttprouter