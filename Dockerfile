FROM golang:alpine

## Add git
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Install develop dependencies
RUN go get -u -v github.com/go-task/task/cmd/task && \
    go get -u -v github.com/axw/gocov/gocov && \
    go get -u -v github.com/AlekSi/gocov-xml && \
    go get -u -v github.com/jstemmer/go-junit-report