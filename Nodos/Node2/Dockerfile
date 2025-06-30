FROM golang:alpine

WORKDIR /app
# 3) Copia el go.mod y go.sum para que Go reconozca tu módulo
COPY go.mod go.sum ./

# 4) Descarga dependencias
RUN go mod download

COPY datos.csv ./datos.csv

COPY Nodos/Node2/ConexionP2P.go main.go
COPY Modelo ./Modelo
COPY cRedis ./cRedis
EXPOSE 9002
CMD ["go", "run", "main.go"]
