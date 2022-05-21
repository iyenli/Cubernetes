#!/bin/bash

# A more robust and clear script, use it with no worry
# It's ok to see:
# - Weave is not running (ignore on Kubernetes).
echo "\033[34m Don't worry about output whatever it is, it'll clear context correctly:) \033[0m"

eval "$(weave env --restore)"
weave reset

# Warning! This may remove all your docker. Please be careful.
docker_count=$(docker ps -a -q | wc -l)
if ((docker_count > 0)); then
  docker stop $(docker ps -a -q)
  docker rm $(docker ps -a -q)
fi

weave stop

# Clear all chains cubernetes created
# far more accurate than -F
ipt_count=$(iptables -t nat -S | grep CUBE-SVC | awk '{print $2}' | wc -l)
if ((ipt_count > 0)); then
  for line in $(iptables -t nat -S | grep CUBE-SVC | awk '{print $2}'); do
    count=$(iptables -t nat -L "${line}" | wc -l)
    count=$((count - 2))
    for ((i = 1; i <= count; i++)); do
      iptables -t nat -D "$line" 1
    done

    iptables -t nat -X "$line"
  done
fi

# clear all rules in service
service_count=$(iptables -t nat -L SERVICE | wc -l)
if ((service_count > 2)); then
  service_count=$((count - 2))
  for ((i = 1; i <= service_count; i++)); do
    iptables -t nat -D SERVICE 1
  done
fi

# restart docker
service docker restart
