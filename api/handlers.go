package api

import (
	"PROJECTFINAL/Modelo"
	"PROJECTFINAL/cRedis"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PredictRequest struct {
	Estrato float64 `json:"estrato"`
	Dominio float64 `json:"dominio"`
}

// Dirección del nodo actual (puedes ponerlo como var global si quieres pasar desde main)
var HostIP = "nodo1"

func TrainHandler(c *gin.Context) {
	// Enviar comando "ENTRENAR" al nodo local (simula trigger)
	go sendMessage("ENTRENAR", HostIP)

	c.JSON(http.StatusOK, gin.H{
		"message": "Entrenamiento iniciado en nodo local",
	})
}

func PredictHandler(c *gin.Context) {
	var req PredictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	pesos := cRedis.LeerPesos("modelo")
	if len(pesos) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pesos no disponibles"})
		return
	}

	minE := cRedis.LeerMinMax("minEstrato")
	maxE := cRedis.LeerMinMax("maxEstrato")
	minD := cRedis.LeerMinMax("minDominio")
	maxD := cRedis.LeerMinMax("maxDominio")

	result := Modelo.Predecir(req.Estrato, req.Dominio, pesos, minE, maxE, minD, maxD)

	c.JSON(http.StatusOK, gin.H{
		"prediccion": fmt.Sprintf("%.4f", result),
	})
}

// Envía un mensaje a través de la red P2P al nodo local (puede ser "ENTRENAR")
func sendMessage(msg string, ip string) {
	conn, err := net.Dial("tcp", ip+":9002")
	if err != nil {
		fmt.Println("Error enviando mensaje al nodo:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintln(conn, msg)
}
