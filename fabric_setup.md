
#fabric 的搭建

##1、环境准备

###1.1、实验环境

- Ubuntu 16.04

- fabric 1.0 release

- go 1.12
###1.2、安装docker

```bash
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.daocloud.io/docker/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=$(dpkg --print-architecture)] https://download.daocloud.io/docker/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install -y -q docker-ce=*
sudo service docker start
sudo service docker status
```
其他环境的docker安装请参考[这里](https://download.daocloud.io/Docker_Mirror/Docker)。

如果ubuntu的源是国外源，请先更换国内源，如何更换，点击[这里](https://www.linuxidc.com/Linux/2017-11/148627.htm)。

查看docker版本。
```bash
docker version
```

将用户加入到docker工作组。
```bash
sudo usermod -aG docker ${USERNAME}
```

设置docker镜像加速。
```bash
curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://f1361db2.m.daocloud.io
```

重启docker。
```bash
sudo systemctl restart docker
```

查看docker是否启用
```bash
docker ps
```

 如果docker 没有成功启动，重启设备。
 
### 1.3、安装docker-comppose.
 
 ```bash
 curl -L https://get.daocloud.io/docker/compose/releases/download/1.24.1/docker-compose-`uname -s`-`uname -m` > ./docker-compose
 
 sudo cp docker-compose /usr/local/bin/
 
 sudo chmod +x /usr/local/bin/docker-compose
 ```
 查看docke-compose版本 
 ```bash
 docker-compose version
 ```
 
 ###1.4、安装golang
```bash
wget https://storage.googleapis.com/golang/go1.12.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go1.12.linux-amd64.tar.gz
```

编辑当前用户的环境变量：
```bash
vi ~/.profile
```

添加如下内容：
```text
export PATH=$PATH:/usr/local/go/bin 
export GOROOT=/usr/local/go 
export GOPATH=$HOME/go 
export PATH=$PATH:$HOME/go/bin
```
编辑保存并退出vi后，环境载入，查看golang版本：
```bash
source ~/.profile

go version
```

 
 ##2、安装Fabric
 
 ###2.1、Fabric源码下载
 
 ```bash
mkdir -p ~/go/src/github.com/hyperledger 

cd ~/go/src/github.com/hyperledger  

git clone https://github.com/hyperledger/fabric.git
 ```
 
由于Fabric一直在更新，所有我们并不需要最新最新的源码，需要切换到v1.0.0版本的源码即可：

```bash
git checkout release-1.0
```

###2.2、Fabric Docker镜像的下载

通过官方提供的脚本下载docker镜像

```bash
cd ~/fabric/examples/e2e_cli

bash download-dockerimages.sh -c x86_64-1.1.0 -f x86_64-1.1.0
``` 
下载完毕后，我们运行以下命令检查下载的镜像列表：
```bash
docker images 
```

得到如下结果：
```
REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
hyperledger/fabric-tools       latest              0403fd1c72c7        24 months ago       1.32GB
hyperledger/fabric-tools       x86_64-1.0.0        0403fd1c72c7        24 months ago       1.32GB
hyperledger/fabric-couchdb     latest              2fbdbf3ab945        24 months ago       1.48GB
hyperledger/fabric-couchdb     x86_64-1.0.0        2fbdbf3ab945        24 months ago       1.48GB
hyperledger/fabric-kafka       latest              dbd3f94de4b5        24 months ago       1.3GB
hyperledger/fabric-kafka       x86_64-1.0.0        dbd3f94de4b5        24 months ago       1.3GB
hyperledger/fabric-zookeeper   latest              e545dbf1c6af        24 months ago       1.31GB
hyperledger/fabric-zookeeper   x86_64-1.0.0        e545dbf1c6af        24 months ago       1.31GB
hyperledger/fabric-orderer     latest              e317ca5638ba        24 months ago       179MB
hyperledger/fabric-orderer     x86_64-1.0.0        e317ca5638ba        24 months ago       179MB
hyperledger/fabric-peer        latest              6830dcd7b9b5        24 months ago       182MB
hyperledger/fabric-peer        x86_64-1.0.0        6830dcd7b9b5        24 months ago       182MB
hyperledger/fabric-javaenv     latest              8948126f0935        24 months ago       1.42GB
hyperledger/fabric-javaenv     x86_64-1.0.0        8948126f0935        24 months ago       1.42GB
hyperledger/fabric-ccenv       latest              7182c260a5ca        24 months ago       1.29GB
hyperledger/fabric-ccenv       x86_64-1.0.0        7182c260a5ca        24 months ago       1.29GB
hyperledger/fabric-ca          latest              a15c59ecda5b        24 months ago       238MB
hyperledger/fabric-ca          x86_64-1.0.0        a15c59ecda5b        24 months ago       238MB
```
###2.3、启动Fabric网络并完成ChainCode的测试

```bash
./network_setup.sh up
```
如果发现以下错误：
```
code = Unknown desc = Error starting container: API error (404): {"message":"network e2ecli_default not found"}
```
请手动修改 ``e2e_cli/base/peer-base.yaml``:
```dockerfile
- CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=e2e_cli_default
```

测试通过后会显示：
```
===================== All GOOD, End-2-End execution completed =====================


 _____   _   _   ____            _____   ____    _____
| ____| | \ | | |  _ \          | ____| |___ \  | ____|
|  _|   |  \| | | | | |  _____  |  _|     __) | |  _|
| |___  | |\  | | |_| | |_____| | |___   / __/  | |___
|_____| |_| \_| |____/          |_____| |_____| |_____|
```

