FROM golang:1.18 as build

COPY . /go/src/github.com/days365/pipeflow

WORKDIR /go/src/github.com/days365/pipeflow

RUN go get && go build

FROM golang:1.18

COPY --from=build /go/src/github.com/days365/pipeflow/pipeflow /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/pipeflow"]
