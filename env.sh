#!/usr/bin/env bash
# Install develop dependencies
go get -u -v github.com/golang/dep/cmd/dep
go get -u -v github.com/axw/gocov/gocov
go get -u -v github.com/AlekSi/gocov-xml
go get -u -v gopkg.in/matm/v1/gocov-html
go get -u -v github.com/go-task/task/cmd/task
go get -u -v github.com/jstemmer/go-junit-report
