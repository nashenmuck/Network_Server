FROM golang:alpine as builder
RUN apk add --no-cache git
COPY . /go/src/github.com/nashenmuck/network_server
RUN go get github.com/nashenmuck/network_server
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/network_server /root/
COPY sql /root/sql
WORKDIR /root
ENTRYPOINT ["./network_server"]
