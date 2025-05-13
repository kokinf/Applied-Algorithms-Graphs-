package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Graph struct {
	AdjacencyList map[int]map[int]bool
}

func NewGraph() *Graph {
	return &Graph{
		AdjacencyList: make(map[int]map[int]bool),
	}
}

func (g *Graph) AddEdge(u, v int) {
	if _, exists := g.AdjacencyList[u]; !exists {
		g.AdjacencyList[u] = make(map[int]bool)
	}
	if _, exists := g.AdjacencyList[v]; !exists {
		g.AdjacencyList[v] = make(map[int]bool)
	}
	g.AdjacencyList[u][v] = true
	g.AdjacencyList[v][u] = true
}

func GenerateKPartiteGraph(partSizes []int, edgeProbability int) *Graph {
	rand.NewSource(time.Now().UnixNano())
	graph := NewGraph()

	start := 0
	partitions := make([][]int, len(partSizes))
	for i, size := range partSizes {
		partitions[i] = make([]int, size)
		for j := 0; j < size; j++ {
			partitions[i][j] = start + j
		}
		start += size
	}

	for i := 0; i < len(partitions); i++ {
		for j := i + 1; j < len(partitions); j++ {
			for _, u := range partitions[i] {
				for _, v := range partitions[j] {
					if rand.Intn(100) < edgeProbability {
						graph.AddEdge(u, v)
					}
				}
			}
		}
	}

	return graph
}

func (g *Graph) SaveToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for u, neighbors := range g.AdjacencyList {
		for v := range neighbors {
			if u < v {
				err := writer.Write([]string{fmt.Sprintf("%d", u), fmt.Sprintf("%d", v)})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Выполнение реберной раскраски графа
func (g *Graph) EdgeColoring() map[[2]int]int {
	colors := make(map[[2]int]int)             // Карта для хранения цветов рёбер
	vertexColors := make(map[int]map[int]bool) // Карта для хранения цветов

	for u := range g.AdjacencyList {
		vertexColors[u] = make(map[int]bool)
	}

	for u, neighbors := range g.AdjacencyList {
		for v := range neighbors {
			if u < v {
				edge := [2]int{u, v}
				availableColors := make(map[int]bool)

				// Проверка цвета соседних рёбер для узла u
				for color := range vertexColors[u] {
					availableColors[color] = true
				}

				// Проверка цвета соседних рёбер для узла v
				for color := range vertexColors[v] {
					availableColors[color] = true
				}

				// Находим минимальный доступный цвет
				color := 0
				for availableColors[color] {
					color++
				}

				// Назначаем цвет ребру
				colors[edge] = color
				vertexColors[u][color] = true
				vertexColors[v][color] = true
			}
		}
	}

	return colors
}

func ValidateEdgeColoring(graph *Graph, colors map[[2]int]int) bool {
	for u, neighbors := range graph.AdjacencyList {
		usedColors := make(map[int]bool)
		for v := range neighbors {
			if u < v {
				edge := [2]int{u, v}
				color := colors[edge]
				if usedColors[color] {
					fmt.Printf("Ошибка: Вершина %d имеет два ребра одного цвета (%d)\n", u, color)
					return false
				}
				usedColors[color] = true
			}
		}
	}
	return true
}

func VisualizeGraph(graph *Graph, colors map[[2]int]int, numNodes int, filename string) error {
	p := plot.New()

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	nodes := make(plotter.XYs, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i].X = rand.Float64() * 10
		nodes[i].Y = rand.Float64() * 10
	}

	nodePlot, err := plotter.NewScatter(nodes)
	if err != nil {
		return err
	}
	p.Add(nodePlot)

	for edge, color := range colors {
		u, v := edge[0], edge[1]
		line := plotter.XYs{
			{X: nodes[u].X, Y: nodes[u].Y},
			{X: nodes[v].X, Y: nodes[v].Y},
		}
		linePlot, err := plotter.NewLine(line)
		if err != nil {
			return err
		}
		linePlot.Color = plotutil.Color(color)
		p.Add(linePlot)
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, filename); err != nil {
		return err
	}

	return nil
}

func main() {
	partSizes := []int{10, 10, 10, 10}
	edgeProbability := 10

	graph := GenerateKPartiteGraph(partSizes, edgeProbability)

	numNodes := 0
	for _, size := range partSizes {
		numNodes += size
	}
	err := graph.SaveToCSV("color-graph.csv")
	if err != nil {
		fmt.Println("Ошибка при сохранении графа:", err)
		return
	}

	colors := graph.EdgeColoring()

	if !ValidateEdgeColoring(graph, colors) {
		fmt.Println("Раскраска некорректна!", err)
		return
	}

	err = VisualizeGraph(graph, colors, numNodes, "graph_vis.png")
	if err != nil {
		fmt.Println("Ошибка при визуализации графа:", err)
		return
	}
}
