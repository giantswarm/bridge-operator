FROM alpine:3.13.1

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./bridge-operator /bridge-operator

ENTRYPOINT ["/bridge-operator"]
