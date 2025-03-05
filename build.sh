#!/usr/bin/env bash

cd danmu-core
go build -o ../bin/core/dcore ./cmd/main
cd ../danmu-http
go build -o ../bin/http/dhttp ./cmd/main

