#!/bin/bash


/etcd/bin/runetcd.sh

export ETCD_LOCAL_ENDPOINT=http://${HOST}:2379

echo "Starting the MrRedis scheduler"
/mrredis/bin/sched -config=/mrredis/bin/config.json
