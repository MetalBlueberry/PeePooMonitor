FROM golang:1.12-alpine3.9 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app
COPY . .

RUN go install -v ./...

FROM alpine:3.9 as telegram_bot_controller
COPY --from=builder /go/bin/telegram_bot_controller /usr/bin/telegram_bot_controller
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates 

CMD ["telegram_bot_controller"]
