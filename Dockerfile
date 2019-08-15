FROM golang:alpine
RUN apk update && apk add --no-cache git
ADD entrypoint.sh /entrypoint.sh
ADD bin /bin
ENTRYPOINT ["/entrypoint.sh"]
