#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
# -------------------------

if ! command -v abigen &> /dev/null; then
    echo "Error: cannot find command: abigen"
    exit
fi

rm -rf abi/*
mkdir abi

echo -e "${HI}>> Generating code for ERC20${NC}"
mkdir -p abi/ERC20
abigen --abi=contracts/erc20abi.json --pkg=ERC20 --out=abi/ERC20/generated_do_not_edit.go

echo -e "${HI}>> Generating code for Uniswap V2 pair${NC}"
mkdir -p abi/univ2pair
abigen --abi=contracts/uniswapv2pair.json --pkg=univ2pair --out=abi/univ2pair/generated_do_not_edit.go