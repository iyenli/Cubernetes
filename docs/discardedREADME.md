# [Discarded]! Docs

Worker1作为master, Worker2作为slave. 假设Worker1 的IP为192.168.1.9, Worker2的IP为192.168.1.5. 首次启动, 位于 ./Cubernetes目录下。

```shell
# worker1
./build/cuberoot init -f ./example/yaml/master-node.yaml
# worker2
./build/cuberoot join 192.168.1.9 -f ./example/yaml/slave-node.yaml
```

运行`weave status peers`, 在每台主机上都应该能看到对方。如果不是首次启动，可以：

````shell
./build/cuberoot start
````

停止Cubernetes集群的工作，可以：

```shell
# 如果想重新启动配置Node
./build/cuberoot reset
./build/cuberoot stop
# if testing, MAKE!
bash ./scripts/clear.sh
```

在Master节点上使用Cubectl控制：

```shell
./build/cubectl apply -f ./example/yaml/test-pod.yaml 
./build/cubectl get pods
```

测试Service:

```shell
./build/cubectl apply -f ./example/yaml/test-pod-2.yaml 
./build/cubectl apply -f ./example/yaml/service.yaml 
```

至此，无论是docker内/外都可以正常访问Service了。而且多机上也能很好的支持。测试rs:

```shell
./build/cubectl apply -f ./example/yaml/test-replicaset.yaml
# 检查可用的副本数
./build/cubectl get rs
# 由RS中的Pod为Service提供服务
./build/cubectl apply -f ./example/yaml/service-rs.yaml
```

测试DNS时，建议启动两个Pod和Service.

```shell
./build/cubectl apply -f ./example/yaml/dns/pod1.yaml 
./build/cubectl apply -f ./example/yaml/dns/pod2.yaml 

./build/cubectl apply -f ./example/yaml/dns/service1.yaml 
./build/cubectl apply -f ./example/yaml/dns/service2.yaml 
# change service ip!
./build/cubectl apply -f ./example/yaml/dns/dns.yaml 
```

如果需要测试GPU：

```shell
./build/cubectl apply -f ./example/yaml/test-gpujob.yaml
```

## 