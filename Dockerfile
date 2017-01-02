FROM golang:1.7.4
MAINTAINER bewiwi <bewiwi@gmail.com>

WORKDIR /go/src/github.com/bewiwi/mta/
COPY / /go/src/github.com/bewiwi/mta/

RUN go build -o mta main.go

ENTRYPOINT ["./mta"]