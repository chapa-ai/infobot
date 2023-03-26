FROM golang:1.17-alpine as build-stage

RUN mkdir -p /app

WORKDIR /app

COPY . /app
RUN go mod download

RUN go build -o info cmd/main.go

ENTRYPOINT [ "/app/info" ]