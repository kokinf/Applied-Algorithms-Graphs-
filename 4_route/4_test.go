package main

import (
	"testing"
)

func TestFindLargestConnectedComponent(t *testing.T) {
	graph := map[int][]int{
		0: {1, 2},
		1: {0},
		2: {0, 3},
		3: {2},
		4: {},
		5: {6},
		6: {5},
	}

	expectedComponent := map[int][]int{
		0: {1, 2},
		1: {0},
		2: {0, 3},
		3: {2},
	}

	largestComponent := findLargestConnectedComponent(graph)

	if len(largestComponent) != len(expectedComponent) {
		t.Errorf("Ожидалось %d вершин в компоненте, получено %d", len(expectedComponent), len(largestComponent))
	}

	for vertex, neighbors := range expectedComponent {
		componentNeighbors, ok := largestComponent[vertex]
		if !ok {
			t.Errorf("Вершина %d отсутствует в максимальной компоненте", vertex)
			continue
		}
		if len(neighbors) != len(componentNeighbors) {
			t.Errorf("Для вершины %d ожидалось %d соседей, получено %d", vertex, len(neighbors), len(componentNeighbors))
		}
		for i, neighbor := range neighbors {
			if neighbor != componentNeighbors[i] {
				t.Errorf("Для вершины %d ожидался сосед %d, получен %d", vertex, neighbor, componentNeighbors[i])
			}
		}
	}
}
