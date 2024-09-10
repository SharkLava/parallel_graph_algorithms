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

func (g *Graph) ToDenseMatrix() [][]float64 {
	dense := make([][]float64, g.Vertices)
	for i := range dense {
		dense[i] = make([]float64, g.Vertices)
		for j := range dense[i] {
			dense[i][j] = float64(g.AdjMatrix[i][j])
		}
	}
	return dense
}

func (g *Graph) AddEdge(v, w, weight int) {
	g.AdjList[v] = append(g.AdjList[v], Edge{To: w, Weight: weight})
	g.AdjList[w] = append(g.AdjList[w], Edge{To: v, Weight: weight})
	g.AdjMatrix[v][w] = weight
	g.AdjMatrix[w][v] = weight
}
