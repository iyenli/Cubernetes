# 验收报告

## 架构介绍



## 成员分工



## 工程实践

Gitee 地址: 

- Cubernetes https://gitee.com/k9-s/Cubernetes
- Serverless Python框架 https://gitee.com/k9-s/Serverless-Python

项目分支介绍:

- master 项目主分支，长期存在。仅由develop branch通过pull request合入，合并时机为完成了一个完整的，经过集成测试的主要功能之后；
- develop 项目开发分支，长期存在。仅由功能分支通过pull request合入，合并时机为完成了一个经过测试的功能点之后；
- 其他功能分支，有明确的生命周期。在开启一项独立的任务之后，由任务负责开发者创建，将在任务开发完毕后合入develop分支(或者使用此分支feature的其他分支)并且删除。

CI/CD介绍：



软件测试方法介绍：

单元测试使用go test. 对于独立有明确接口的组件，如Kafka go client, etcd go client 使用go tes进行单元测试。

新Feature开发流程介绍：



## 功能介绍



