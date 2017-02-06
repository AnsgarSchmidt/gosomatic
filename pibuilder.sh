#!/usr/bin/env bash

GOOS=linux GOARCH=arm GOARM=6 go build -v /Users/ansi/development/go/src/github.com/ansgarschmidt/gosomatic/infofetcher/weatherunderground.go
GOOS=linux GOARCH=arm GOARM=6 go build -v /Users/ansi/development/go/src/github.com/ansgarschmidt/gosomatic/infofetcher/radiation.go
GOOS=linux GOARCH=arm GOARM=6 go build -v /Users/ansi/development/go/src/github.com/ansgarschmidt/gosomatic/infofetcher/worldofwarcraft.py.go
