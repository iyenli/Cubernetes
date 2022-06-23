# Cubernetes
Course project for SE3356. 

## Architecture

<img src="docs/Cubernetes验收报告.assets/overview.png" alt="overview" style="zoom: 25%;" />

![serverless1](docs\Cubernetes验收报告.assets\serverless1.png)

<img src="docs/Cubernetes验收报告.assets/serverless2.png" alt="serverless2" style="zoom:28%;" />

## Quick start

Requirements: 

- Etcd(suggested path: /usr/local/bin)
- Nginx
- Resolvconf
- Kafka

Enjoy Cubernetes!

```shell
./build/cuberoot init -f ./example/yaml/master-node.yaml
./build/cuberoot join $master_ip -f ./example/yaml/slave-node.yaml
```
