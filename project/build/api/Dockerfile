FROM golang:1.19-alpine

WORKDIR /go/src/app

RUN apk update && \
    apk add git && \
    apk add --no-cache alpine-conf

ADD ./app /go/src/app

RUN cd /go/src/app && \
    go install app

ENTRYPOINT ["/go/bin/app"]

WORKDIR "/go"
