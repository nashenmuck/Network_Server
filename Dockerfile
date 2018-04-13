FROM golang:alpine as builder
RUN apk add --no-cache git
COPY . /go/src/github.com/nashenmuck/network_server
RUN go get github.com/nashenmuck/network_server
FROM alpine
RUN apk add --no-cache ca-certificates
RUN adduser -S network
USER network
COPY --from=builder /go/bin/network_server /network/bin/
COPY sql /network/bin/sql
WORKDIR /network/bin
ENTRYPOINT ["./network_server"]
