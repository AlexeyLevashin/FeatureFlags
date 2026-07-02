FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sv_flags_api ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sv_flags_api .

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./sv_flags_api"]