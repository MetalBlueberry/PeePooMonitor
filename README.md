# PeePoo Monitor

The objective of this application is to monitor when my cats are going to the toilet.

## How it works

A PIR motion sensor [HC-SR501](https://www.makerlab-electronics.com/product/pir-motion-sensor-hc-sr501/) is connected to the Raspberry GPIO to detect when the cats are close to the sand.

## Secrets

You must create .txt files inside secrets folder with the credentials for external services.

## TODO list

- [ ] Run everything in docker-compose environment. (problems with GPIO)
- [ ] Remove telegram credentials from git history
- [ ] Add database to store data
- [ ] Allow the user to temporally disable the system from telegram
- [ ] Send Daily/Weekly reports
- [ ] Notify of abnormal usage
- [ ] Manually register when the sand is clean and notify after X usages to clean it again
- [ ] Take a picture when movement is detected.
