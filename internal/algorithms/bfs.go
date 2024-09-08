package algorithms

import (
	"parallel_graph_algorithms/internal/graph"
	"runtime"
	"sync"
)

func RegularBFS(g *graph.Graph, start int) []int {
	visited := make([]bool, g.Vertices)
	result := make([]int, 0, g.Vertices)
	queue := make([]int, 0, g.Vertices)

	visited[start] = true
	queue = append(queue, start)
	result = append(result, start)

	for len(queue) > 0 {
		levelSize := len(queue)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < levelSize; i++ {
			v := queue[i]
			wg.Add(1)
			go func(v int) {
				defer wg.Done()
				localQueue := make([]int, 0, len(g.AdjList[v]))
				for _, edge := range g.AdjList[v] {
					if !visited[edge.To] {
						visited[edge.To] = true
						localQueue = append(localQueue, edge.To)
					}
				}
				if len(localQueue) > 0 {
					mu.Lock()
					result = append(result, localQueue...)
					queue = append(queue, localQueue...)
					mu.Unlock()
				}
			}(v)
		}
		wg.Wait()
		queue = queue[levelSize:]
	}

	return result
}

func MatrixBFS(g *graph.Graph, start int) []int {
	visited := make([]bool, g.Vertices)
	result := make([]int, 0, g.Vertices)
	currentLevel := make([]int, g.Vertices)

	visited[start] = true
	currentLevel[start] = 1
	result = append(result, start)

	for len(result) < g.Vertices {
		nextLevel := make([]int, g.Vertices)
		var wg sync.WaitGroup
		var mu sync.Mutex

		chunkSize := g.Vertices / runtime.NumCPU()
		if chunkSize == 0 {
			chunkSize = 1
		}

		for chunk := 0; chunk < g.Vertices; chunk += chunkSize {
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				localResult := make([]int, 0)
				for v := start; v < end && v < g.Vertices; v++ {
					if currentLevel[v] > 0 {
						for u := 0; u < g.Vertices; u++ {
							if g.AdjMatrix[v][u] != 0 && !visited[u] {
								nextLevel[u] = 1
								visited[u] = true
								localResult = append(localResult, u)
							}
						}
					}
				}
				if len(localResult) > 0 {
					mu.Lock()
					result = append(result, localResult...)
					mu.Unlock()
				}
			}(chunk, chunk+chunkSize)
		}
		wg.Wait()
		currentLevel = nextLevel
	}

	return result
}
