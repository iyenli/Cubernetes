#!/bin/bash

# Warning! This may remove all your docker. Please be careful.
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

# Warning! This may cut down your connection between you and your host:)
iptables -t nat -F

weave reset
weave stop

# Clear all chains cubernetes created
for line in `iptables -t nat -S | grep CUBE-SVC | awk '{print $2}'`
do
  iptables -t nat -X "${line}"
done

# restart docker
service docker restart
