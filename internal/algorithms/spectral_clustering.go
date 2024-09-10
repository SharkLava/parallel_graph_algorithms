package algorithms

import (
	"math"
	"math/rand"
	"parallel_graph_algorithms/internal/graph"
	"sync"
)

// SpectralClustering performs spectral clustering on the graph
func SpectralClustering(g *graph.Graph, k int) []int {
	// Step 1: Compute the Laplacian matrix
	laplacian := computeLaplacian(g)

	// Step 2: Compute the k smallest eigenvectors
	eigenvectors := computeEigenvectors(laplacian, k)

	// Step 3: Perform k-means clustering on the eigenvectors
	return kMeansClustering(eigenvectors, k)
}

// computeLaplacian computes the normalized Laplacian matrix
func computeLaplacian(g *graph.Graph) [][]float64 {
	n := g.Vertices
	laplacian := make([][]float64, n)
	for i := range laplacian {
		laplacian[i] = make([]float64, n)
	}

	// Compute degree matrix and adjacency matrix
	degree := make([]float64, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g.AdjMatrix[i][j] != 0 {
				degree[i] += float64(g.AdjMatrix[i][j])
			}
		}
	}

	// Compute normalized Laplacian
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < n; j++ {
				if i == j {
					laplacian[i][j] = 1
				} else if g.AdjMatrix[i][j] != 0 {
					laplacian[i][j] = -float64(g.AdjMatrix[i][j]) / math.Sqrt(degree[i]*degree[j])
				}
			}
		}(i)
	}
	wg.Wait()

	return laplacian
}

// computeEigenvectors computes the k smallest eigenvectors using power iteration
func computeEigenvectors(matrix [][]float64, k int) [][]float64 {
	n := len(matrix)
	eigenvectors := make([][]float64, k)
	for i := range eigenvectors {
		eigenvectors[i] = make([]float64, n)
		for j := range eigenvectors[i] {
			eigenvectors[i][j] = rand.Float64()
		}
	}

	numIterations := 100
	var wg sync.WaitGroup

	for iter := 0; iter < numIterations; iter++ {
		for i := 0; i < k; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				// Power iteration
				newVector := make([]float64, n)
				for j := 0; j < n; j++ {
					for l := 0; l < n; l++ {
						newVector[j] += matrix[j][l] * eigenvectors[i][l]
					}
				}

				// Gram-Schmidt orthogonalization
				for j := 0; j < i; j++ {
					dot := 0.0
					for l := 0; l < n; l++ {
						dot += newVector[l] * eigenvectors[j][l]
					}
					for l := 0; l < n; l++ {
						newVector[l] -= dot * eigenvectors[j][l]
					}
				}

				// Normalize
				norm := 0.0
				for j := 0; j < n; j++ {
					norm += newVector[j] * newVector[j]
				}
				norm = math.Sqrt(norm)
				for j := 0; j < n; j++ {
					eigenvectors[i][j] = newVector[j] / norm
				}
			}(i)
		}
		wg.Wait()
	}

	return eigenvectors
}

// kMeansClustering performs k-means clustering on the eigenvectors
func kMeansClustering(vectors [][]float64, k int) []int {
	n := len(vectors[0])
	dim := len(vectors)
	centroids := make([][]float64, k)
	for i := range centroids {
		centroids[i] = make([]float64, dim)
		for j := range centroids[i] {
			centroids[i][j] = vectors[j][rand.Intn(n)]
		}
	}

	assignments := make([]int, n)
	numIterations := 100

	for iter := 0; iter < numIterations; iter++ {
		// Assign points to nearest centroid
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				minDist := math.MaxFloat64
				for j := 0; j < k; j++ {
					dist := 0.0
					for l := 0; l < dim; l++ {
						diff := vectors[l][i] - centroids[j][l]
						dist += diff * diff
					}
					if dist < minDist {
						minDist = dist
						assignments[i] = j
					}
				}
			}(i)
		}
		wg.Wait()

		// Update centroids
		newCentroids := make([][]float64, k)
		counts := make([]int, k)
		for i := range newCentroids {
			newCentroids[i] = make([]float64, dim)
		}

		for i := 0; i < n; i++ {
			cluster := assignments[i]
			counts[cluster]++
			for j := 0; j < dim; j++ {
				newCentroids[cluster][j] += vectors[j][i]
			}
		}

		for i := 0; i < k; i++ {
			if counts[i] > 0 {
				for j := 0; j < dim; j++ {
					centroids[i][j] = newCentroids[i][j] / float64(counts[i])
				}
			}
		}
	}

	return assignments
}
