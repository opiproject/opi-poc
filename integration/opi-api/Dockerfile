# syntax=docker/dockerfile:1

FROM golang:1.18.4-alpine

ARG OPI_PLUGIN

WORKDIR ~/go/src/opi-api

COPY go.mod ./
COPY go.sum ./
COPY ./ ./
RUN go mod download
RUN go get

RUN go build -o opi-api
RUN go build -o opi-api-plugin ./plugin-opi

CMD [ "KV_PLUGIN="~/go/src/opi-api/opi-api-plugin", "~/go/src/opi-api/opi-api" ]
