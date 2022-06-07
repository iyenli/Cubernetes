# 验收报告



## 架构介绍

### 使用的依赖

- 没有使用k8s.io组织下的任何依赖和库，没有使用Knative, openwhisk有关的任何依赖和库
- 没有使用CNI库，但容器间网络连接是利用Weave实现的
- 主要使用了etcd, kafka, iptables client封装对这些外部组件的操作

<img src="https://s2.loli.net/2022/06/07/UVjySh7TCYsdG9P.png" alt="image-20220607193721317" style="zoom: 67%;" />

- 自行实现了简易的时序数据库作为Serverless动态扩缩容的依据

### 多机部署

- cuberoot通过是否存在元信息持久化文件判断应该新创建并且加入一个Cubernetes节点还是恢复之前的数据
- Scheduler调度优先满足Pod selector, 在匹配的Node集合中，Scheduler会进行Round Robin调度
- Worker通过监测Watch channel的状态判断Master的存货状态，周期性的重连实现容错
- Master通过心跳检测探测worker状态，更新ETCD中存储的Node状态

### Cubernetes网络

- 以Weave Plugin为基础，用户**无感知**的下载与配置插件，不需要用户额外配置任何网络
  - Pod中容器共享Pause提供的网络，Pause加入Weave net, Pod中容器共享Localhost
  - Container内可以通过任意节点的IP，任意Pod IP和Service IP 以及注册的域名访问对应服务
  - 动态节点加入与删除
- Service: 基于IPtables, 自动化的配置
  - 自定义多层规则链，相互引用，便于更新，删除，新增Service
- 容器化Nginx实现DNS服务
  - 观察到不同Path对应不同的IP超出了DNS职责，通过Nginx实现不同path到不同IP的转换
  - 复用Service进行Pod负载均衡

### GPU服务

通过负载均衡的部署GPU Server到各个节点上，复用原有的Pod管理机制。具体来说，由Cubernetes提供镜像，用户提交Slurm文件后，将运行相应的GPU Server Pod与远程HPC服务器进行交互，然后直接和API Server反馈任务状态和结果。用户通过cubectl进行查询。

### Serverless

#### Gateway

- 提供HTTP服务，接受HTTP Trigger并且返回Serverless计算结果
- Gateway构建为镜像，可以实现无感知更新；复用AutoScaler + Service做负载均衡，对外暴露固定的Service IP
- 通过Kafka Topic Partition机制实现消息队列的可扩展性
- 每个请求通过Go Channel和long-running的请求返回队列监听者通信，减少切换开销



### Cubelet

Cubelet组件运行在每个Node上，负责对本地运行的Pod对象进行管理：包括从API-Server监测对Pod对象的状态操作（如创建、更新和删除），以及定时监控、计算本地所有Pod对象的状态，并将其同步到API-Server中。

出于通用性和效率的权衡考虑，Cubelet将管理的“类Pod对象”分为三种：

+   最为泛用的Pod对象，一般用于持续运行的workload，具有刷新Spec，监控资源利用状态等所有功能，运行开销也最大，是项目使用与演示中主要用到的对象；
+   GPU应用专属的GpuJob对象，是专门为运行一段GPU计算任务（Job）实现的抽象；
+   Serverless部分中运行Action实例的Actor，使用一个内部定义的Python-Runtime镜像，并且会自动生成其专属的cmd，也不监控其资源利用率，有较小的运行开销和较快的启动速度。

Cubelet组件主要由以下几个层次组成（调用栈从上到下）：

1.   Informer：

     各个对象的informer组件使用API-Server提供的List-And-Watch接口，监控其中对象的状态变化，并暴露为统一的Event管道，以便调用各种事件（Create、Update、Remove）的Handler。在这层，基于NodeID过滤掉尚未Schedule/非本地的Pod对象。

     由于informer监控了所有属于当前Node的类Pod资源的状态变化，因此也在其中维护了每个资源的本地Cache，可以避免Cubelet大量向API-Server发送查询请求，减少其负载。

     List-And-Watch接口也可用来监控本地节点与API-Server的连接情况，如果控制节点下线，Cubelet将在这层反复尝试与控制面重新建立连接，达成容错的效果。

2.   Handler：

     这一层是Cubelet的主要逻辑处理部分：通过Select机制监听Informer层提供的事件channel，调用各种Pod对象的Create、Update、Remove事件的响应处理逻辑函数；同时对于每种类Pod对象都有一个定期Routine，监控、计算出本地运行对象的状态，将其更新到API-Server中。

