FROM --platform=linux/amd64 golang:1.20-buster

RUN mkdir /app
WORKDIR /app

RUN mkdir logs
COPY ./server /app/server
COPY ./.env /app/.env
COPY ./proto /app/proto
COPY ./config /app/config
COPY ./migrations /app/migrations

RUN apt-get update -y
RUN apt install -y protobuf-compiler