FROM golang:1.17-alpine

WORKDIR /go-chat

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 8080