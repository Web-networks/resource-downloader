FROM golang:alpine as build
COPY . /usr/src/resource-downloader
WORKDIR /usr/src/resource-downloader/src
ENV GO111MODULE=on
RUN go build -o resource-downloader

FROM alpine:latest as runtime
COPY --from=build /usr/src/resource-downloader/src/resource-downloader resource-downloader
CMD ["./resource-downloader"]