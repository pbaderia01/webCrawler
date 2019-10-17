FROM golang:1.13.1
LABEL maintainer="Piyush Baderia <piyush.baderia@outlook.com"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY crawler.go ./
ENV CRAWL_URL https://default.com
ENV ROOT_PATH abc
ENV DISPLAY_URI false
ENV THREAD_COUNT 5
ENV STORE_ON_DISK false
ENTRYPOINT ["go","run","crawler.go"]