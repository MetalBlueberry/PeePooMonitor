# FROM golang:1.12-alpine3.9 as builder

# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh

# WORKDIR /app
# COPY . .

# RUN go install -v ./...

FROM alpine:3.9 as sensor
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates 
COPY --from=peepoomonitor_compile /go/bin/sensor /usr/bin/sensor

CMD ["sensor"]
