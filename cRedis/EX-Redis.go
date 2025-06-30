package cRedis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

// Inicializa la conexión a Redis
func IniciarRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Nombre del servicio en docker-compose
		Password: "",           // Sin contraseña
		DB:       0,
	})
}

// Guarda los pesos en Redis bajo una clave específica
func GuardarPesos(clave string, pesos []float64) {
	var strBuild strings.Builder

	for i, val := range pesos {
		strBuild.WriteString(fmt.Sprintf("%f", val))

		if i < len(pesos)-1 {
			strBuild.WriteString(",")
		}
	}

	status := rdb.Set(ctx, clave, strBuild.String(), 0).Err()

	if status != nil {
		log.Printf("No se guardo nada: %v", status)
	}
}

// Lee pesos desde Redis y los devuelve como slice de float64
func LeerPesos(clave string) []float64 {
	val, _ := rdb.Get(ctx, clave).Result()

	valores := strings.Split(val, ",")

	var pesos []float64

	for _, v := range valores {
		var num float64
		fmt.Sscanf(v, "%f", &num)
		pesos = append(pesos, num)
	}

	return pesos
}

func GuardarMinMax(key string, val float64) {
	rdb.Set(ctx, key, val, 0)
}

func LeerMinMax(key string) float64 {
	valStr, _ := rdb.Get(ctx, key).Result()
	val, _ := strconv.ParseFloat(valStr, 64)
	return val
}
