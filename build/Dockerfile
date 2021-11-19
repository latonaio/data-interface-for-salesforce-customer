# syntax = docker/dockerfile:1.0.0-experimental
# Build Container
FROM golang:1.17.3 as builder

ENV GO111MODULE on
ENV GOPRIVATE "github.com/latonaio"
WORKDIR /go/src/github.com/latonaio

COPY go.mod .

RUN go build -o data-interface-for-salesforce-contract-get .

COPY . .

RUN go build -o data-interface-for-salesforce-customer .

# Runtime Container
FROM alpine:3.14

RUN apk add --no-cache libc6-compat tzdata

COPY --from=builder /go/src/github.com/latonaio/data-interface-for-salesforce-customer .

CMD ["./data-interface-for-salesforce-customer"]

