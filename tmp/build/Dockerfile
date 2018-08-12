FROM alpine:3.6

RUN adduser -D website-operator
USER website-operator

ADD tmp/_output/bin/website-operator /usr/local/bin/website-operator
