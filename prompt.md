# ./cmd/graph_algo/main.go
```go
package main

import (
	"flag"
	"fmt"
	"parallel_graph_algorithms/internal/algorithms"
	"parallel_graph_algorithms/internal/graph"
	"math/rand"
	"time"
)

func main() {
	density := flag.Float64("density", 0.01, "Graph density")
	size := flag.Int("size", 1000, "Number of vertices")
	algorithm := flag.String("algo", "bfs", "Algorithm to run (bfs, bellman-ford, or floyd-warshall)")
	flag.Parse()

	g := graph.NewGraph(*size, *density)
	start := 0

	switch *algorithm {
	case "bfs":
		runBFS(g, start)
	case "bellman-ford":
		runBellmanFord(g, start)
	case "floyd-warshall":
		runFloydWarshall(g)
	default:
		fmt.Println("Invalid algorithm. Please choose 'bfs', 'bellman-ford', or 'floyd-warshall'.")
	}
}

func runBFS(g *graph.Graph, start int) {
	regularStart := time.Now()
	regularResult := algorithms.RegularBFS(g, start)
	regularDuration := time.Since(regularStart)

	matrixStart := time.Now()
	matrixResult := algorithms.MatrixBFS(g, start)
	matrixDuration := time.Since(matrixStart)

	// fmt.Printf("BFS Results:\nRegular: %v (%d nodes)\nMatrix: %v (%d nodes)\n",
	// 	regularDuration, len(regularResult), matrixDuration, len(matrixResult))
	fmt.Printf("BFS Results:\n")
	fmt.Printf("Regular: %v (%.2f ms)\n", regularDuration, float64(regularDuration.Nanoseconds())/1e6)
	fmt.Printf("Matrix:  %v (%.2f ms)\n", matrixDuration, float64(matrixDuration.Nanoseconds())/1e6)
	fmt.Printf("Speedup: %.2fx\n", float64(regularDuration)/float64(matrixDuration))

	// Verify results (compare the first 10 visited nodes)
	fmt.Println("Comparing first 10 visited nodes:")
	for i := 0; i < 10 && i < len(regularResult) && i < len(matrixResult); i++ {
		fmt.Printf("Node %d: Regular = %d, Matrix = %d\n", i, regularResult[i], matrixResult[i])
	}
}

func runBellmanFord(g *graph.Graph, start int) {
	regularStart := time.Now()
	regularResult := algorithms.RegularBellmanFord(g, start)
	regularDuration := time.Since(regularStart)

	matrixStart := time.Now()
	matrixResult := algorithms.MatrixBellmanFord(g, start)
	matrixDuration := time.Since(matrixStart)

	// fmt.Printf("Bellman-Ford Results:\nRegular: %v\nMatrix: %v\n", regularDuration, matrixDuration)
	// fmt.Printf("Bellman-Ford Results:\nRegular: %v (%d nodes)\nMatrix: %v (%d nodes)\n",
	// 	regularDuration, len(regularResult), matrixDuration, len(matrixResult))

	fmt.Printf("Bellman-Ford Results:\n")
	fmt.Printf("Regular: %v (%.2f ms)\n", regularDuration, float64(regularDuration.Nanoseconds())/1e6)
	fmt.Printf("Matrix:  %v (%.2f ms)\n", matrixDuration, float64(matrixDuration.Nanoseconds())/1e6)
	fmt.Printf("Speedup: %.2fx\n", float64(regularDuration)/float64(matrixDuration))
	// Verify results (compare the first 10 distances)
	fmt.Println("Comparing first 10 distances:")
	for i := 0; i < 10 && i < len(regularResult); i++ {
		fmt.Printf("Node %d: Regular = %d, Matrix = %d\n", i, regularResult[i], matrixResult[i])
	}
}

func runFloydWarshall(g *graph.Graph) {
	regularStart := time.Now()
	regularResult := algorithms.RegularFloydWarshall(g)
	regularDuration := time.Since(regularStart)

	matrixStart := time.Now()
	matrixResult := algorithms.MatrixFloydWarshall(g)
	matrixDuration := time.Since(matrixStart)

	fmt.Printf("Floyd-Warshall Results:\n")
	fmt.Printf("Regular: %v (%.2f ms)\n", regularDuration, float64(regularDuration.Nanoseconds())/1e6)
	fmt.Printf("Matrix:  %v (%.2f ms)\n", matrixDuration, float64(matrixDuration.Nanoseconds())/1e6)
	fmt.Printf("Speedup: %.2fx\n", float64(regularDuration)/float64(matrixDuration))

	// Verify results (compare a few random distances)
	fmt.Println("Comparing a few random distances:")
	for i := 0; i < 5; i++ {
		x, y := rand.Intn(g.Vertices), rand.Intn(g.Vertices)
		fmt.Printf("Distance from %d to %d: Regular = %d, Matrix = %d\n",
			x, y, regularResult[x][y], matrixResult[x][y])
	}
}
```

# ./internal/algorithms/bfs.go
```go
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
```

# ./internal/algorithms/floyd_warshall.go
```go
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
```

# ./internal/algorithms/bellman_ford.go
```go
package algorithms

import (
	"parallel_graph_algorithms/internal/graph"
	"math"
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
```

# ./internal/graph/graph.go
```go
package graph

import (
	"math/rand"
)

type Graph struct {
	AdjList   map[int][]Edge
	AdjMatrix [][]int
	Vertices  int
}

type Edge struct {
	To     int
	Weight int
}

func NewGraph(vertices int, density float64) *Graph {
	g := &Graph{
		AdjList:   make(map[int][]Edge, vertices),
		AdjMatrix: make([][]int, vertices),
		Vertices:  vertices,
	}
	for i := range g.AdjMatrix {
		g.AdjMatrix[i] = make([]int, vertices)
	}

	for i := 0; i < vertices; i++ {
		for j := i + 1; j < vertices; j++ {
			if rand.Float64() < density {
				weight := rand.Intn(10) + 1 // Random weight between 1 and 10
				g.AddEdge(i, j, weight)
			}
		}
	}
	return g
}

func (g *Graph) AddEdge(v, w, weight int) {
	g.AdjList[v] = append(g.AdjList[v], Edge{To: w, Weight: weight})
	g.AdjList[w] = append(g.AdjList[w], Edge{To: v, Weight: weight})
	g.AdjMatrix[v][w] = weight
	g.AdjMatrix[w][v] = weight
}
```

