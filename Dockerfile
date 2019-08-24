FROM golang:alpine
ADD entrypoint.sh /entrypoint.sh
RUN apk add git
ADD go ./go
RUN cd ./go && go build -o /bin/prwatch prwatch.go
ENTRYPOINT ["/entrypoint.sh"]
