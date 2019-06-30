#!/bin/sh
cd /app/sensor
go install 
cd /app/motion_processor
go install 
cd /app/telegram_bot_controller
go install 