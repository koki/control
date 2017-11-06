FROM golang

ADD . /go/src/github.com/koki/control

RUN go install github.com/koki/control/controller

ENTRYPOINT ["/go/bin/controller"]
CMD ["--help"]
