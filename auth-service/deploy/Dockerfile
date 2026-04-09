FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main cmd/auth-service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /src/.env .
COPY --from=builder /src/build/main .

ENTRYPOINT [ "./main" ]