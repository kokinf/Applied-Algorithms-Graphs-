package main

import (
	"bytes"
	"os"
	"testing"
)

func compareSlices(a, b []int) bool {
	// Сравниваем длины слайсов
	if len(a) != len(b) {
		return false
	}
	// Сравниваем элементы слайсов
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDFS(t *testing.T) {

	graph := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {6},
		4: {},
		5: {7},
		6: {},
		7: {},
	}

	// Ожидаемый результат DFS обхода
	expected := []int{1, 3, 6, 2, 5, 7, 4}

	output := captureOutput(func() {
		DFS(graph, 1)
	})

	result := parseOutput(output)

	if !compareSlices(result, expected) {
		t.Errorf("DFS обход неверный. Ожидалось: %v, Получено: %v", expected, result)
	}
}

// Тест для BFS
func TestBFS(t *testing.T) {

	graph := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {6},
		4: {},
		5: {7},
		6: {},
		7: {},
	}

	// Ожидаемый результат BFS обхода
	expected := []int{1, 2, 3, 4, 5, 6, 7}

	output := captureOutput(func() {
		BFS(graph, 1)
	})

	result := parseOutput(output)

	if !compareSlices(result, expected) {
		t.Errorf("BFS обход неверный. Ожидалось: %v, Получено: %v", expected, result)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func parseOutput(output string) []int {
	var result []int
	var num int
	for _, char := range output {
		if char >= '0' && char <= '9' {
			num = num*10 + int(char-'0')
		} else if char == ' ' {
			if num != 0 {
				result = append(result, num)
				num = 0
			}
		}
	}
	if num != 0 {
		result = append(result, num)
	}
	return result
}
