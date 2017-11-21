FROM golang

ADD . /go/src/github.com/koki/control

RUN go install github.com/koki/control/cli

ENTRYPOINT ["/go/bin/cli"]
CMD ["--help"]
