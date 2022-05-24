package indexer

type Indexer interface {
	Start() error
	Stop() error
	Init() error
	Status() interface{}
}
