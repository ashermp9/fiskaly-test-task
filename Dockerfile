FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o signer ./cmd/signer/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/signer .

EXPOSE 8080

# Set a default config path
ENV CONFIG_PATH /app/config.yaml

CMD ["sh", "-c", "./signer -config ${CONFIG_PATH}"]