3.   Runtime：

     对于三种不同的类Pod对象，基于Dockershim暴露的docker相关接口，封装了对象独有的启动、同步、停止、监控等操作，供上层Handler进行调用，使得编写Handler时只用关心业务逻辑而无需直接与docker接口交互。

4.   Dockershim：

     对于docker sdk的简单封装，提供给Runtime使用，用于对容器、镜像的各种操作。



### Controller-Manager

Controller-Manager运行在控制面上，负责保证高级对象（ReplicaSet、AutoScaler）的运行状态符合其定义。该组件只与API-Server组件交互，从其中监测各类对象的状态变化，计算出为了达到目标状态所需进行的操作，并将操作通知到API-Server中。由于交互链较为简单，因此该组件的分层逻辑也很简洁：

1.   Informer：

     与Cubelet组件中的Informer相似，也通过List-And-Watch机制监控了API-Server，提供了组件内的对象Cache（Pod、ReplicaSet、AutoScaler），向下暴露了事件通知channel，提供容错机制；由于多个Controller可能需要监控同一个资源，Informer可以同时建立多个事件通知channel，多个Controller也能共用同样的Cache，保证了资源信息一致性。

2.   Controller：

     每个高级对象（本项目中是ReplicaSet和AutoScaler）都对应了一个Controller，可以按需传入Informer作为参数创建，以此监控对应的事件通知channel。Controller中通过handler处理每种资源对象的Create、Update、Remove事件，也会通过定时Routine检查高级对象是否符合目标状态，并通过合理的更新API对象使其达到目标状态。

以上结构可以方便地扩展：对于新的高级对象（如ReplicaSet），只需要实现其Controller中的控制逻辑，再传入所需的Informer即可开始监视-控制该类对象；对于新的被监控对象（Pod），也只需要添加Informer即可。

### Action-Brain

Action-Brain运行在控制面上，可以看作为Serverless特化的Controller。其处理逻辑和Controller-Manager类似，不同的是，Action-Brain主要控制Action对象（如其名），还通过一个Monitor子组件从kafka消息队列中监听对函数的每一次调用，记录在轻量级时序数据库tstorage中，提供Action扩容的RPM数据。

## 成员分工



## 工程实践

### Gitee 地址: 

- Cubernetes https://gitee.com/k9-s/Cubernetes
- Serverless Python框架 https://gitee.com/k9-s/Serverless-Python

### 项目分支

- master 项目主分支，长期存在。仅由develop branch通过pull request合入，合并时机为完成了一个完整的，经过集成测试的主要功能之后；
- develop 项目开发分支，长期存在。仅由功能分支通过pull request合入，合并时机为完成了一个经过测试的功能点之后；
- 其他功能分支，有明确的生命周期。在开启一项独立的任务之后，由任务负责开发者创建，将在任务开发完毕后合入develop分支(或者使用此分支feature的其他分支)并且删除。

### CI/CD

