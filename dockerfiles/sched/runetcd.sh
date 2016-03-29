#!/bin/bash

echo "***************************************************************************************"
set
echo "***************************************************************************************"

./etcd/bin/etcd --advertise-client-urls=http://${HOST}:2379,http://${HOST}:4001 --listen-client-urls=http://${HOST}:2379,http://${HOST}:4001 >./etcd.log 2>&1 &
sleep 1
head -10 ./etcd.log
echo "*** etcd started ***"
