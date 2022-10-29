package gograph

import "errors"

type Route[K comparable, V any] struct {
	Vertices []K
	Edges    []WeightedEdge[uint64, V]
	Distance uint64
}

var (
	ErrIncompatibleRoutes = errors.New("routes incompatible for concatenation")
)

func (r *Route[K, V]) AppendRoute(second *Route[K, V]) error {
	if second == nil || len(second.Vertices) == 0 {
		return nil
	}

	if r.Vertices[len(r.Vertices)-1] != second.Vertices[0] {
		return ErrIncompatibleRoutes
	}

	r.Vertices = append(r.Vertices, second.Vertices[1:]...)
	r.Edges = append(r.Edges, second.Edges...)
	r.Distance += second.Distance
	return nil
}
