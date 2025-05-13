package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Edge struct {
	to   int
	rev  int
	cap  int
	flow int
}

type Dinic struct {
	graph [][]Edge
	level []int
	ptr   []int
	queue []int
}

// NewDinic создает новый экземпляр для n вершин
func NewDinic(n int) *Dinic {
	return &Dinic{
		graph: make([][]Edge, n),
		level: make([]int, n),
		ptr:   make([]int, n),
		queue: make([]int, n),
	}
}

// AddEdge добавляет ориентированное ребро и обратное ребро
func (d *Dinic) AddEdge(from, to, cap int) {
	forward := Edge{to, len(d.graph[to]), cap, 0}
	backward := Edge{from, len(d.graph[from]), 0, 0}
	d.graph[from] = append(d.graph[from], forward)
	d.graph[to] = append(d.graph[to], backward)
}

// BFS строит слоистую сеть и возвращает достижимость стока
func (d *Dinic) BFS(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	d.level[s] = 0
	q := []int{s}

	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range d.graph[v] {
			if e.flow < e.cap && d.level[e.to] == -1 {
				d.level[e.to] = d.level[v] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.level[t] != -1
}

// DFS ищет блокирующий поток
func (d *Dinic) DFS(v, t, minFlow int) int {
	if v == t {
		return minFlow
	}
	for ; d.ptr[v] < len(d.graph[v]); d.ptr[v]++ {
		e := &d.graph[v][d.ptr[v]]
		if e.flow < e.cap && d.level[v]+1 == d.level[e.to] {
			pushed := d.DFS(e.to, t, min(minFlow, e.cap-e.flow))
			if pushed > 0 {
				e.flow += pushed
				d.graph[e.to][e.rev].flow -= pushed
				return pushed
			}
		}
	}
	return 0
}

// MaxFlow вычисляет максимальный поток
func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.BFS(s, t) {
		for i := range d.ptr {
			d.ptr[i] = 0
		}
		for {
			pushed := d.DFS(s, t, math.MaxInt)
			if pushed == 0 {
				break
			}
			flow += pushed
		}
	}
	return flow
}

// generateNetwork создает сеть с гарантированным путем
func generateNetwork(numNodes, numEdges int, filename string) {
	rand.NewSource(time.Now().UnixNano())

	// Создаем директории при необходимости
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ' '

	writer.Write([]string{strconv.Itoa(numNodes), strconv.Itoa(numEdges)})

	source := 0
	sink := numNodes - 1
	writer.Write([]string{strconv.Itoa(source), strconv.Itoa(sink)})

	// Гарантируем путь от истока к стоку
	pathLength := 5 + rand.Intn(5)
	if pathLength > numNodes {
		pathLength = numNodes - 1
	}

	prev := source
	for i := 0; i < pathLength; i++ {
		next := i + 1
		if next >= numNodes {
			next = sink
		}
		capacity := 50 + rand.Intn(50) // Пропускная способность 50-100
		writer.Write([]string{
			strconv.Itoa(prev),
			strconv.Itoa(next),
			strconv.Itoa(capacity),
		})
		prev = next
	}
	numEdges -= pathLength

	for i := 0; i < numEdges; i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		for to == from || (from == source && to == sink) {
			to = rand.Intn(numNodes)
		}
		capacity := 1 + rand.Intn(100) // Пропускная способность 1-100
		writer.Write([]string{
			strconv.Itoa(from),
			strconv.Itoa(to),
			strconv.Itoa(capacity),
		})
	}

	writer.Flush()
}

func readNetwork(filename string) (int, int, int, [][]int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ' '
	reader.FieldsPerRecord = -1

	record, err := reader.Read()
	if err != nil {
		panic(err)
	}
	numNodes, _ := strconv.Atoi(record[0])
	numEdges, _ := strconv.Atoi(record[1])

	record, err = reader.Read()
	if err != nil {
		panic(err)
	}
	source, _ := strconv.Atoi(record[0])
	sink, _ := strconv.Atoi(record[1])

	edges := make([][]int, 0, numEdges)
	for i := 0; i < numEdges; i++ {
		record, err = reader.Read()
		if err != nil {
			break
		}
		from, _ := strconv.Atoi(record[0])
		to, _ := strconv.Atoi(record[1])
		cap, _ := strconv.Atoi(record[2])
		edges = append(edges, []int{from, to, cap})
	}

	return numNodes, source, sink, edges
}

func main() {
	numNodes := 1000  // Узлы
	numEdges := 20000 // Рёбера
	filename := filepath.Join("network.csv")

	generateNetwork(numNodes, numEdges, filename)

	n, source, sink, edges := readNetwork(filename)

	dinic := NewDinic(n)
	for _, edge := range edges {
		dinic.AddEdge(edge[0], edge[1], edge[2])
	}

	if !dinic.BFS(source, sink) {
		fmt.Println("ОШИБКА: нет пути от источника к приемнику!")
		return
	}

	maxFlow := dinic.MaxFlow(source, sink)
	fmt.Printf("Максимальный поток: %d\n", maxFlow)

}
