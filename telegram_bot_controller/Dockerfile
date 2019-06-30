FROM alpine:3.9 as telegram_bot_controller
COPY --from=peepoomonitor_compile /go/bin/telegram_bot_controller /usr/bin/telegram_bot_controller
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates 

CMD ["telegram_bot_controller"]
