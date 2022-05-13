#!/bin/bash

# Warning! This may remove all your docker. Please be careful.
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

# Warning! This may cut down your connection between you and your host:)
iptables -t nat -F