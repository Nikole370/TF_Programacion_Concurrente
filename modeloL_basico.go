package main

import (
	"fmt"
	"math"
	"sync"
)

type Registro struct {
	Estrato int
	Dominio int
	UsaTIC  int // 1 o 0
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// Calcula el gradiente para un batch de datos
func calcularGradiente(dataset []Registro, pesos []float64, out chan<- []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	grad := make([]float64, len(pesos))

	for _, r := range dataset {
		x := []float64{1.0, float64(r.Estrato), float64(r.Dominio)} // 1, x1, x2
		y := float64(r.UsaTIC)

		z := 0.0
		for i := range pesos {
			z += pesos[i] * x[i]
		}
		pred := sigmoid(z)
		for i := range grad {
			grad[i] += (pred - y) * x[i]
		}
	}
	out <- grad
}

func main() {
	dataset := []Registro{
		{1, 1, 1}, {2, 0, 0}, {3, 1, 1}, {1, 0, 0}, {2, 1, 1},
	}

	pesos := []float64{0.0, 0.0, 0.0} // w0, w1, w2
	learningRate := 0.1
	epochs := 100

	for ep := 0; ep < epochs; ep++ {
		var wg sync.WaitGroup
		out := make(chan []float64, 1)

		wg.Add(1)
		go calcularGradiente(dataset, pesos, out, &wg)

		go func() {
			wg.Wait()
			close(out)
		}()

		// Acumular gradiente y actualizar pesos
		for grad := range out {
			for i := range pesos {
				pesos[i] -= learningRate * grad[i] / float64(len(dataset)) // promedio
			}
		}
	}

	fmt.Println("Pesos entrenados:")
	fmt.Printf("w0: %.4f, w1: %.4f, w2: %.4f\n", pesos[0], pesos[1], pesos[2])
}
