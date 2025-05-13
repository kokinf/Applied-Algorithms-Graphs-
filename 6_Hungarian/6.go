package main

import (
	"fmt"
	"math/rand"
	"time"

	hungarian "github.com/oddg/hungarian-algorithm"
)

func generateRandomGraph(uSize, vSize, edgeProbability int) [][]int {
	rand.NewSource(time.Now().UnixNano())

	graph := make([][]int, uSize)
	for i := range graph {
		graph[i] = make([]int, vSize)
		for j := range graph[i] {
			if rand.Intn(100) < edgeProbability {
				graph[i][j] = 1
			} else {
				graph[i][j] = 0
			}
		}
	}

	return graph
}

func main() {
	const uSize = 2000
	const vSize = 2000
	const edgeProbability = 30

	graph := generateRandomGraph(uSize, vSize, edgeProbability)

	costMatrix := make([][]int, uSize)
	for i := range graph {
		costMatrix[i] = make([]int, vSize)
		for j := range graph[i] {
			if graph[i][j] == 1 {
				costMatrix[i][j] = 1
			} else {
				costMatrix[i][j] = 0
			}
		}
	}

	// Венгерский алгоритм для поиска максимального паросочетания
	result, _ := hungarian.Solve(costMatrix)
	// индекс = узлу из доли U, значение = узлу из доли V
	fmt.Println("Максимальное паросочетание:")
	for u, v := range result {
		if v != -1 {
			fmt.Printf("U[%d] -> V[%d]\n", u, v)
		}
	}

}
