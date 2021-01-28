FROM golang:1.14.8-buster AS builder

WORKDIR /go/src/covid-hub-cmpt474/services/email
COPY . .

ENV MYSQL_HOST=covidhub.cluster-cscg2ncxunhq.us-west-2.rds.amazonaws.com
ENV MYSQL_PORT=3306
ENV MYSQL_USER=admin
ENV MYSQL_PASSWORD=Shopify7224
ENV DEFAULT_EMAIL=maniyarmanan1996@gmail.com
ENV DEFAULT_EMAIL_PASSWORD=Shopify9700


RUN set -xe; \
    apt-get update; \
    go build -o /go/bin ./

EXPOSE 8001

ENTRYPOINT ["/go/bin/email"]
