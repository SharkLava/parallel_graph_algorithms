package algorithms

import (
	"math"
	"parallel_graph_algorithms/internal/graph"
	"runtime"
	"sync"
)

func RegularBellmanFord(g *graph.Graph, start int) []int {
	distances := make([]int, g.Vertices)
	for i := range distances {
		distances[i] = math.MaxInt32
	}
	distances[start] = 0

	changed := true
	for i := 0; i < g.Vertices-1 && changed; i++ {
		changed = false
		var wg sync.WaitGroup
		var mu sync.Mutex

		for v := range g.AdjList {
			wg.Add(1)
			go func(v int) {
				defer wg.Done()
				localChanged := false
				for _, edge := range g.AdjList[v] {
					if distances[v] != math.MaxInt32 && distances[v]+edge.Weight < distances[edge.To] {
						distances[edge.To] = distances[v] + edge.Weight
						localChanged = true
					}
				}
				if localChanged {
					mu.Lock()
					changed = true
					mu.Unlock()
				}
			}(v)
		}
		wg.Wait()
	}

	return distances
}

func MatrixBellmanFord(g *graph.Graph, start int) []int {
	distances := make([]int, g.Vertices)
	for i := range distances {
		distances[i] = math.MaxInt32
	}
	distances[start] = 0

	changed := true
	for i := 0; i < g.Vertices-1 && changed; i++ {
		changed = false
		newDistances := make([]int, g.Vertices)
		copy(newDistances, distances)

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
				localChanged := false
				for v := start; v < end && v < g.Vertices; v++ {
					for u := 0; u < g.Vertices; u++ {
						if g.AdjMatrix[u][v] != 0 && distances[u] != math.MaxInt32 {
							newDist := distances[u] + g.AdjMatrix[u][v]
							if newDist < newDistances[v] {
								newDistances[v] = newDist
								localChanged = true
							}
						}
					}
				}
				if localChanged {
					mu.Lock()
					changed = true
					mu.Unlock()
				}
			}(chunk, chunk+chunkSize)
		}
		wg.Wait()

		distances = newDistances
	}

	return distances
}
