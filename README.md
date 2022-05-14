# Cubernetes
Course project for SE3356.

> - 在首次启动，需要下载weave, 可能需要20-30秒。尽量不要在这时候apply api object:)
> - 测试时可以用Script下的clear.sh清理残留的docker容器。[风险提示: Docker和IPtables会被清空]

## Quick start

在开始之前，您需要安装ETCD, 放到合适的目录下(Suggest: `/usr/local/bin`). 建议安装Nginx以获取默认的配置作为Volume.

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
bash ./scripts/clear.sh
```

<img src="https://s2.loli.net/2022/05/12/Uy61jQ9cpbK2ZR4.png" alt="image-20220512094310090" style="zoom: 80%;" />

在Master节点上使用Cubectl控制：
```shell
./build/cubectl apply -f ./example/yaml/test-pod.yaml 
./build/cubectl get pods
```

至此，`curl ip:8085`可以看到Nginx主页的HTML界面。尚未测试*带命令的容器，容器内的localhost访问*。测试Service:

```shell
./build/cubectl apply -f ./example/yaml/test-pod-2.yaml 
./build/cubectl apply -f ./example/yaml/service.yaml 
```

`iptables -t nat -L`, 与K8s保持一致。

![image-20220514131850160](https://s2.loli.net/2022/05/14/iprcM7wYNFL1moR.png)

然后`./build/cubectl describe svc uid`得到Cluster IP后，你可以通过它访问：

<img src="C:/Users/11796/AppData/Roaming/Typora/typora-user-images/image-20220514103639327.png" alt="image-20220514103639327" style="zoom:67%;" />

至此，无论是docker内/外都可以正常访问Service了。而且多机上也能很好的支持。测试rs:

```shell
./build/cubectl apply -f ./example/yaml/test-replicaset.yaml
# 检查可用的副本数
./build/cubectl get rs
# 由RS中的Pod为Service提供服务
./build/cubectl apply -f ./example/yaml/service-rs.yaml
```

如果想测试负载均衡，可以修改Nginx config. 将`example/html`下的`nginx.conf`copy到运行环境的`/etc/nginx/nginx.conf`. 再将两个HTML文件copy到`/var/www/html/`下。启动`test-pod(2).yaml`，并启动`service.yaml`. 然后访问Cluster IP，就可以发现是哪个Pod提供的服务。





