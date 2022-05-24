# Blockpour Geth Indexer
`bgidx` program aims to index an EVM like chain using go-ethereum clients via RPCs.

## Building
Build the application using:
```
make build
```
To force rebuilding, use `-B` make flag
```
make build -B
```
This builds `build/bgidx`

To install the built `bgidx` into `/usr/local/bin/bgidx`, use `make install`

To remove the artifact located at `/usr/local/bin/bgidx` use `make uninstall`

## Asserting version of a built artifact
`bgidx` while building is stamped with build information such as git commit, CI process, time of build etc. so as to allow identification while in production.
```
$ bgidx --version
bgidx version v0.0.1 build main@5c6c6c5
compiled at 23-05-2022 04:29:22 by s20y671
using go version go1.17.8 linux/amd64
```

## Running indexer "on-head"
Run indexer in `realtime` mode using
```
./build/bgids realtime -c <config_file.yaml>
```

## Running indexer as a "backfiller"
TODO