package main

import (
	"fmt"
	"sync"
)

// función para simular el paso de pipeline: cuadrado del número
func square(num int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	out <- num * num
}

// función para simular el reduce: suma de todos los valores
func reducer(in <-chan int, result *int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range in {
		mutex.Lock()
		*result += val
		mutex.Unlock()
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	outChan := make(chan int, len(nums))

	var wgMap sync.WaitGroup
	var wgReduce sync.WaitGroup
	var mutex sync.Mutex
	finalResult := 0

	// lanzar mappers
	for _, n := range nums {
		wgMap.Add(1)
		go square(n, outChan, &wgMap)
	}

	// cerrar canal cuando terminen los mappers
	go func() {
		wgMap.Wait()
		close(outChan)
	}()

	// lanzar reducer
	wgReduce.Add(1)
	go reducer(outChan, &finalResult, &mutex, &wgReduce)

	// esperar que termine el reduce
	wgReduce.Wait()

	fmt.Printf("Resultado final (suma de cuadrados): %d\n", finalResult)
}
