# Usa la imagen oficial de Go
FROM golang:1.20-alpine

# Crea el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el código fuente del cliente al contenedor
COPY ./cliente.go .

# Inicializa el módulo de Go si no existe (puedes omitir si ya copiaste go.mod)
RUN go mod init cliente && go mod tidy

# Comando para ejecutar el cliente
CMD ["go", "run", "cliente.go"]
