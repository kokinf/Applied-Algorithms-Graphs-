package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func generateRandomGraph(numVertices int, edgeProbability float64) map[int][]int {
	rand.NewSource(time.Now().UnixNano())
	graph := make(map[int][]int)

	for i := 0; i < numVertices; i++ {
		for j := i + 1; j < numVertices; j++ {
			if rand.Float64() < edgeProbability {
				graph[i] = append(graph[i], j)
				graph[j] = append(graph[j], i)
			}
		}
	}

	return graph
}

func writeGraphToCSV(graph map[int][]int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for vertex, neighbors := range graph {
		record := []string{strconv.Itoa(vertex)}
		if len(neighbors) == 0 {
			record = append(record, "")
		} else {
			for _, neighbor := range neighbors {
				record = append(record, strconv.Itoa(neighbor))
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func readGraphFromCSV(filename string) (map[int][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	graph := make(map[int][]int)
	for _, record := range records {
		if len(record) == 0 {
			continue
		}

		vertex, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		graph[vertex] = []int{}
		for _, neighborStr := range record[1:] {
			if neighborStr != "" {
				neighbor, err := strconv.Atoi(neighborStr)
				if err != nil {
					return nil, err
				}
				graph[vertex] = append(graph[vertex], neighbor)
			}
		}
	}

	return graph, nil
}

func DFS(graph map[int][]int, start int, visited map[int]bool) []int {
	stack := []int{start}
	component := []int{}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if !visited[node] {
			visited[node] = true
			component = append(component, node)

			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					stack = append(stack, neighbor)
				}
			}
		}
	}

	return component
}

// Нахождение максимально связанный компонент
func findLargestConnectedComponent(graph map[int][]int) map[int][]int {
	visited := make(map[int]bool)
	var largestComponent map[int][]int
	maxSize := 0

	for vertex := range graph {
		if !visited[vertex] {
			componentVertices := DFS(graph, vertex, visited)
			if len(componentVertices) > maxSize {
				maxSize = len(componentVertices)
				largestComponent = make(map[int][]int)
				for _, v := range componentVertices {
					largestComponent[v] = graph[v]
				}
			}
		}
	}

	return largestComponent
}

func main() {
	numVertices := 10
	edgeProbability := 0.3
	graph := generateRandomGraph(numVertices, edgeProbability)

	if err := writeGraphToCSV(graph, "input.csv"); err != nil {
		fmt.Println("Ошибка при записи графа в файл:", err)
		return
	}

	readGraph, err := readGraphFromCSV("input.csv")
	if err != nil {
		fmt.Println("Ошибка при чтении графа из файла:", err)
		return
	}

	largestComponent := findLargestConnectedComponent(readGraph)

	if err := writeGraphToCSV(largestComponent, "output.csv"); err != nil {
		fmt.Println("Ошибка при записи максимальной связной компоненты в файл:", err)
		return
	}

}
