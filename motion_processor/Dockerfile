FROM golang:1.12-alpine3.9 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app
COPY . .

RUN go install -v ./...
FROM alpine:3.9 as motion_processor
COPY --from=peepoomonitor_compile /go/bin/motion_processor /usr/bin/motion_processor
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates 

CMD ["motion_processor"]
