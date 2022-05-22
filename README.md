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

## Asserting version of a built artifact
`bgidx` while building is stamped with build information such as git commit, CI process, time of build etc. so as to allow identification while in production.
```
$ bgidx --version
bgidx version 1.0.0 build main@1923df1876d27fa0092089d665d38e0bfbe31aac
compiled at 22-05-2022_03-14-39@UTC by travis
```

## Running indexer "on-head"
TODO

## Running indexer as a "backfiller"
TODO