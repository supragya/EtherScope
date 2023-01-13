FROM golang:1.19.3 as build

ARG buildCommit
ARG gittag
ARG buildTime
ARG builder
ARG gover

WORKDIR /geth-indexer

COPY ./assets assets
COPY ./cmd cmd
COPY ./go.mod .
COPY ./go.sum .
COPY ./libs libs
COPY ./main.go .
COPY ./_pgdata _pgdata
COPY ./scripts scripts
COPY ./services services
COPY ./types types
COPY ./version version

RUN go build -ldflags="\
  -X 'github.com/Blockpour/Blockpour-Geth-Indexer/version.buildCommit=$buildCommit' \
  -X 'github.com/Blockpour/Blockpour-Geth-Indexer/version.gittag=$gittag' \
  -X 'github.com/Blockpour/Blockpour-Geth-Indexer/version.buildTime=$buildtime' \
  -X 'github.com/Blockpour/Blockpour-Geth-Indexer/version.builder=$builder' \
  -X 'github.com/Blockpour/Blockpour-Geth-Indexer/version.gover=$gover'" \
  -o build/bgidx

FROM golang:1.19.3

WORKDIR /geth-indexer

COPY --from=build /geth-indexer/build/bgidx /geth-indexer/bgidx

# Files requiring mounting (?):
# 1. config.yaml
# 2. chainlink_oracle_dumpfile.csv
# 3. dex_dumpfile.csv

# Default to realtime and expect --entrypoint for backfill? 
# Can configuration be supplied by env rather than a file that would need to be
# mounted at runtime?
CMD ["./bgidx", "realtime", "-c", "config.yaml"]
