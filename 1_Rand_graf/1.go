package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func gRGraph(numNodes int, connectivity float64) map[int][]int {
	rand.NewSource(time.Now().UnixNano())
	graph := make(map[int][]int)
	for i := 0; i < numNodes; i++ {
		for j := i + 1; j < numNodes; j++ {
			if rand.Float64() < connectivity {
				graph[i] = append(graph[i], j)
				graph[j] = append(graph[j], i)
			}
		}
	}
	return graph
}

func saveGrCSV(graph map[int][]int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for node, neighbors := range graph {
		record := []string{fmt.Sprintf("%d", node)}
		for _, neighbor := range neighbors {
			record = append(record, fmt.Sprintf("%d", neighbor))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	numNodes := 10      // кол-во узлов
	connectivity := 0.3 // вероятность

	graph := gRGraph(numNodes, connectivity)

	if err := saveGrCSV(graph, "rand_graf.csv"); err != nil {
		fmt.Println("ошибка сохранения в CSV:", err)
	}

	for node, neighbors := range graph {
		fmt.Printf("%d: %v\n", node, neighbors)
	}
}
