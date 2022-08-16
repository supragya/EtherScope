#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
# -------------------------

if ! command -v abigen &> /dev/null; then
    echo -e "${CI}Warning: cannot find command: abigen, setting up${NC}"
    cd /tmp
    rm -rf go-ethereum
    git clone https://github.com/ethereum/go-ethereum.git --depth 1
    cd go-ethereum
    make devtools
    sudo cp build/bin/abigen /usr/bin/abigen
    cd $PWD
fi

rm -rf abi/*
mkdir abi

echo -e "${HI}>> Generating code for ERC20${NC}"
mkdir -p abi/ERC20
abigen --abi=contracts/erc20abi.json --pkg=ERC20 --out=abi/ERC20/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V2 pair${NC}"
mkdir -p abi/univ2pair
abigen --abi=contracts/uniswapv2pair.json --pkg=univ2pair --out=abi/univ2pair/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V3 pair${NC}"
mkdir -p abi/univ3pair
abigen --abi=contracts/uniswapv3pair.json --pkg=univ3pair --out=abi/univ3pair/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V3 positions NFT${NC}"
mkdir -p abi/univ3positionsnft
abigen --abi=contracts/uniswapv3positionsNFT.json --pkg=univ3positionsnft --out=abi/univ3positionsnft/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Chainlink Oracle${NC}"
mkdir -p abi/chainlink
abigen --abi=contracts/chainlink.json --pkg=chainlink --out=abi/chainlink/generated_do_not_edit.go