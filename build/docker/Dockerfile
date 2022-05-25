FROM golang:1.17

RUN apt-get update && \
    go get github.com/githubnemo/CompileDaemon

ENTRYPOINT make run-compile-daemon