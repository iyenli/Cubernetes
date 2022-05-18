#!/bin/bash

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