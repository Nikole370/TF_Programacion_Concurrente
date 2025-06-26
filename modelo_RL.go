package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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

func calcularGradiente(dataset []Registro, pesos []float64, out chan<- []float64, wg *sync.WaitGroup) {
	defer wg.Done()
	grad := make([]float64, len(pesos))

	for _, r := range dataset {
		x := []float64{1.0, float64(r.Estrato), float64(r.Dominio)}
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

func cargarDatosCSV(path string) ([]Registro, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var dataset []Registro
	for i, row := range records {
		if i == 0 {
			continue // saltar cabecera
		}
		dominio, _ := strconv.Atoi(row[6])
		estrato, _ := strconv.Atoi(row[7])
		p612 := row[9]
		if p612 != "1" && p612 != "2" {
			continue
		}
		label := 1
		if p612 == "2" {
			label = 0
		}
		dataset = append(dataset, Registro{Estrato: estrato, Dominio: dominio, UsaTIC: label})
	}
	return dataset, nil
}

func main() {
	dataset, err := cargarDatosCSV("Enaho01-2022-612.csv")
	if err != nil {
		log.Fatal("Error cargando CSV:", err)
	}
	fmt.Println("Total registros válidos:", len(dataset))

	pesos := []float64{0.0, 0.0, 0.0}
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

		for grad := range out {
			for i := range pesos {
				pesos[i] -= learningRate * grad[i] / float64(len(dataset))
			}
		}
	}

	fmt.Println("Pesos entrenados:")
	fmt.Printf("w0: %.4f, w1: %.4f, w2: %.4f\n", pesos[0], pesos[1], pesos[2])
}
