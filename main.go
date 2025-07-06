package main

import (
	"PROJECTFINAL/api"
	"PROJECTFINAL/cRedis"

	"github.com/gin-gonic/gin"
)

func main() {
	cRedis.IniciarRedis()

	router := gin.Default()

	// 🛡️ Middleware CORS para que Angular pueda acceder a la API
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	api.HostIP = "nodo1" // cambia esto según el nodo

	router.POST("/train", api.TrainHandler)
	router.POST("/predict", api.PredictHandler)

	router.Run(":8080")
}
