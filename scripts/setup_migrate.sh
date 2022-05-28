#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
RAND_STR=`cat /proc/sys/kernel/random/uuid | sed 's/[-]//g' | head -c 5; echo;`
# -------------------------

echo -e "${HI}>> Checking if migrate is installed${NC}"
if command -v migrate &> /dev/null; then
    MIGRATELOC=`command -v migrate`
    echo -e "migrate is already installed, found at ${MIGRATELOC}"
    exit
fi

echo -e "${HI}>> Checking if go exists${NC}"
if ! command -v go &> /dev/null; then
    echo "Error: cannot find command: go"
    exit
else
    GOLOC=`command -v go`
    GOVER=`go version`
    echo "Found go at ${GOLOC}: ${GOVER}"
fi

TMPDIR="/tmp/migrate_${RAND_STR}"
echo -e "${HI}>> Cloning into${NC} ${TMPDIR}"
git clone https://github.com/golang-migrate/migrate.git ${TMPDIR}

echo -e "${HI}>> Building migrate and installing...${NC}"
cd ${TMPDIR}
make build
sudo mv migrate /usr/local/bin/migrate
rm -rf ${TMPDIR}