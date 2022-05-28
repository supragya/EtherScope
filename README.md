# Blockpour Geth Indexer
`bgidx` program aims to index an EVM like chain using go-ethereum clients via RPCs.

## Building
Build the application using:
- `make build` or `make build -B` (force) to build `build/bgidx`
- `make install` to install the built `bgidx` into `/usr/local/bin/bgidx` bringing it into $PATH.
- `make uninstall` removes `/usr/local/bin/bgidx`

## Asserting version of a built artifact
`bgidx` while building is stamped with build information such as git commit, CI process, time of build etc. so as to allow identification while in production.
```
$ bgidx --version
bgidx version v0.0.1 build main@5c6c6c5
compiled at 23-05-2022 04:29:22 by travisci
using go version go1.17.8 linux/amd64
```

## Indexer modes
`bgidx` runs in two different modes: 
- **Realtime**: Aims to stay on head of the concerned blockchain and update the backend database in realtime. Can be run using `bgidx realtime -c <config.yaml>`
- **Backfill**: Aims to backfill a range of blocks in the past and update the backend database, rewriting the entries for concerned blocks. Can be run using `bgidx backfill -c <config.yaml>`

Example config file(s) is available at `test/configs/testcfg.yaml`

## Checklist
- [x] Config checks for mandatory fields
- [x] Realtime subcommand init
- [x] Backfill subcommand init
- [ ] Connection to backend DB: postresql