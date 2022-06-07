# Cubernetes
Course project for SE3356.

> - 在首次启动，需要下载weave, 可能需要20-30秒。尽量不要在这时候apply api object:)
> - 测试时可以用Script下的clear.sh清理残留的docker容器。[风险提示: Docker和IPtables会被清空]

## Quick start

在开始之前，您需要安装ETCD, 放到合适的目录下(Suggest: `/usr/local/bin`). 建议安装Nginx以获取默认的配置作为Volume. 

```shell
apt-get install resolvconf
```

此外，需要安装Kafka. (在Master)

```shell
wget https://dlcdn.apache.org/kafka/3.2.0/kafka_2.13-3.2.0.tgz
tar -xzvf kafka_2.13-3.2.0.tgz
cp -r ./kafka_2.13-3.2.0/* /usr/local/kafka/
# copy config *service to 
vim /etc/systemd/system/kafka.service
vim /etc/systemd/system/zookeeper.service

systemctl daemon-reload
systemctl restart zookeeper
systemctl status zookeeper
systemctl restart kafka
systemctl status kafka

systemctl enable kafka
systemctl enable zookeeper
```

## Test Routines

测试用例在./example/presentation下。



