#!/bin/bash

start_cluster() {
    etcd > .test_log/etcd.txt &
    sleep 3
    echo "etcd pid = $!"

    build/apiserver > .test_log/apiserver.txt &
    sleep 3
    echo "apiserver pid = $!"

    build/cubelet > .test_log/cubelet.txt &
    sleep 3
    echo "cubelet pid = $!"
}

start_cluster