package gograph

// A weightedEdge connects vertices of type V with
// weights defined by W with metadata hinted by H
// and provided by M
type WeightedEdge[V comparable, W comparable, H any, M any] struct {
	IsReverseEdge bool
	VertexFrom    V
	VertexTo      V
	Weight        W
	Hint          H
	Metadata      M
}

// Connections is set of weighted edges for any source
// vertex already contrained. All weighted edges in a
// particular "Connections" will always have VertexFrom
// A weightedEdge connects vertices of type V with
// weights defined by W with metadata hinted by H
// and provided by M
type Connections[V, W comparable, H, M any] map[V]WeightedEdge[V, W, H, M]

func CopyConnections[V, W comparable, H, M any](src Connections[V, W, H, M]) Connections[V, W, H, M] {
	newConnections := make(Connections[V, W, H, M], len(src))
	for vertexTo, edge := range src {
		newConnections[vertexTo] = edge
	}
	return newConnections
}

func (s *Connections[V, W, H, M]) Exists(vertex V) bool {
	if s == nil {
		return false
	}
	_, exists := (*s)[vertex]
	return exists
}

func (s *Connections[V, W, H, M]) AddWeightedEdge(vertexFrom V, vertexTo V,
	edgeWeight W, hint H, metadata M, isReverseEdge bool) *Connections[V, W, H, M] {
	edge := WeightedEdge[V, W, H, M]{
		IsReverseEdge: isReverseEdge,
		VertexFrom:    vertexFrom,
		VertexTo:      vertexTo,
		Weight:        edgeWeight,
		Hint:          hint,
		Metadata:      metadata,
	}

	if s == nil {
		sNew := make(Connections[V, W, H, M])
		sNew[vertexTo] = edge
		return &sNew
	}
	(*s)[vertexTo] = edge
	return s
}
