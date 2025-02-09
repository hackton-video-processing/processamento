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

# Etapa 2: Adicionar ffmpeg
FROM debian:bullseye-slim AS runtime

# Instalar ffmpeg
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*

# Etapa 3: Imagem final
FROM scratch

WORKDIR /

# Copiar binário compilado e ffmpeg
COPY --from=builder /main .
COPY --from=runtime /usr/bin/ffmpeg /usr/bin/ffmpeg

# Definir o ponto de entrada
ENTRYPOINT ["/main"]
