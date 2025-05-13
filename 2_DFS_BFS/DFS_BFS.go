package main

import (
	"container/list"
	"fmt"
	"time"
)

// Функция (DFS)
func DFS(graph map[int][]int, start int) {
	stack := list.New()
	visited := make(map[int]bool)

	stack.PushBack(start)

	for stack.Len() > 0 {
		element := stack.Back()
		stack.Remove(element)
		node := element.Value.(int)

		if !visited[node] {
			visited[node] = true
			fmt.Printf("%d ", node)

			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					stack.PushBack(neighbor)
				}
			}
		}
	}

	fmt.Println()
}

// Функция (BFS)
func BFS(graph map[int][]int, start int) {
	queue := list.New()
	visited := make(map[int]bool)

	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		node := element.Value.(int)

		fmt.Printf("%d ", node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushBack(neighbor)
			}
		}
	}

	fmt.Println()
}

func main() {
	graph := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {6},
		4: {},
		5: {7},
		6: {},
		7: {},
	}

	// Время DFS
	start := time.Now()
	fmt.Println("DFS обход:")
	DFS(graph, 1)
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения DFS: %s\n", elapsed)

	// Время BFS
	start = time.Now()
	fmt.Println("BFS обход:")
	BFS(graph, 1)
	elapsed = time.Since(start)
	fmt.Printf("Время выполнения BFS: %s\n", elapsed)
}
