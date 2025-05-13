package main

import (
	"os"
	"testing"
)

func TestGenerateRandomGraph(t *testing.T) {
	numVertices := 5
	numEdges := 7
	edges, vertices := GenerateRandomGraph(numVertices, numEdges)

	if len(vertices) != numVertices {
		t.Errorf("Ожидалось %d вершин, получено %d", numVertices, len(vertices))
	}

	if len(edges) != numEdges {
		t.Errorf("Ожидалось %d рёбер, получено %d", numEdges, len(edges))
	}

	edgeMap := make(map[string]bool)
	for _, edge := range edges {
		key := edge.Start + "-" + edge.End
		if edgeMap[key] {
			t.Errorf("Обнаружено дублирующееся ребро: %s", key)
		}
		edgeMap[key] = true
	}
}

func TestWriteAndReadGraph(t *testing.T) {
	edges := []Edge{
		{Start: "V1", End: "V2", Weight: 10},
		{Start: "V2", End: "V3", Weight: 20},
		{Start: "V3", End: "V4", Weight: 30},
	}
	vertices := []string{"V1", "V2", "V3", "V4"}

	testFile := "test_graph.csv"
	err := WriteGraph(testFile, edges)
	if err != nil {
		t.Fatalf("Ошибка при записи графа в файл: %v", err)
	}
	defer os.Remove(testFile)

	readEdges, readVertices, err := ReadGraph(testFile)
	if err != nil {
		t.Fatalf("Ошибка при чтении графа из файла: %v", err)
	}

	if len(readEdges) != len(edges) {
		t.Errorf("Ожидалось %d рёбер, получено %d", len(edges), len(readEdges))
	}

	for i, edge := range edges {
		if readEdges[i].Start != edge.Start || readEdges[i].End != edge.End || readEdges[i].Weight != edge.Weight {
			t.Errorf("Ожидалось ребро %v, получено %v", edge, readEdges[i])
		}
	}

	if len(readVertices) != len(vertices) {
		t.Errorf("Ожидалось %d вершин, получено %d", len(vertices), len(readVertices))
	}
}

// TestKruskal проверяет корректность работы алгоритма Крускала.
func TestKruskal(t *testing.T) {
	// Создаем тестовый граф
	edges := []Edge{
		{Start: "V1", End: "V2", Weight: 10},
		{Start: "V2", End: "V3", Weight: 20},
		{Start: "V3", End: "V4", Weight: 30},
		{Start: "V1", End: "V4", Weight: 40},
	}
	vertices := []string{"V1", "V2", "V3", "V4"}

	// Ожидаемое минимальное остовное дерево
	expectedMST := []Edge{
		{Start: "V1", End: "V2", Weight: 10},
		{Start: "V2", End: "V3", Weight: 20},
		{Start: "V3", End: "V4", Weight: 30},
	}

	mst := Kruskal(edges, vertices)

	if len(mst) != len(expectedMST) {
		t.Errorf("Ожидалось %d рёбер в MST, получено %d", len(expectedMST), len(mst))
	}

	for i, edge := range expectedMST {
		if mst[i].Start != edge.Start || mst[i].End != edge.End || mst[i].Weight != edge.Weight {
			t.Errorf("Ожидалось ребро %v, получено %v", edge, mst[i])
		}
	}
}

// TestKruskalWithDisconnectedGraph проверяет работу алгоритма Крускала на несвязном графе.
func TestKruskalWithDisconnectedGraph(t *testing.T) {
	// Создаем тестовый граф с двумя компонентами связности
	edges := []Edge{
		{Start: "V1", End: "V2", Weight: 10},
		{Start: "V3", End: "V4", Weight: 20},
	}
	vertices := []string{"V1", "V2", "V3", "V4"}

	// Ожидаемое минимальное остовное дерево
	expectedMST := []Edge{
		{Start: "V1", End: "V2", Weight: 10},
		{Start: "V3", End: "V4", Weight: 20},
	}

	mst := Kruskal(edges, vertices)

	// Проверяем, что результат совпадает с ожидаемым
	if len(mst) != len(expectedMST) {
		t.Errorf("Ожидалось %d рёбер в MST, получено %d", len(expectedMST), len(mst))
	}

	for i, edge := range expectedMST {
		if mst[i].Start != edge.Start || mst[i].End != edge.End || mst[i].Weight != edge.Weight {
			t.Errorf("Ожидалось ребро %v, получено %v", edge, mst[i])
		}
	}
}
