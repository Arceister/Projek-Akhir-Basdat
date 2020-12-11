FROM golang:1.15-alpine

ENV MYSQL_DATABASE Mihaya
ENV MYSQL_USER root

WORKDIR /app

RUN apk update \
&& apk add mysql mysql-client \
&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app
COPY . .
RUN go get ./... 
RUN go build 

EXPOSE 14045
EXPOSE 3306
CMD [ "go", "run", "*.go" ]