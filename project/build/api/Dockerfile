FROM golang:1.20 AS builder

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/app

# Copy go.mod and go.sum
COPY app/go.* ./

RUN go mod download

COPY app/ ./

RUN go build -o main .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /go/src/app/main /main

EXPOSE 8080

CMD [ "./main" ]
