FROM golang:alpine
ADD entrypoint.sh /entrypoint.sh
RUN apk add git
ADD bin /bin
ENTRYPOINT ["/entrypoint.sh"]
