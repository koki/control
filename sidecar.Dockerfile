FROM golang

ADD . /go/src/github.com/koki/control

RUN go install github.com/koki/control/sidecar

ENTRYPOINT ["/go/bin/sidecar"]
CMD ["--help"]
