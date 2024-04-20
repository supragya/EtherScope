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
COPY ./scripts scripts
COPY ./services services
COPY ./types types
COPY ./version version
COPY ./algorand algorand

RUN go build -ldflags="\
  -X 'github.com/supragya/EtherScope/version.buildCommit=$buildCommit' \
  -X 'github.com/supragya/EtherScope/version.gittag=$gittag' \
  -X 'github.com/supragya/EtherScope/version.buildTime=$buildtime' \
  -X 'github.com/supragya/EtherScope/version.builder=$builder' \
  -X 'github.com/supragya/EtherScope/version.gover=$gover'" \
  -o build/escope

# App did not start correctly in Alpine, but perhaps a solution for this could be found
FROM golang:1.19.3

WORKDIR /geth-indexer

COPY --from=build /geth-indexer/build/escope /geth-indexer/escope

RUN "printf \"%s\" \"$CONFIG\" > ./config.yaml"

CMD ["./escope", "realtime", "-c", "config.yaml"]
