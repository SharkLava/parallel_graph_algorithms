package algorithms

import (
	"parallel_graph_algorithms/internal/graph"
	"math"
	"runtime"
	"sync"
)

const INF = math.MaxInt32

func RegularFloydWarshall(g *graph.Graph) [][]int {
	dist := make([][]int, g.Vertices)
	for i := range dist {
		dist[i] = make([]int, g.Vertices)
		copy(dist[i], g.AdjMatrix[i])
		for j := range dist[i] {
			if i != j && dist[i][j] == 0 {
				dist[i][j] = INF
			}
		}
	}

	numCPU := runtime.NumCPU()
	chunkSize := g.Vertices / numCPU
	if chunkSize == 0 {
		chunkSize = 1
	}

	for k := 0; k < g.Vertices; k++ {
		var wg sync.WaitGroup
		for chunk := 0; chunk < g.Vertices; chunk += chunkSize {
			wg.Add(1)
			go func(start, end, k int) {
				defer wg.Done()
				for i := start; i < end && i < g.Vertices; i++ {
					for j := 0; j < g.Vertices; j++ {
						if dist[i][k] != INF && dist[k][j] != INF {
							newDist := dist[i][k] + dist[k][j]
							if newDist < dist[i][j] {
								dist[i][j] = newDist
							}
						}
					}
				}
			}(chunk, chunk+chunkSize, k)
		}
		wg.Wait()
	}

	return dist
}

func MatrixFloydWarshall(g *graph.Graph) [][]int {
	dist := make([][]int, g.Vertices)
	for i := range dist {
		dist[i] = make([]int, g.Vertices)
		copy(dist[i], g.AdjMatrix[i])
		for j := range dist[i] {
			if i != j && dist[i][j] == 0 {
				dist[i][j] = INF
			}
		}
	}

	numCPU := runtime.NumCPU()
	chunkSize := g.Vertices / numCPU
	if chunkSize == 0 {
		chunkSize = 1
	}

	for k := 0; k < g.Vertices; k++ {
		kDist := dist[k]
		var wg sync.WaitGroup
		for chunk := 0; chunk < g.Vertices; chunk += chunkSize {
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end && i < g.Vertices; i++ {
					iDist := dist[i]
					ikDist := iDist[k]
					if ikDist == INF {
						continue
					}
					for j := 0; j < g.Vertices; j++ {
						kjDist := kDist[j]
						if kjDist != INF {
							newDist := ikDist + kjDist
							if newDist < iDist[j] {
								iDist[j] = newDist
							}
						}
					}
				}
			}(chunk, chunk+chunkSize)
		}
		wg.Wait()
	}

	return dist
}
