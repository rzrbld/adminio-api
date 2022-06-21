FROM golang:1.17-alpine

LABEL maintainer="rzrbld <razblade@gmail.com>"

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org

RUN  \
     apk add --no-cache git && \
     git clone https://github.com/rzrbld/adminio-api && cd adminio-api/src/ && go build main.go && cp main /go/bin/adminio

FROM alpine:3.16

EXPOSE 8080

COPY --from=0 /go/bin/adminio /usr/bin/adminio

RUN  \
     apk add --no-cache ca-certificates 'curl>7.61.0' 'su-exec>=0.2' && \
     echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

CMD ["adminio"]
