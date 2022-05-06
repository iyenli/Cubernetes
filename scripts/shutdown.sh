#!/bin/bash

shutdown() {
    kill -9 $(ps -a | grep cubelet | awk '{print $1}')
    sleep 1

    kill -9 $(ps -a | grep apiserver | awk '{print $1}')
    sleep 1

    kill -9 $(ps -a | grep etcd | awk '{print $1}')
}