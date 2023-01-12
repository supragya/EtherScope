FROM golang:1.19.3

WORKDIR /geth-indexer

# Copy configuration files
COPY ./build/bgidx .


# Files requiring mounting (?):
# 1. config.yaml
# 2. chainlink_oracle_dumpfile.csv
# 3. dex_dumpfile.csv

# Default to realtime and expect --entrypoint for backfill? 
# Can configuration be supplied by env rather than a file that would need to be
# mounted at runtime?
ENTRYPOINT ["/geth-indexer/bgidx", "realtime", "-c", "config.yaml"]

