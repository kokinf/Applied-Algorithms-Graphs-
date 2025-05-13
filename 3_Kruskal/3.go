package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type Edge struct {
	Start  string
	End    string
	Weight int
}

type DisjointSet struct {
	parent map[string]string
	rank   map[string]int
}

// NewDisjointSet создает новую систему непересекающихся множеств
func NewDisjointSet() *DisjointSet {
	return &DisjointSet{
		parent: make(map[string]string),
		rank:   make(map[string]int),
	}
}

// Find находит корень множества для элемента
func (ds *DisjointSet) Find(node string) string {
	if ds.parent[node] != node {
		ds.parent[node] = ds.Find(ds.parent[node])
	}
	return ds.parent[node]
}

// Union объединяет два множества
func (ds *DisjointSet) Union(node1, node2 string) {
	root1 := ds.Find(node1)
	root2 := ds.Find(node2)

	if root1 != root2 {
		if ds.rank[root1] > ds.rank[root2] {
			ds.parent[root2] = root1
		} else if ds.rank[root1] < ds.rank[root2] {
			ds.parent[root1] = root2
		} else {
			ds.parent[root2] = root1
			ds.rank[root1]++
		}
	}
}

// Kruskal реализует алгоритм Крускала для нахождения MST
func Kruskal(edges []Edge, vertices []string) []Edge {
	ds := NewDisjointSet()
	for _, v := range vertices {
		ds.parent[v] = v
		ds.rank[v] = 0
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	mst := []Edge{}
	for _, edge := range edges {
		root1 := ds.Find(edge.Start)
		root2 := ds.Find(edge.End)

		if root1 != root2 {
			mst = append(mst, edge)
			ds.Union(edge.Start, edge.End)
		}
	}

	return mst
}

func ReadGraph(filename string) ([]Edge, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	edges := []Edge{}
	vertices := make(map[string]bool)

	for _, record := range records {
		if len(record) != 3 {
			return nil, nil, fmt.Errorf("неверный формат строки: %v", record)
		}

		start := record[0] // Начальный узел
		end := record[1]   // Конечный узел
		weight, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, nil, fmt.Errorf("неверный вес ребра: %v", record[2])
		}

		edges = append(edges, Edge{Start: start, End: end, Weight: weight})
		vertices[start] = true
		vertices[end] = true
	}

	vertexList := []string{}
	for v := range vertices {
		vertexList = append(vertexList, v)
	}

	return edges, vertexList, nil
}

func WriteGraph(filename string, edges []Edge) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, edge := range edges {
		record := []string{edge.Start, edge.End, strconv.Itoa(edge.Weight)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func GenerateRandomGraph(numVertices, numEdges int) ([]Edge, []string) {
	rand.NewSource(time.Now().UnixNano())

	vertices := make([]string, numVertices)
	for i := 0; i < numVertices; i++ {
		vertices[i] = fmt.Sprintf("V%d", i+1)
	}

	edges := []Edge{}
	edgeMap := make(map[string]bool)

	for len(edges) < numEdges {
		start := vertices[rand.Intn(numVertices)]
		end := vertices[rand.Intn(numVertices)]
		if start == end {
			continue
		}

		// Убедимся, ребро не дублируется
		edgeKey1 := fmt.Sprintf("%s-%s", start, end)
		edgeKey2 := fmt.Sprintf("%s-%s", end, start)
		if edgeMap[edgeKey1] || edgeMap[edgeKey2] {
			continue
		}

		weight := rand.Intn(100) + 1
		edges = append(edges, Edge{Start: start, End: end, Weight: weight})
		edgeMap[edgeKey1] = true
		edgeMap[edgeKey2] = true
	}

	return edges, vertices
}

func main() {
	numVertices := 9 // Количество вершин
	numEdges := 13   // Количество рёбер
	edges, vertices := GenerateRandomGraph(numVertices, numEdges)

	err := WriteGraph("input.csv", edges)
	if err != nil {
		fmt.Println("Ошибка при записи графа в файл:", err)
		return
	}
	fmt.Println("Случайный граф записан в input.csv")

	edges, vertices, err = ReadGraph("input.csv")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	mst := Kruskal(edges, vertices)

	err = WriteGraph("output.csv", mst)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Println("MST записано в output.csv")
}
