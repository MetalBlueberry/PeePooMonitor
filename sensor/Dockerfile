FROM golang:1.12 as dev
WORKDIR /app
ENV SHELL /bin/bash

FROM golang:1.12-alpine3.9 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app
COPY . .

RUN go install -v ./...

FROM alpine:3.9 as monitor
COPY --from=builder /go/bin/monitor /usr/bin/monitor
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates


CMD ["monitor"]
