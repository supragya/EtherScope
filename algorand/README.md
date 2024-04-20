## Usage

Requirements
- Rabbitmq running on specified port

Create a config file
`cp config-local.yaml config.yaml`



Start the indexer
`go run main.go indexer`


## Node Setup

Verify the node is running: `goal node status` 

Update the node config in the `node/config.json` file. 

`config.json`
{
    "Archival": true,
    "EndpointAddress": 0.0.0.0:8080
}

Be sure to restart the node after modifying the configuration. 

## Indexer Daemon Setup 

./algorand-indexer daemon --data-dir /tmp -P postgres://<db_user>:<db_password><db_host>/authentication --token $(cat $ALGORAND_DATA/algod.token)

## Indexer Config 

If using self hosted node: 
algodUrl: "http://algorand.supragya.local:8080"
indexerUrl:  "http://algorand.supragya.local:8980"

If using the quicknode, the algodUrl and indexerUrl should be in this format (including token and suffix)
algodUrl: "https://<project-moniker>.algorand-mainnet.discover.quiknode.pro/<token>/algod"
indexerUrl: "https://<project-moniker>.algorand-mainnet.discover.quiknode.pro/<token>/indexer"


## Troubleshooting

`curl algorand.supragya.local:8980/health`
