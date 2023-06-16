#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
# -------------------------

# if ! command -v abigen &> /dev/null; then
#     echo -e "${CI}Warning: cannot find command: abigen, setting up${NC}"
#     cd /tmp
#     rm -rf go-ethereum
#     git clone https://github.com/ethereum/go-ethereum.git --depth 1
#     cd go-ethereum
#     make devtools
#     cd $PWD
# fi

rm -rf assets/abi/*
mkdir assets/abi

echo -e "${HI}>> Generating code for ERC20${NC}"
mkdir -p assets/abi/ERC20
abigen --abi=assets/contracts/erc20abi.json --pkg=ERC20 --out=assets/abi/ERC20/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V2 pair${NC}"
mkdir -p assets/abi/univ2pair
abigen --abi=assets/contracts/uniswapv2pair.json --pkg=univ2pair --out=assets/abi/univ2pair/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V3 pair${NC}"
mkdir -p assets/abi/univ3pair
abigen --abi=assets/contracts/uniswapv3pair.json --pkg=univ3pair --out=assets/abi/univ3pair/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V3 positions NFT${NC}"
mkdir -p assets/abi/univ3positionsnft
abigen --abi=assets/contracts/uniswapv3positionsNFT.json --pkg=univ3positionsnft --out=assets/abi/univ3positionsnft/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Chainlink Oracle${NC}"
mkdir -p assets/abi/chainlink
abigen --abi=assets/contracts/chainlink.json --pkg=chainlink --out=assets/abi/chainlink/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Trader Joe V2${NC}"
mkdir -p assets/abi/traderjoev2
abigen --abi=assets/contracts/traderjoev2.json --pkg=traderjoev2 --out=assets/abi/traderjoev2/generated_do_not_edit.go