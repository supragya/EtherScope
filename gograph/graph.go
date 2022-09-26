package gograph

import (
	"errors"
	"math"
	"os"

	"github.com/alecthomas/binary"
)

type Tuple2[K comparable] struct {
	One K
	Two K
}

type Graph[K comparable, V any] struct {
	IsBidirectional      bool                            `json:"IsBidrectional"`
	Vertices             map[K]Connections[K, uint64, V] `json:"Vertices"`
	VertexCount          int                             `json:"VertexCount"`
	EdgeCount            int                             `json:"EdgeCount"`
	AllPairShortestPaths map[Tuple2[K]]Route[K, V]       `json:"AllPairShortestPath"`
}

var (
	ErrEdgeExists = errors.New("edge exists between given vertices")
)

func NewGraphStringUintString(bidirectional bool) *Graph[string, string] {
	return &Graph[string, string]{
		IsBidirectional:      bidirectional,
		Vertices:             make(map[string]Connections[string, uint64, string]),
		VertexCount:          0,
		EdgeCount:            0,
		AllPairShortestPaths: nil,
	}
}

func (g *Graph[K, V]) ensureVertexAvailable(vertex K) {
	_, isAvailable := g.Vertices[vertex]
	if !isAvailable {
		g.Vertices[vertex] = make(Connections[K, uint64, V])
		g.VertexCount++
	}
}

func (g *Graph[K, V]) GetConnectedVertices(vertex K) Connections[K, uint64, V] {
	g.ensureVertexAvailable(vertex)
	connectedVertices := g.Vertices[vertex]
	return connectedVertices
}

func (g *Graph[K, V]) AddEdge(from K, to K, weight uint64, edge V) error {
	cFrom := g.GetConnectedVertices(from)
	cTo := g.GetConnectedVertices(to)

	if cFrom.Exists(to) || cTo.Exists(from) {
		return ErrEdgeExists
	}

	g.Vertices[from] = *cFrom.Added(to, weight, edge)
	if g.IsBidirectional {
		g.Vertices[to] = *cTo.Added(from, weight, edge)
	}

	g.EdgeCount++

	return nil
}

func (g *Graph[K, V]) GetVertexCount() int {
	return g.VertexCount
}

func (g *Graph[K, V]) GetEdgeCount() int {
	return g.EdgeCount
}

func (g *Graph[K, V]) CalculateAllPairShortestPath() {
	if g.AllPairShortestPaths != nil {
		// Already calculated, available in cache
		return
	}

	// Create a map
	routeMap := make(map[Tuple2[K]]Route[K, V],
		g.GetVertexCount()*g.GetVertexCount())

	for from, connections := range g.Vertices {
		for to := range g.Vertices {
			if connections.Exists(to) {
				weightedEdge := connections[to]
				routeMap[Tuple2[K]{from, to}] = Route[K, V]{
					Vertices: []K{from, to},
					Edges:    []WeightedEdge[uint64, V]{weightedEdge},
					Distance: weightedEdge.Weight,
				}
				continue
			}
			routeMap[Tuple2[K]{from, to}] = Route[K, V]{
				Vertices: []K{from, to},
				Edges:    []WeightedEdge[uint64, V]{},
				Distance: math.MaxUint64,
			}
		}
	}

	for intermediate := range g.Vertices {
		for from := range g.Vertices {
			for to := range g.Vertices {
				var (
					routeFI = routeMap[Tuple2[K]{from, intermediate}]
					distFI  = routeFI.Distance
					routeIT = routeMap[Tuple2[K]{intermediate, to}]
					distIT  = routeIT.Distance
					distFT  = routeMap[Tuple2[K]{from, to}].Distance
				)

				isValidIntermediate := (distFI != math.MaxUint64) && (distIT != math.MaxUint64)
				isDetourBetter := distFT > (distFI + distIT)

				if isValidIntermediate && isDetourBetter {
					err := routeFI.AppendRoute(&routeIT)
					if err != nil {
						panic(err)
					}
					routeMap[Tuple2[K]{from, to}] = routeFI // includes appended route to To
				}
			}
		}
	}

	g.AllPairShortestPaths = routeMap
}

func (g *Graph[K, V]) GetShortestRoute(from K, to K) Route[K, V] {
	g.CalculateAllPairShortestPath()
	route := g.AllPairShortestPaths[Tuple2[K]{from, to}]
	return route
}

func (g *Graph[K, V]) SaveToDisk(fileLocation string) error {
	g.CalculateAllPairShortestPath()

	fo, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	bin, err := binary.Marshal(g)
	if err != nil {
		return err
	}
	_, err = fo.Write(bin)
	return err
}

func (g *Graph[K, V]) ReadFromDisk(fileLocation string) error {
	dat, err := os.ReadFile(fileLocation)
	if err != nil {
		return err
	}

	return binary.Unmarshal(dat, g)
}
