FROM golang:1.23 AS builder

WORKDIR /src

COPY go.mod .
RUN go mod download

COPY . .

RUN mkdir -p .build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/ cmd/post-service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /src/build .
COPY --from=builder /src/.env .
COPY --from=builder /src/migrations ./migrations

ENTRYPOINT [ "./main" ]