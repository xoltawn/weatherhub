FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /weather-api ./cmd/server/main.go

# STAGE 2
FROM scratch

# Copy SSL certs for OpenWeatherMap API calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /weather-api /weather-api

EXPOSE 8080

ENTRYPOINT ["/weather-api"]