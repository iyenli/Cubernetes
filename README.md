# Cubernetes
Course project for SE3356.

> - 在首次启动，需要下载weave, 可能需要20-30秒。尽量不要在这时候apply api object:)
> - 测试时可以用Script下的clear.sh清理残留的docker容器。[风险提示: Docker和IPtables会被清空]

## Quick start

在开始之前，您需要安装ETCD, 放到合适的目录下(Suggest: `/usr/local/bin`). 建议安装Nginx以获取默认的配置作为Volume. 安装`apt-get install resolvconf`

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

然后`./build/cubectl describe svc uid`得到Cluster IP后，你可以通过它访问：

<img src="https://s2.loli.net/2022/05/19/SUB9DoVuRTFgCbm.png" alt="image-20220514103639327" style="zoom: 50%;" />

至此，无论是docker内/外都可以正常访问Service了。而且多机上也能很好的支持。测试rs:

```shell
./build/cubectl apply -f ./example/yaml/test-replicaset.yaml
# 检查可用的副本数
./build/cubectl get rs
# 由RS中的Pod为Service提供服务
./build/cubectl apply -f ./example/yaml/service-rs.yaml
```

如果想测试负载均衡，可以修改Nginx config. 将`example/html`下的`nginx.conf`copy到运行环境的`/etc/nginx/nginx.conf`. 再将两个HTML文件copy到`/var/www/html/`下。启动`test-pod(2).yaml`，并启动`service.yaml`. 然后访问Cluster IP，就可以发现是哪个Pod提供的服务。

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

## Test Routines

测试用例在./example/presentation下。

### Init

```
./build/cuberoot init -f ./example/yaml/master-node.yaml
./build/cuberoot join 192.168.1.6 -f ./example/yaml/slave-node.yaml

# 可以展示
weave status peers
```

### Pod

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

### GPU

由于GPU任务可能发生排队，可以提前检查。

```shell
./build/cubectl apply -f ./example/yaml/presentation/gpu/gpu-add.yaml
./build/cubectl apply -f ./example/yaml/presentation/gpu/gpu-mult.yaml

# Job也有基本的RR负载均衡
```

### Svc

这部分将和RS, DNS一起检查，先检查SVC与负载均衡，再检查杀死Pod之后的行为；最后检查DNS.

```shell
./build/cubectl apply -f ./example/yaml/presentation/rs.yaml

./build/cubectl get rs
./build/cubectl apply -f ./example/yaml/presentation/svc-1.yaml
./build/cubectl apply -f ./example/yaml/presentation/svc-2.yaml

curl 172.16.0.0:8080
curl 172.16.0.1:80
./build/cubectl apply -f ./example/yaml/presentation/dns.yaml
./build/cubectl describe dns xxx

curl example.cubernetes.nb.weave.local/test/cubernetes/nb
curl example.cubernetes.nb.weave.local/test/cubernetes/very/nb
```



### AS



### 容错





### Schedule

这部分将运行2个符合仅符合一个node标签的Pod, 观察他们是否被部署到了唯一符合的node上。





### Finally

```shell
./build/cuberoot reset
./build/cuberoot stop
```



## Bug

Delete RS不成功

<img src="C:/Users/11796/AppData/Roaming/Typora/typora-user-images/image-20220522111638328.png" alt="image-20220522111638328" style="zoom:67%;" />
