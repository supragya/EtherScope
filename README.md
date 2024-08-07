# EtherScope Go-Ethereum Indexer
`escope` program aims to index an EVM like chain using go-ethereum clients via RPCs.


![Dashboard](assets/bpour_dashboard.png)

## Building
Build the application using:
- Setup golang abi modules via abigen using `./scripts/setup_abi.sh`.
- `make build` or `make build -B` (force) to build `build/escope`
- `make install` to install the built `escope` into `/usr/local/bin/escope` bringing it into $PATH.
- `make uninstall` removes `/usr/local/bin/escope`

## Asserting version of a built artifact
`escope` while building is stamped with build information such as git commit, CI process, time of build etc. so as to allow identification while in production.
```
$ escope --version
escope version v0.0.1 build main@5c6c6c5
compiled at 23-05-2022 04:29:22 by travisci
using go version go1.17.8 linux/amd64
```

## Indexer modes
`escope` runs in two different modes: 
- **Realtime**: Aims to stay on head of the concerned blockchain and update the backend database in realtime. Can be run using `escope realtime -c <config.yaml>`
- **Backfill**: Aims to backfill a range of blocks in the past and update the backend database, rewriting the entries for concerned blocks. Can be run using `escope backfill -c <config.yaml>`

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

## Docker 

### Building
To provide parity with existing build steps, scripts have been added to the `Makefile` for building the docker image. Use `make docker` to build the project or `make dockerbuildx` to execute a multiarch build (see [Multiarchitecture Builds](#multiarchitecture-builds) for details).

### Running the image
Currently the image expects that the runner will mount files it needs for operation. These files are `config.yaml`, `chainlink_oracle_dumpefile.csv`, and `dex_dumpfile.csv`. Additionally, the application stores stateful data in the folder `lb.badger.db` which is managed by the local backend. Currently this is also being mounted in order for the state to be maintained between application runs. A more architecturally stable solution for this should be explored in the future. For instance, if your current working directory is the project root and these files are also located at the project root, then you could run the container like this:

```
docker run \
--name=geth-indexer \
--mount type=bind,source="$(pwd)"/config.yaml,target=/geth-indexer/config.yaml \
--mount type=bind,source="$(pwd)"/chainlink_oracle_dumpfile.csv,target=/geth-indexer/chainlink_oracle_dumpfile.csv \
--mount type=bind,source="$(pwd)"/dex_dumpfile.csv,target=/geth-indexer/dex_dumpfile.csv \
--mount type=bind,source="$(pwd)"/lb.badger.db,target=/geth-indexer/lb.badger.db \
geth-indexer
```

### Multiarchitecture Builds

#### Buildx
Building for multiple architectures currently requires the docker buildx CLI plugin.  See [Docker docs: Install Docker Buildx](https://docs.docker.com/build/install-buildx/) for details on setup.

#### Limitations
Building for multiple architectures requires that the created artifacts be immediatly pushed to a docker remote. This means that the multiarchitecture builds cannot be used for local testing. Instead, for testing using your local docker server you should do the default single architecture build. See this [Github issue](https://github.com/docker/buildx/issues/59) for discussion on this topic and multiarchitecture builds should be run in a pipeline where the resulting images can be automatically pushed.
