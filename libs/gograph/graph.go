package gograph

import (
	"errors"
)

// A Graph is connects vertices of type V with edges weighted by
// W, with metadata hinted by H and provided by M
type Graph[V comparable, W comparable, H any, M any] struct {
	isBidirectional bool                          `json:"IsBidrectional"`
	graph           map[V]Connections[V, W, H, M] `json:"Graph"`
	vertexCount     int                           `json:"VertexCount"`
	edgeCount       int                           `json:"EdgeCount"`
}

var (
	ErrEdgeExists = errors.New("edge exists between given vertices")
)

// Creates a new graph with vertices of type V with edges weighted by
// W, with metadata hinted by H and provided by M
func NewGraph[V comparable, W comparable, H any, M any](isBidirectional bool) *Graph[V, W, H, M] {
	return &Graph[V, W, H, M]{
		isBidirectional: isBidirectional,
		graph:           make(map[V]Connections[V, W, H, M]),
		vertexCount:     0,
		edgeCount:       0,
	}
}

// Creates a new graph with by deep copy
func CopyGraph[V comparable, W comparable, H any, M any](src *Graph[V, W, H, M]) *Graph[V, W, H, M] {
	newGraph := make(map[V]Connections[V, W, H, M], len(src.graph))
	for vertexFrom, connections := range src.graph {
		newGraph[vertexFrom] = CopyConnections(connections)
	}
	return &Graph[V, W, H, M]{
		isBidirectional: src.isBidirectional,
		graph:           newGraph,
		vertexCount:     src.vertexCount,
		edgeCount:       src.edgeCount,
	}
}

func (g *Graph[V, W, H, M]) ensureVertexAvailable(vertex V) {
	_, isAvailable := g.graph[vertex]
	if !isAvailable {
		g.graph[vertex] = make(Connections[V, W, H, M])
		g.vertexCount++
	}
}

// Get map of connected edges to vertex
func (g *Graph[V, W, H, M]) GetConnectedVertices(vertex V) Connections[V, W, H, M] {
	g.ensureVertexAvailable(vertex)
	connectedVertices := g.graph[vertex]
	return connectedVertices
}

func (g *Graph[V, W, H, M]) AddWeightedEdge(vertexFrom V, vertexTo V,
	edgeWeight W, hint H, metadata M) error {
	cFrom := g.GetConnectedVertices(vertexFrom)
	cTo := g.GetConnectedVertices(vertexTo)

	if cFrom.Exists(vertexTo) || cTo.Exists(vertexFrom) {
		return ErrEdgeExists
	}

	g.graph[vertexFrom] = *cFrom.AddWeightedEdge(vertexFrom, vertexTo, edgeWeight, hint, metadata, false)
	g.edgeCount++

	if g.isBidirectional {
		g.graph[vertexTo] = *cTo.AddWeightedEdge(vertexTo, vertexFrom, edgeWeight, hint, metadata, true)
		g.edgeCount++
	}

	return nil
}

func (g *Graph[V, W, H, M]) GetVertexCount() int {
	return g.vertexCount
}

func (g *Graph[V, W, H, M]) GetEdgeCount() int {
	return g.edgeCount
}
