# Cubernetes
Course project for SE3356.

> 注意，在现在，运行Cuberoot的Shell不能被关闭。
>
> 测试时可以用Script下的clear.sh清理残留的docker容器。[风险提示]

## Quick start

在开始之前，可能的前置工作：
- 安装Nginx和ETCD

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
./build/cuberoot stop
# 如果想重新启动配置Node
./build/cuberoot reset
```

<img src="https://s2.loli.net/2022/05/12/Uy61jQ9cpbK2ZR4.png" alt="image-20220512094310090" style="zoom: 80%;" />

在Master节点上使用Cubectl控制：
```shell
./build/cubectl apply -f ./example/yaml/test-pod.yaml 
./build/cubectl get pods
```

至此，`curl ip:8085`可以看到Nginx主页的HTML界面。尚未测试*带命令的容器，容器内的localhost访问*。测试Service:

```shell
./build/cubectl apply -f ./example/yaml/test-pod2.yaml 
./build/cubectl apply -f ./example/yaml/service.yaml 
```

`iptables -t nat -L`, 与K8s保持一致。

<img src="https://s2.loli.net/2022/05/13/XCQvMTdOr5eyZPV.png" alt="image-20220513215801877" style="zoom: 50%;" />

然后`./build/cubectl describe svc uid`得到Cluster IP后，你可以通过它访问：

<img src="C:/Users/11796/AppData/Roaming/Typora/typora-user-images/image-20220514103639327.png" alt="image-20220514103639327" style="zoom:67%;" />

至此，无论是docker内/外都可以正常访问Service了。



