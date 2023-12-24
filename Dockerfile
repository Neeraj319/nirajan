FROM golang:alpine

WORKDIR /app

# ENV GOPATH /go

COPY go.mod go.sum ./

RUN go mod download

# RUN go get github.com/cosmtrek/air@latest

COPY . .

CMD ["go", "run", "."]
