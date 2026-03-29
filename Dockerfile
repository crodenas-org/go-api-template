# Stage 1 — build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

# Stage 2 — run
FROM alpine:3.21

# ca-certificates: required for TLS to Azure AD OIDC endpoints and Postgres
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /bin/api /app/api

EXPOSE 8080
ENTRYPOINT ["/app/api"]
