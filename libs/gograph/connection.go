package gograph

type WeightedEdge[W comparable, V any] struct {
	Weight   W
	Metadata V
}
type Connections[K comparable, W comparable, V any] map[K]WeightedEdge[W, V]

func (s *Connections[K, W, V]) Exists(item K) bool {
	if s == nil {
		return false
	}
	_, exists := (*s)[item]
	return exists
}

func (s *Connections[K, W, V]) Added(item K, weight W, edge V) *Connections[K, W, V] {
	if s == nil {
		sNew := make(Connections[K, W, V])
		sNew[item] = WeightedEdge[W, V]{weight, edge}
		return &sNew
	}
	(*s)[item] = WeightedEdge[W, V]{weight, edge}
	return s
}
