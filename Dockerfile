FROM golang:alpine
ADD entrypoint.sh /entrypoint.sh
RUN apk add git
ADD go /go
ENTRYPOINT ["/entrypoint.sh"]
