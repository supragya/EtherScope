package indexer

type Indexer interface {
	Start() error
	Stop() error
	Init(interface{}) error
	Status() interface{}
}
