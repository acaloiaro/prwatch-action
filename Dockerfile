FROM golang:alpine
ADD entrypoint.sh /entrypoint.sh
ADD bin /bin
ENTRYPOINT ["/entrypoint.sh"]
