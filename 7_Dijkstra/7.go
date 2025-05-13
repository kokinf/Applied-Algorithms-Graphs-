package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Edge struct {
	to     int
	weight int
}

// Graph представляет граф как список смежности
type Graph struct {
	Vertices int
	Edges    [][]Edge
}

func NewGraph(vertices int) *Graph {
	return &Graph{
		Vertices: vertices,
		Edges:    make([][]Edge, vertices),
	}
}

// AddEdge добавляет направленное ребро в граф
func (g *Graph) AddEdge(from, to, weight int) {
	g.Edges[from] = append(g.Edges[from], Edge{to, weight})
}

func (g *Graph) SaveToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{strconv.Itoa(g.Vertices)}); err != nil {
		return err
	}

	for from, edges := range g.Edges {
		for _, edge := range edges {
			record := []string{
				strconv.Itoa(from),
				strconv.Itoa(edge.to),
				strconv.Itoa(edge.weight),
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func GenerateRandomGraph(vertices int, edgeProbability float64, maxWeight int) *Graph {
	rand.NewSource(time.Now().UnixNano())
	g := NewGraph(vertices)

	for i := 0; i < vertices; i++ {
		for j := 0; j < vertices; j++ {
			if i != j && rand.Float64() < edgeProbability {
				weight := rand.Intn(maxWeight) + 1
				g.AddEdge(i, j, weight)
			}
		}
	}

	return g
}

func (g *Graph) Dijkstra(start int) ([]int, []int) {
	dist := make([]int, g.Vertices)
	prev := make([]int, g.Vertices)
	visited := make([]bool, g.Vertices)

	for i := range dist {
		dist[i] = math.MaxInt32
		prev[i] = -1
	}
	dist[start] = 0

	for {
		u := -1
		minDist := math.MaxInt32
		for v := 0; v < g.Vertices; v++ {
			if !visited[v] && dist[v] < minDist {
				minDist = dist[v]
				u = v
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		for _, edge := range g.Edges[u] {
			v := edge.to
			if !visited[v] && dist[v] > dist[u]+edge.weight {
				dist[v] = dist[u] + edge.weight
				prev[v] = u
			}
		}
	}

	return dist, prev
}

// HasNegativeCycle проверяет наличие отрицательных циклов
func (g *Graph) HasNegativeCycle() bool {
	dist := make([]int, g.Vertices)
	for i := range dist {
		dist[i] = math.MaxInt32
	}
	dist[0] = 0

	for i := 0; i < g.Vertices-1; i++ {
		for u := 0; u < g.Vertices; u++ {
			for _, edge := range g.Edges[u] {
				v := edge.to
				if dist[u] != math.MaxInt32 && dist[v] > dist[u]+edge.weight {
					dist[v] = dist[u] + edge.weight
				}
			}
		}
	}

	for u := 0; u < g.Vertices; u++ {
		for _, edge := range g.Edges[u] {
			v := edge.to
			if dist[u] != math.MaxInt32 && dist[v] > dist[u]+edge.weight {
				return true
			}
		}
	}

	return false
}

// GetPath восстанавливает путь
func GetPath(prev []int, start, end int) []int {
	path := make([]int, 0)
	for at := end; at != -1; at = prev[at] {
		path = append(path, at)
		if at == start {
			break
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func main() {
	vertices := 100
	edgeProbability := 0.3
	maxWeight := 100

	g := GenerateRandomGraph(vertices, edgeProbability, maxWeight)

	err := g.SaveToCSV("graph.csv")
	if err != nil {
		fmt.Println("Ошибка при сохранении графа:", err)
		return
	}

	if g.HasNegativeCycle() {
		fmt.Println("Граф содержит отрицательные циклы!")
		return
	}

	start := 0
	end := vertices - 1

	dist, prev := g.Dijkstra(start)

	if dist[end] == math.MaxInt32 {
		fmt.Printf("Пути от вершины %d до вершины %d не существует\n", start, end)
	} else {
		path := GetPath(prev, start, end)
		fmt.Printf("Кратчайший путь от вершины %d до вершины %d: %v\n", start, end, path)
		fmt.Printf("Длина пути: %d\n", dist[end])
	}
}
