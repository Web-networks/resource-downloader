FROM golang:alpine as build
COPY src /go
RUN go build -o resource-downloader

FROM alpine:latest as runtime
COPY --from=build /go/resource-downloader resource-downloader
CMD ["./resource-downloader"]