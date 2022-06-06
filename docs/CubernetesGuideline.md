## Test

### Clear

```shell
./build/cuberoot reset
./build/cuberoot stop

bash ./scripts/clear.sh
```

### Case 1

```shell
./build/cuberoot init -f ./example/yaml/master-node.yaml
./build/cuberoot join 192.168.1.6 -f ./example/yaml/slave-node.yaml

./build/cubectl get nodes
```

### Case 2

```shell
# 展示基本Pod操作，command和资源限制
./build/cubectl apply -f ./example/yaml/presentation/pod/stress.yaml
# top / log看下CPU占用，应该是2个VM大致平分一个Core

./build/cubectl get pods
./build/cubectl describe pod

# 展示停止pod
./build/cubectl delete pod f8f5ce4e-5819-493e-a9c0-408b5c9b4560
# 已经没有这个Docker了

# 展示挂载卷，容器内localhost，暴露端口
# 记得回忆下tc的index.html有没有写入
# 这个Pod会被调度到另一台机器上
./build/cubectl apply -f ./example/yaml/presentation/pod/tn.yaml
# 展示Volume的网页
curl 127.0.0.1:8095
curl 127.0.0.1:8090

docker exec -it xxx bash
curl 127.0.0.1:8080
# 至此，Pod要求已经全部实现了
```

### Case 3

```shell
./build/cubectl apply -f ./example/yaml/presentation/rs.yaml

./build/cubectl get rs
./build/cubectl apply -f ./example/yaml/presentation/svc-1.yaml
./build/cubectl apply -f ./example/yaml/presentation/svc-2.yaml

# 回填svc ID
curl 172.16.0.0:8080
curl 172.16.0.1:80

./build/cubectl apply -f ./example/yaml/presentation/dns.yaml
./build/cubectl describe dns xxx

curl example.cubernetes.weave.local/test/cubernetes/nb
curl example.cubernetes.weave.local/test/cubernetes/very/nb

# 展示docker内也可以访问DNS与Svc
docker exec it ...
```

### Case 4

```shell
./build/cuberoot stop
# 检查Pod IP和Service
./build/cuberoot start
# 检查Pod IP和Service, 检查已有功能是否正常

./build/cubectl apply -f ./example/yaml/presentation/onemoreRS.yaml
```

### Case 5

```shell
./build/cubectl apply -f ./example/yaml/presentation/pod/sharedv.yaml
```

###   Case 6

```shell
./build/cubectl apply -f ./example/yaml/presentation/multipod/pod1.yaml
./build/cubectl apply -f ./example/yaml/presentation/multipod/pod2.yaml
./build/cubectl apply -f ./example/yaml/presentation/multipod/multipod-service.yaml
```

### Case 7

```shell
./build/cubectl apply -f ./example/yaml/presentation/rs-selfkilled.yaml
./build/cubectl apply -f ./example/yaml/presentation/svc-selfkilled.yaml
./build/cubectl get rs
./build/cubectl describe rs 
./build/cubectl describe svc 
```

### Case 8

```shell
./build/cubectl apply -f ./example/yaml/presentation/gpu.yaml -j ./build/matadd.tar.gz
./build/cubectl apply -f ./example/yaml/presentation/gpu.yaml -j ./build/matmult.tar.gz
```

## Bug

// TODO: Delete it

Delete RS不成功

<img src="https://s2.loli.net/2022/05/23/GCL7goy9HXOA1Sr.png" alt="image-20220522111638328" style="zoom: 50%;" />
