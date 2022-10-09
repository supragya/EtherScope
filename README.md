# Blockpour Geth Indexer
`bgidx` program aims to index an EVM like chain using go-ethereum clients via RPCs.

![Test run status](https://github.com/Blockpour/Blockpour-Geth-Indexer/actions/workflows/gotest.yaml/badge.svg?branch=feat/v0.3.0)
## Building
Build the application using:
- Setup golang abi modules via abigen using `./scripts/setup_abi.sh`.
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
- Install docker via `paru -S docker` (arch linux) and run it using `sudo systemctl start docker`
- Install psycopg2 using `sudo pacman -S python-psycopg2`
- Install `pgcli` for interacting with a postgres instance via cli using `pip3 install pgcli`
- Install migrate using `./scripts/setup_migrate.sh`
- Setup a pgsql docker container using `./scripts/start_db.sh`. This will generate a new `pgdata/.pgdata_XXXX` directory (XXXXX being random for each invocation) which will be used by postgresql. Every time this script is invoked, the DB is launched anew with no data.
- To start a docker container with a previous directory, invoke using `./scripts/start_db.sh a6df1` if concerned data directory is `pgdata/.pgdata_a6df1`. In this mode, db migrations are not run.

You should now have two users:
- **devuser**: Accessible via `pgcli postgresql://devuser:devpass@localhost:5432/devdb` for DB superuser access.
- **proguser**: Accessible via `pgcli postgresql://proguser:progpass@localhost:5432/devdb` for insert only access to `blocks` and `pool_actions_geth` tables.

## Development rabbit mq
- Install docker via `paru -S docker` (arch linux) and run it using `sudo systemctl start docker`
- Setup a rabbit MQ container using `./scripts/start_rmq.sh`. This will create a fresh cluster each time it is invoked.

Single user for all access to rmq in dev mode. Acess management console in browser using: `http://devuser:devpass@localhost:15672/#/queues`