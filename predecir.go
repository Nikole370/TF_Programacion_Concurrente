package main

import (
	"fmt"
	"math"
)

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func predecir(estrato, dominio int, pesos []float64) float64 {
	z := pesos[0] + pesos[1]*float64(estrato) + pesos[2]*float64(dominio)
	return sigmoid(z)
}

func main() {
	pesos := []float64{-0.1163, -0.2015, -0.0790} // <- tus pesos entrenados

	// ejemplo de predicción con datos nuevos
	test := []struct {
		Estrato int
		Dominio int
	}{
		{1, 1}, {2, 0}, {3, 1},
	}

	for i, r := range test {
		prob := predecir(r.Estrato, r.Dominio, pesos)
		fmt.Printf("Registro %d -> Estrato: %d, Dominio: %d => Probabilidad uso TIC: %.2f\n", i+1, r.Estrato, r.Dominio, prob)
	}
}
