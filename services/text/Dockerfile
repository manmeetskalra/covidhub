FROM golang:1.14.8-buster AS builder

WORKDIR /go/src/covid-hub-cmpt474/services/text
COPY . .

ENV MYSQL_HOST=covidhub.cluster-cscg2ncxunhq.us-west-2.rds.amazonaws.com
ENV MYSQL_PORT=3306
ENV MYSQL_USER=admin
ENV MYSQL_PASSWORD=Shopify7224
ENV DEFAULT_NUMBER=12044106266
ENV NEXMO_API_KEY=942ed9ee
ENV NEXMO_API_SECRET=jI8FftZVmwl6VIHw

RUN set -xe; \
    apt-get update; \
    go build -o /go/bin ./

EXPOSE 8002

ENTRYPOINT ["/go/bin/text"]
