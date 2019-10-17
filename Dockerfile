FROM golang:1.13.1
LABEL maintainer="Piyush Baderia <piyush.baderia@outlook.com"
WORKDIR /app
COPY crawler.go ./
