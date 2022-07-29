#!/bin/sh
HI='\e[1;34m'
CI='\e[0;33m'
NC='\e[0m' # No Color
PWD=`pwd`
RAND_STR=`cat /proc/sys/kernel/random/uuid | sed 's/[-]//g' | head -c 5; echo;`
# -------------------------

echo -e "${HI}>> (Re)spawning rabbit mq instance${NC}"
EXISTING_CONTAINERS=`sudo docker ps -f "name=rmq_*" -a -q`
# sudo docker stop ${EXISTING_CONTAINERS}
sudo docker kill ${EXISTING_CONTAINERS} &> /dev/null
sleep 1
sudo docker rm ${EXISTING_CONTAINERS} &> /dev/null

sudo docker run -d \
    --name rmq_${RAND_STR} \
    --hostname rmq_${RAND_STR} \
    -e RABBITMQ_DEFAULT_USER=devuser \
    -e RABBITMQ_DEFAULT_PASS=devpass \
    -p 5672:5672 \
    -p 15672:15672 \
    rabbitmq:3-management &> /dev/null

sleep 1
sudo docker ps -f "name=rmq_*"

echo -e "${HI}>> Access management console in browser: ${NC}\n\thttp://devuser:devpass@localhost:15672/#/queues\n"