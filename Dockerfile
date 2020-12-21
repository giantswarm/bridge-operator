FROM alpine:3.12.3

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./bridge-operator /bridge-operator

ENTRYPOINT ["/bridge-operator"]
