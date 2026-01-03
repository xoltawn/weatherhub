FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /weather-api ./cmd/server/main.go

# STAGE 2
FROM scratch

COPY --from=builder /weather-api /weather-api

EXPOSE 8080

ENTRYPOINT ["/weather-api"]