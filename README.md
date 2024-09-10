# Graph Algorithms in Go

This project implements and compares the performance of various graph algorithms using both regular and matrix-based approaches. The algorithms included are:

- Breadth-First Search (BFS)
- Bellman-Ford
- Floyd-Warshall
- Spectral Clustering

All implementations are optimized for parallel execution to take advantage of multi-core processors.

## Requirements

- Go 1.18 or higher

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/SharkLava/parallel_graph_algorithms.git
   cd parallel_graph_algorithms
   ```

2. Build the project:
   ```
   go build ./cmd/graph_algo
   ```

## Usage

Run the program using the following command:

```
./graph_algo -algo <algorithm> -size <graph_size> -density <graph_density>
```

Where:
- `<algorithm>` is one of: `bfs`, `bellman-ford`, `spectral-clustering`, or `floyd-warshall`
- `<graph_size>` is the number of vertices in the graph
- `<graph_density>` is a float between 0 and 1 representing the density of edges in the graph

Example:
```
./graph_algo -algo floyd-warshall -size 1000 -density 0.01
```

## Project Structure

- `cmd/graph_algo/main.go`: Main entry point of the application
- `internal/graph/graph.go`: Graph data structure implementation
- `internal/algorithms/`:
  - `bfs.go`: BFS algorithm implementations
  - `bellman_ford.go`: Bellman-Ford algorithm implementations
  - `floyd_warshall.go`: Floyd-Warshall algorithm implementations
  - `spectral_clustering.go`: Spectral Clustering algorithm implementation

## Performance

The program compares the performance of regular and matrix-based implementations for each algorithm. Results are displayed in milliseconds, and a speedup factor is calculated for easy comparison.
