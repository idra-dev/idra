ARG GO_VERSION=1.18.2

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf

RUN mkdir -p /cdc
WORKDIR /cdc

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN ls -R
RUN go build -tags musl -o ./app ./cdc_agent/main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /cdc
WORKDIR /cdc
COPY --from=builder /cdc/app .
RUN mkdir -p /cdc/cdc_agent/
RUN mkdir -p /cdc/cdc_agent/static_configs/

ENTRYPOINT ["./app"]