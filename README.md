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

## Development database
- Install docker via `paru -S docker` (arch linux)
- Install `pgcli` for interacting with a postgres instance via cli using `pip3 install pgcli`
- Setup a pgsql docker container using `./scripts/start_db.sh`. This will generate a new `pgdata/.pgdata_XXXX` directory (XXXXX being random for each invocation) which will be used by postgresql. Every time this script is invoked, the DB is launched anew with no data.
- To start a docker container with a previous directory, invoke using `./scripts/start_db.sh a6df1` if `a6df1` is concerned data directory is `pgdata/.pgdata_a6df1`
- Install migrate using `./scripts/setup_migrate.sh`
- Run dev migrations using `migrate -database postgresql://devuser:devpass@localhost:5432/devdb?sslmode=disable -path db/migrations up` assuming defaults being used from scripts above.

## Checklist
- [x] Config checks for mandatory fields
- [x] Realtime subcommand init
- [x] Backfill subcommand init
- [ ] Connection to backend DB: postresql
- [ ] Custom datadir loading into start_db.sh