# Blockpour Geth Indexer
`bgidx` program aims to index an EVM like chain using go-ethereum clients via RPCs.

## Building
Build the application using (substitute `1.0.0` with version number of build):
```
make build VER=1.0.0
```
To force rebuilding, use `-B` make flag
```
make build VER=1.0.0 -B
```
This builds `build/bgidx`

To install the built `bgidx` into `/usr/local/bin/bgidx`, use `make install`

To remove the artifact located at `/usr/local/bin/bgidx` use `make uninstall`

## Running indexer "on-head"
TODO

## Running indexer as a "backfiller"
TODO