
# Base Dockerfile
FROM golang:1.23.5 AS base

# Dev Dockerfile
FROM base AS development

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

RUN go mod download

CMD ["air"]

# Builder Dockerfile
FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bechdelproxy

# Production docker
FROM scratch AS production

LABEL org.opencontainers.image.source=https://github.com/arthur-trt/bechdelproxy
LABEL org.opencontainers.image.description="Bechdel Proxy"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /app

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/bechdelproxy ./

CMD ["/app/bechdelproxy"]
