FROM alpine:3.9 as motion_processor
COPY --from=peepoomonitor_compile /go/bin/motion_processor /usr/bin/motion_processor
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates 

CMD ["motion_processor"]
