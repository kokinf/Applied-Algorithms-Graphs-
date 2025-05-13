package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	INF = 999999 // Отсутствия пути
)

func generateRandomGraph(vertices int, maxWeight int, density float64) [][]int {
	rand.NewSource(time.Now().UnixNano())
	g := make([][]int, vertices)
	for i := range g {
		g[i] = make([]int, vertices)
		for j := range g[i] {
			if i == j {
				g[i][j] = 0
			} else {
				g[i][j] = INF
			}
		}
	}

	maxEdges := vertices * (vertices - 1)
	edges := int(float64(maxEdges) * density)

	for e := 0; e < edges; e++ {
		i := rand.Intn(vertices)
		j := rand.Intn(vertices)
		if i != j && g[i][j] == INF {
			weight := rand.Intn(maxWeight * 2)
			g[i][j] = weight
		} else {

		}
	}

	return g
}

func saveGraphToCSV(graph [][]int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	vertices := len(graph)

	if err := writer.Write([]string{strconv.Itoa(vertices)}); err != nil {
		return err
	}

	for i := 0; i < vertices; i++ {
		for j := 0; j < vertices; j++ {
			if i != j && graph[i][j] != INF {
				record := []string{
					strconv.Itoa(i),
					strconv.Itoa(j),
					strconv.Itoa(graph[i][j]),
				}
				if err := writer.Write(record); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func floydWarshall(graph [][]int) ([][]int, bool) {
	vertices := len(graph)

	dist := make([][]int, vertices)
	for i := range dist {
		dist[i] = make([]int, vertices)
		copy(dist[i], graph[i])
	}

	// Перебираем все возможные промежуточные вершины
	for k := 0; k < vertices; k++ {
		for i := 0; i < vertices; i++ {
			for j := 0; j < vertices; j++ {
				if dist[i][k] != INF && dist[k][j] != INF {
					if dist[i][k]+dist[k][j] < dist[i][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
					}
				}
			}
		}
	}

	hasNegativeCycle := false
	for i := 0; i < vertices; i++ {
		if dist[i][i] < 0 {
			hasNegativeCycle = true
			break
		}
	}

	return dist, hasNegativeCycle
}

func printDistances(dist [][]int) {
	vertices := len(dist)
	fmt.Println("Матрица кратчайших расстояний:")
	for i := 0; i < vertices; i++ {
		for j := 0; j < vertices; j++ {
			if dist[i][j] == INF {
				fmt.Print("INF ")
			} else {
				fmt.Printf("%3d ", dist[i][j])
			}
		}
		fmt.Println()
	}
}

func main() {
	vertices := 6
	maxWeight := 20
	density := 0.4

	graph := generateRandomGraph(vertices, maxWeight, density)

	err := saveGraphToCSV(graph, "graph.csv")
	if err != nil {
		fmt.Println("Ошибка при сохранении графа:", err)
		return
	}

	distances, hasNegativeCycle := floydWarshall(graph)
	if hasNegativeCycle {
		fmt.Println("В графе есть отрицательные циклы!")
	} else {
		fmt.Println("В графе нет отрицательных циклов.")
	}

	printDistances(distances)
}
