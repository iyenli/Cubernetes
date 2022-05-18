#!/bin/bash

# Warning! This may remove all your docker. Please be careful.
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

# Warning! This may cut down your connection between you and your host:)
# iptables -t nat -F

weave reset
weave stop

# Clear all chains cubernetes created
# far more accurate than -F
for line in `iptables -t nat -S | grep CUBE-SVC | awk '{print $2}'`
do
  count=`iptables -t nat -L "${line}" | wc -l`
  count=`expr $count - 2`
  for((i=1;i<=$count;i++));
  do
    iptables -t nat -D $line 1
  done

  iptables -t nat -X "${line}"
done

# clear all rules in service
service_count=`iptables -t nat -L SERVICE | wc -l`
service_count=`expr $count - 2`
for((i=1;i<=$service_count;i++));
do
  iptables -t nat -D SERVICE 1
done

# restart docker
service docker restart
