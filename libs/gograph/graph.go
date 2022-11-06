package gograph

import (
	"errors"
)

// A Graph is connects vertices of type V with edges weighted by
// W, with metadata hinted by H and provided by M
type Graph[V comparable, W comparable, H any, M any] struct {
	IsBidirectional bool                          `json:"IsBidrectional"`
	Graph           map[V]Connections[V, W, H, M] `json:"Graph"`
	VertexCount     int                           `json:"VertexCount"`
	EdgeCount       int                           `json:"EdgeCount"`
}

var (
	ErrEdgeExists = errors.New("edge exists between given vertices")
)

// Creates a new graph with vertices of type V with edges weighted by
// W, with metadata hinted by H and provided by M
func NewGraph[V comparable, W comparable, H any, M any](isBidirectional bool) *Graph[V, W, H, M] {
	return &Graph[V, W, H, M]{
		IsBidirectional: isBidirectional,
		Graph:           make(map[V]Connections[V, W, H, M]),
		VertexCount:     0,
		EdgeCount:       0,
	}
}

// Creates a new graph with by deep copy
func CopyGraph[V comparable, W comparable, H any, M any](src *Graph[V, W, H, M]) *Graph[V, W, H, M] {
	newGraph := make(map[V]Connections[V, W, H, M], len(src.Graph))
	for vertexFrom, connections := range src.Graph {
		newGraph[vertexFrom] = CopyConnections(connections)
	}
	return &Graph[V, W, H, M]{
		IsBidirectional: src.IsBidirectional,
		Graph:           newGraph,
		VertexCount:     src.VertexCount,
		EdgeCount:       src.EdgeCount,
	}
}

func (g *Graph[V, W, H, M]) ensureVertexAvailable(vertex V) {
	_, isAvailable := g.Graph[vertex]
	if !isAvailable {
		g.Graph[vertex] = make(Connections[V, W, H, M])
		g.VertexCount++
	}
}

// Get map of connected edges to vertex
func (g *Graph[V, W, H, M]) GetConnectedVertices(vertex V) Connections[V, W, H, M] {
	g.ensureVertexAvailable(vertex)
	connectedVertices := g.Graph[vertex]
	return connectedVertices
}

func (g *Graph[V, W, H, M]) AddWeightedEdge(vertexFrom V, vertexTo V,
	edgeWeight W, hint H, metadata M) error {
	cFrom := g.GetConnectedVertices(vertexFrom)
	cTo := g.GetConnectedVertices(vertexTo)

	if cFrom.Exists(vertexTo) || cTo.Exists(vertexFrom) {
		return ErrEdgeExists
	}

	g.Graph[vertexFrom] = *cFrom.AddWeightedEdge(vertexFrom, vertexTo, edgeWeight, hint, metadata, false)
	g.EdgeCount++

	if g.IsBidirectional {
		g.Graph[vertexTo] = *cTo.AddWeightedEdge(vertexTo, vertexFrom, edgeWeight, hint, metadata, true)
		g.EdgeCount++
	}

	return nil
}

func (g *Graph[V, W, H, M]) GetVertexCount() int {
	return g.VertexCount
}

func (g *Graph[V, W, H, M]) GetEdgeCount() int {
	return g.EdgeCount
}
