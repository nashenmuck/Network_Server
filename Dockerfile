FROM golang as builder
COPY . /go/src/github.com/nashenmuck/network_server
RUN go get github.com/nashenmuck/network_server
FROM ubuntu:xenial
COPY --from=builder /go/bin/network_server /root/
WORKDIR /root
ENTRYPOINT ["./network_server"]
