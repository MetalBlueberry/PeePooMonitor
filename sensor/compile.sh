#!/bin/bash
cmd cmd/monitor
env GOOS=linux GOARCH=arm GOARM=5 go build 
sleep 1
ssh pi@192.168.1.131 rm monitor
scp monitor pi@192.168.1.131:
sleep 1
ssh pi@192.168.1.131 ./monitor

