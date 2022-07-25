#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
RAND_STR=`cat /proc/sys/kernel/random/uuid | sed 's/[-]//g' | head -c 5; echo;`
# -------------------------

if [ "$#" -ne 0 ]; then
    RAND_STR=$1
    echo -e "${HI}>> Reusing previous data dir ${NC}_pgdata/.pgdata_${RAND_STR}"
else
    PSQLDIR="${PWD}/_pgdata/.pgdata_${RAND_STR}"
    echo -e "${HI}>> Create data dir for postgres ${NC}_pgdata/.pgdata_${RAND_STR}"
    mkdir -p ${PSQLDIR}
fi

echo -e "${HI}>> (Re)spawning postgres instance${NC}"
EXISTING_CONTAINERS=`sudo docker ps -f "name=pg_*" -a -q`
# sudo docker stop ${EXISTING_CONTAINERS}
sudo docker kill ${EXISTING_CONTAINERS} &> /dev/null
sleep 1
sudo docker rm ${EXISTING_CONTAINERS} &> /dev/null

sudo docker run -d \
    --name pg_${RAND_STR} \
    -e POSTGRES_USER=devuser \
    -e POSTGRES_PASSWORD=devpass \
    -e POSTGRES_DB=devdb \
    -v ${PSQLDIR}:/var/lib/postgresql/data \
    -p 5432:5432 \
    postgres &> /dev/null

sleep 1
sudo docker ps -f "name=pg_*"

sleep 2
echo -e "${HI}>> Running DB migrations ${NC}"
migrate \
    -database postgresql://devuser:devpass@localhost:5432/devdb?sslmode=disable \
    -path db/migrations up

echo -e "${HI}>> Connect using pgcli for superuser: ${NC}\n\tpgcli postgresql://devuser:devpass@localhost:5432/devdb\n"
echo -e "${HI}>> Connect using pgcli for realtime programmatic user: ${NC}\n\tpgcli postgresql://proguser:progpass@localhost:5432/devdb\n"