使用Gitee[仓库镜像](https://gitee.com/help/articles/4336)到Github私有仓库，使用Github Campus(白嫖)的Github Action time(3K minutes/month)进行CI. 使用的Workflow配置yaml文件可以参考附件1. 主要完成的工作：

* 运行`go build`确保编译通过
* 运行`go vet`捕获是否存在编译器未发现的错误
* 运行`staticcheck`进行静态代码检查
* 运行`go test`完成单元测试
* 如果CI不通过，会通过邮件通知小组

Github Action在镜像同步后会自动触发：

<img src="https://s2.loli.net/2022/06/07/u81hBGRo9WanevA.png" alt="image-20220607190844772" style="zoom:80%;" />

<img src="https://s2.loli.net/2022/06/07/aYtXODCPKb3y5Bp.png" alt="image-20220607191603092" style="zoom: 50%;" />

由于本身云操作系统并非持续服务的Web server或其他服务，因此没有进行持续部署。用于集成测试的部署方式是基于Jetbrain Deployment和Linux rsync进行同步与部署的：

- 在Jetbrain worker group 中创建两台JCloud服务器，填入ssh连接信息
- 配置rsync连接，将本地Cubernetes开发目录映射到部署目录，这一步需要排除build等目录
- 在Jetbrain Goland中配置rsync，在git merge/代码本地修改时触发rsync到远程服务器Group
- 在JCloud机器上进行构建和部署

### 软件测试

- 单元测试使用go test. test目录一般在相应组件目录下方，命名为`function_test.go`, 内部函数Pattern为`TestSomeFunction(t *testing.t)`签名的即为单元测试函数；
- 集成测试则使用`/example`下的yaml文件进行测试，并且通过docker, log等观察是否运行正常；
- 回归测试依靠CI进行；
- log 是进行Debug的主要工具，如果使用`cuberoot`启动Cubernetes，log会被重定向到`/var/log`，同时通过sync映射到本地便于查看。

###  新Feature开发流程

- 分工，如果和其他模块有交互部分需要讨论接口，涉及多个模块的需要会议讨论
  - 需要讨论接口的情形：如何向`iptables utils`发起创建规则链调用
  - 需要腾讯会议的情形：Serverless组件之间的进程间通信使用消息队列还是HTTP
  - 需要通过Git Wiki的情形：预研GPU CUDA程序的成员将Note记录在Wiki上
- 从Develop分支创建Branch并且推送至远程存储库
- 编写Feature所需的组件或者修改原有组件完成新Feature
- 编写单元测试，检查测试是否通过
- 编写对应的集成测试用例，在JCloud上部署并进行集成测试
- 推送到Gitee，Pull request, Git message 应该遵守[Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/), 即`<type>: <subject>`. 
  - 如果提交比较复杂，还要提交一份含有Body的Git message.
  - `<type>`应为`feat, fix, docs, refactor, test, chore, style`中的一种
- 进行code review, 然后合入Develop分支，Develop分支则需要经过完整的回归集成测试后才能合入主线

### 迭代流程

- 基础架构迭代
- 端到端迭代
- 高级功能迭代

## 功能介绍

### Init

```shell
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
./build/cubectl get gpuJobs
./build/cubectl describe GpuJob
```

### Svc

这部分将和RS, DNS一起检查，先检查SVC与负载均衡，再检查杀死Pod之后的行为；最后检查DNS.

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

# 删除Pod尝试
./build/cubectl get rs
./build/cubectl describe rs x
./build/cubectl delete pod
# Waiting...
./build/cubectl describe rs
```

### AS

```shell
# Keep scaling
./build/cubectl apply -f ./example/yaml/test-autoscaler.yaml

./build/cubectl get autoscaler 
./build/cubectl describe autoscaler <autoscaler_id>

# access all replica
./build/cubectl apply -f ./example/yaml/test-autoscaler-svc.yaml
./build/cubectl describe svc <svc_id>

curl <svc_cluster_ip>:8086/hostname # repeat to see different hostname

# repeat kill stress process to scale down
curl <svc_cluster_ip>:8086/kill-stress

# waiting...
./build/cubectl describe autoscaler
```

演示时使用的配置文件使用了默认的伸缩速度（最小间隔20秒），下面一组例子演示了更慢的伸缩速度（最小间隔80秒），可以对比体现“通过配置文件改变扩缩容速度”的功能点。

```shell
# Keep scaling (but slower)
./build/cubectl apply -f ./example/yaml/slow-autoscaler.yaml

```

### 容错

```shell
./build/cuberoot stop
# 检查Pod IP和Service
curl podIP
curl serviceIP

./build/cuberoot start
# 检查Pod IP和Service
curl podIP
curl serviceIP
```

### Schedule

这部分将运行2个符合仅符合一个node标签的Pod, 观察他们是否被部署到了唯一符合的node上。同时作为容错恢复后的检查。

```shell
./build/cubectl apply -f ./example/yaml/presentation/onemoreRS.yaml

./build/cubectl describe pod
# Only stay in worker 1
```

### Serverless

需要保证所有机器上拥有镜像：`yiyanleee/python-runtime:v1.5` & `yiyanleee/serverless-gateway:v1`

```shell
./build/cuberoot serverless enable
./build/cubectl apply -f example/serverless/hello/hello.yaml
./build/cubectl apply -f example/serverless/hello/ingress.yaml

curl "172.16.0.0:6810/hello?name=serverless" # Change to service IP
```

### Finally

```shell
./build/cuberoot reset
./build/cuberoot stop
bash ./scripts/clear.sh
```



## 附件1 CI Yaml

```yaml
name: Go

on: [pull_request, push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.18

      - name: Build
        run: make

      - name: Run go vet
        run: go vet ./...

      - name: Install staticcheck
        run: go install honnef.co/go/tools/cmd/staticcheck@latest

      - name: Run staticcheck
        run: staticcheck ./...

      - name: Go Test
        run: go test -v ./...
```

