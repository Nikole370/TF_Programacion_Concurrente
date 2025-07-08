package Modelo

import (
	"PROJECTFINAL/cRedis"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

// ----------- Funciones comunes -----------

func sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

// predict clasifica el ejemplo X y devuelve una etiqueta en texto.
func predict(X []float64, weights []float64) float64 {
	var z float64
	for i := 0; i < len(X); i++ {
		z += X[i] * weights[i]
	}
	return sigmoid(z)
}

func loadCSVData(path string, partitionIdx, totalPartitions int) ([][]float64, []float64, float64, float64, float64, float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, 0, 0, 0, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, 0, 0, 0, 0, err
	}

	totalRows := len(records) - 1
	rowsPerPartition := totalRows / totalPartitions
	start := partitionIdx * rowsPerPartition
	end := start + rowsPerPartition
	if partitionIdx == totalPartitions-1 {
		end = totalRows
	}

	var X [][]float64
	var y []float64
	minEstrato, maxEstrato := math.MaxFloat64, -math.MaxFloat64
	minDominio, maxDominio := math.MaxFloat64, -math.MaxFloat64

	for i := start + 1; i < end+1 && i < len(records); i++ {
		row := records[i]
		dominio, err1 := strconv.ParseFloat(row[6], 64)
		estrato, err2 := strconv.ParseFloat(row[7], 64)
		labelStr := row[9]
		if err1 != nil || err2 != nil || (labelStr != "1" && labelStr != "2") {
			continue
		}

		if dominio < minDominio {
			minDominio = dominio
		}
		if dominio > maxDominio {
			maxDominio = dominio
		}
		if estrato < minEstrato {
			minEstrato = estrato
		}
		if estrato > maxEstrato {
			maxEstrato = estrato
		}

		label := 0.0
		if labelStr == "1" {
			label = 1.0
		}

		xi := []float64{1.0, estrato, dominio}
		X = append(X, xi)
		y = append(y, label)
	}

	return X, y, minEstrato, maxEstrato, minDominio, maxDominio, nil
}

func normalizeFeatures(X [][]float64, minEstrato, maxEstrato, minDominio, maxDominio float64) {
	for i := 0; i < len(X); i++ {
		X[i][1] = (X[i][1] - minEstrato) / (maxEstrato - minEstrato)
		X[i][2] = (X[i][2] - minDominio) / (maxDominio - minDominio)
	}
}

// ----------- Entrenamiento Concurrente -----------

func trainConcurrent(X [][]float64, y []float64, learningRate float64, iterations int, batchSize int) []float64 {
	features := len(X[0])
	weights := make([]float64, features)
	dataLen := len(X)

	for iter := 0; iter < iterations; iter++ {
		var wg sync.WaitGroup
		var mutex sync.Mutex

		for i := 0; i < dataLen; i += batchSize {
			wg.Add(1)

			start := i
			end := i + batchSize
			if end > dataLen {
				end = dataLen
			}

			go func(start, end int) {
				defer wg.Done()
				partialGradients := make([]float64, features)

				for j := start; j < end; j++ {
					pred := predict(X[j], weights)
					error := pred - y[j]
					for k := 0; k < features; k++ {
						partialGradients[k] += error * X[j][k]
					}
				}

				mutex.Lock()
				for k := 0; k < features; k++ {
					weights[k] -= learningRate * partialGradients[k] / float64(end-start)
				}
				mutex.Unlock()
			}(start, end)
		}

		wg.Wait()
	}
	return weights
}

func calculateAccuracy(X [][]float64, y []float64, weights []float64) float64 {
	correct := 0
	for i := 0; i < len(X); i++ {
		pred := predict(X[i], weights)
		if (pred >= 0.5 && y[i] == 1.0) || (pred < 0.5 && y[i] == 0.0) {
			correct++
		}
	}
	return float64(correct) / float64(len(X)) * 100
}

func EntrenarModelo(csvP string, partitionIdx, totalPartitions int) (
	[]float64, float64, float64, float64, float64, error) {
	X, y, minEstrato, maxEstrato, minDominio, maxDominio, _ := loadCSVData(csvP, partitionIdx, totalPartitions)

	normalizeFeatures(X, minEstrato, maxEstrato, minDominio, maxDominio)

	learningRate := 0.1
	iterations := 1000
	batchSize := 100

	pesos := trainConcurrent(X, y, learningRate, iterations, batchSize)
	accuracy := calculateAccuracy(X, y, pesos)

	fmt.Printf("Precisión (nodo %d): %.2f%%\n", partitionIdx, accuracy)

	cRedis.GuardarMinMax("minEstrato", minEstrato)
	cRedis.GuardarMinMax("maxEstrato", maxEstrato)
	cRedis.GuardarMinMax("minDominio", minDominio)
	cRedis.GuardarMinMax("maxDominio", maxDominio)

	return pesos, minEstrato, maxEstrato, minDominio, maxDominio, nil
}

func Predecir(estrato, dominio float64, pesos []float64, minEstrato, maxEstrato, minDominio, maxDominio float64) float64 {
	estratoNorm := (estrato - minEstrato) / (maxEstrato - minEstrato)
	dominioNorm := (dominio - minDominio) / (maxDominio - minDominio)

	input := []float64{1.0, estratoNorm, dominioNorm}
	return predict(input, pesos)
}
