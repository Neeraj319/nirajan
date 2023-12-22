FROM golang:alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

RUN go get github.com/cosmtrek/air@latest

COPY . .


ENTRYPOINT ["air", "server.go"]
