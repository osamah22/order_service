FROM golang:1.26-alpine

WORKDIR /app

# install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# generate swagger docs
RUN swag init -g ./cmd/server/main.go -o ./docs

# build the service
RUN go build -o app ./cmd/server

EXPOSE 8080

CMD ["./app"]
