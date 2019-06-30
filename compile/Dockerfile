FROM golang:1.12-alpine3.9 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app
COPY . .

WORKDIR /app/sensor
RUN go install 
WORKDIR /app/motion_processor
RUN go install 
WORKDIR /app/telegram_bot_controller
RUN go install 
WORKDIR /app

RUN chmod +x /app/compile/build.sh
CMD [ "/app/compile/build.sh" ]