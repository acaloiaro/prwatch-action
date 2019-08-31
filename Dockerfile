FROM golang:alpine

LABEL com.github.actions.name="PRwatch"
LABEL com.github.actions.description="Let Github monitor Pull Requests and manage development processes for you."
LABEL com.github.actions.icon="watch"
LABEL com.github.actions.color="green"

LABEL repository="https://github.com/acaloiaro/prwatch"
LABEL homepage="https://github.com/acaloiaro/prwatch"
LABEL maintainer="Adriano Caloiaro <adriano@caloiaro.com>"

ADD entrypoint.sh /entrypoint.sh
RUN apk add git
ADD go ./go
RUN cd ./go && go build -o /bin/prwatch prwatch.go
ENTRYPOINT ["/entrypoint.sh"]
