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
