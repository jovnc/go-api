FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
