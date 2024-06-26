FROM golang:1.22 as build

COPY . /go/src/github.com/days365/pipeflow

WORKDIR /go/src/github.com/days365/pipeflow

RUN go get && go build

FROM golang:1.22

COPY --from=build /go/src/github.com/days365/pipeflow/pipeflow /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/pipeflow"]
