# Build stage
FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Definir diretório de trabalho
WORKDIR /app

# Copiar go.mod e go.sum primeiro (para cache das dependências)
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s' \
    -o main ./main.go

# Production stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata libc6-compat

RUN adduser -D -s /bin/sh appuser

WORKDIR  /home/appuser/

COPY --from=builder /app/main .

COPY migrations ./migrations

RUN chmod +x main

USER appuser

# Expor a porta (ajuste conforme necessário)
EXPOSE 5000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5000/health || exit 1
# Comando para executar a aplicação
CMD ["./main"]