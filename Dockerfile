# Etapa 1: Build
FROM golang:1.23.4 AS builder

WORKDIR /app

# Copiar e instalar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o binário de forma estática
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Etapa 2: Adicionar ffmpeg e certificados
FROM debian:bullseye-slim AS runtime

# Instalar ffmpeg e atualizar certificados
RUN apt-get update && apt-get install -y ffmpeg ca-certificates && rm -rf /var/lib/apt/lists/*

# Etapa 3: Imagem final
FROM scratch

WORKDIR /

# Copiar binário, ffmpeg e certificados
COPY --from=builder /main .
COPY --from=runtime /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=runtime /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Definir o ponto de entrada
ENTRYPOINT ["/main"]
