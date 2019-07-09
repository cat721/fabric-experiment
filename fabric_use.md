#Fabric 使用部署
==============================

在上一个实验中，我们讲解了fabric的搭建和集成测试，本次实验具体讲解fabric的使用过程。

##1、生成公私钥和证书
==============================
Fabric中有两种类型的公私钥和证书，一种是给节点之前通讯安全而准备的TLS证书，另一种是用户登录和权限控制的用户证书。这些证书本来应该是由CA来颁发，但是我们这里是测试环境，并没有启用CA节点，所以Fabric帮我们提供了一个工具：cryptogen

###1.1、编译生成cryptogen
==============================
Fabric官方提供了专门编译cryptogen的入口，运行以下命令编译cryptogen：
```bash
cd ~/go/src/github.com/hyperledger/fabric

make cryptogen
```

如果报错：
```bash
vendor/github.com/miekg/pkcs11/pkcs11.go:29:18: fatal error: ltdl.h: No such file or directory
```
请安装

ubuntu:
```bash
sudo apt install libltdl3-dev
```

Centos:
```bash
yum install libtool-ltdl-devel
```

运行后系统返回结果：
```bash
build/bin/cryptogen 
CGO_CFLAGS=" " GOBIN=/home/studyzy/go/src/github.com/hyperledger/fabric/build/bin go install -tags "" -ldflags "-X github.com/hyperledger/fabric/common/tools/cryptogen/metadata.Version=1.0.0" github.com/hyperledger/fabric/common/tools/cryptogen 
Binary available as build/bin/cryptogen
```

cryptogen程序在build/bin文件夹下。

###1.2、配置crypto-config.yaml
==============================
examples/e2e_cli/crypto-config.yaml已经提供了一个Orderer Org和两个Peer Org的配置，该模板中也对字段进行了注释。

```yaml
- Name: Org2 
  Domain: org2.example.com 
  Template: 
    Count: 2 
  Users: 
    Count: 1
```

Name和Domain就是关于这个组织的名字和域名，这主要是用于生成证书的时候，证书内会包含该信息。而Template Count=2说明生成2套公私钥和证书，一套是peer0.org2的，还有一套是peer1.org2的。最后Users. Count=1是说每个Template下面会有几个普通User（注意，Admin是Admin，不包含在这个计数中），这里配置了1，也就是需要一个普通用户User1@org2.example.com 可以根据实际需要调整这个配置文件，增删Org Users等。

###1.3、生成公私钥和证书
==============================

用cryptogen去读取该crypto-config.yaml，并生成对应的公私钥和证书：
```bash
cd examples/e2e_cli/
 
../../build/bin/cryptogen generate --config=./crypto-config.yaml
```
生成的文件保存到crypto-config文件夹中。

##2、生成创世区块和Channel配置区块
================================

###2.1、编译生成configtxgen
==============================
通过make命令生成configtxgen程序：
 ```bash
 cd ~/go/src/github.com/hyperledger/fabric
 
 make configtxgen
 ```
###2.2、配置configtx.yaml
==============================

examples/e2e_cli/configtx.yaml这个文件里面配置了由2个Org参与的Orderer共识配置TwoOrgsOrdererGenesis，以及由2个Org参与的Channel配置：TwoOrgsChannel。
Orderer可以设置共识的算法是Solo还是Kafka，以及共识时区块大小，超时时间等，本实验使用默认值即可，不用更改。Peer节点的配置包含了MSP的配置，锚节点的配置。
如果有更多的Org，或者有更多的Channel，可以根据模板进行对应的修改。

###2.3、生成创世区块
==============================
配置修改好后，我们就用configtxgen 生成创世区块。
并把这个区块保存到本地channel-artifacts文件夹中：
```bash
cd examples/e2e_cli/

../../build/bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```
锚节点的更新信息：
```bash
../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

../../build/bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
```

###2.4、生成Channel配置区块
==============================
```bash
../../build/bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
```
 
 最终，我们在channel-artifacts文件夹中有4个文件。
```bash
channel-artifacts/ 
├── channel.tx 
├── genesis.block 
├── Org1MSPanchors.tx 
└── Org2MSPanchors.tx
``` 
 
##3、配置Fabric环境的docker-compose文件
=======================================
配置docker-compose的yaml文件，启动Fabric的Docker环境。

###3.1、配置Orderer
==============================
Orderer的配置是在base/docker-compose-base.yaml里：
```dockerfile
orderer.example.com: 
  container_name: orderer.example.com 
  image: hyperledger/fabric-orderer 
  environment: 
    - ORDERER_GENERAL_LOGLEVEL=debug 
    - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0 
    - ORDERER_GENERAL_GENESISMETHOD=file 
    - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block 
    - ORDERER_GENERAL_LOCALMSPID=OrdererMSP 
     - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp 
     # enabled TLS 
    - ORDERER_GENERAL_TLS_ENABLED=true 
    - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key 
    - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt 
    - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt] 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric 
  command: orderer 
  volumes: 
  - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block 
  - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp 
  - ../crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/:/var/hyperledger/orderer/tls 
  ports: 
    - 7050:7050
```
其中``ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block``，这个创世区块就是我们之前创建的创世区块，
  ``- ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block``是Host到Docker的映射。
另外的配置主要是TL，Log等，最后暴露出服务端口7050。

###3.2、配置Peer
==============================
Peer的配置是在base/docker-compose-base.yaml和peer-base.yaml里面，摘取其中的peer0.org1：
```dockerfile
peer-base: 
  image: hyperledger/fabric-peer 
  environment: 
    - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock 
    # the following setting starts chaincode containers on the same 
    # bridge network as the peers 
    # https://docs.docker.com/compose/networking/ 
    - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=e2ecli_default 
    #- CORE_LOGGING_LEVEL=ERROR 
    - CORE_LOGGING_LEVEL=DEBUG 
    - CORE_PEER_TLS_ENABLED=true 
    - CORE_PEER_GOSSIP_USELEADERELECTION=true 
    - CORE_PEER_GOSSIP_ORGLEADER=false 
    - CORE_PEER_PROFILE_ENABLED=true 
    - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt 
    - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key 
    - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer 
  command: peer node start

peer0.org1.example.com: 
  container_name: peer0.org1.example.com 
  extends: 
    file: peer-base.yaml 
    service: peer-base 
  environment: 
    - CORE_PEER_ID=peer0.org1.example.com 
    - CORE_PEER_ADDRESS=peer0.org1.example.com:7051 
    - CORE_PEER_CHAINCODELISTENADDRESS=peer0.org1.example.com:7052 
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051 
    - CORE_PEER_LOCALMSPID=Org1MSP 
  volumes: 
      - /var/run/:/host/var/run/ 
      - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/fabric/msp 
      - ../crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls:/etc/hyperledger/fabric/tls 
  ports: 
    - 7051:7051 
    - 7052:7052 
    - 7053:7053
```
在Peer的配置中，主要是给Peer分配好各种服务的地址，以及TLS和MSP信息。

###3.3、配置CLI
==============================
CLI在整个Fabric网络中扮演客户端的角色，我们在开发测试的时候可以用CLI来代替SDK，
执行各种SDK能执行的操作。CLI会和Peer相连，把指令发送给对应的Peer执行。
CLI的配置在docker-compose-cli.yaml中：
```dockerfile
cli: 
  container_name: cli 
  image: hyperledger/fabric-tools 
  tty: true 
  environment: 
    - GOPATH=/opt/gopath 
    - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock 
    - CORE_LOGGING_LEVEL=DEBUG 
    - CORE_PEER_ID=cli 
    - CORE_PEER_ADDRESS=peer0.org1.example.com:7051 
    - CORE_PEER_LOCALMSPID=Org1MSP 
    - CORE_PEER_TLS_ENABLED=true 
    - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt 
    - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key 
    - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
    - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
  working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer 
#  command: /bin/bash -c './scripts/script.sh ${CHANNEL_NAME}; sleep $TIMEOUT' 
  volumes: 
      - /var/run/:/host/var/run/ 
      - ../chaincode/go/:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go 
      - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ 
       - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/ 
      - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts 
  depends_on: 
    - orderer.example.com 
    - peer0.org1.example.com 
    - peer1.org1.example.com 
    - peer0.org2.example.com 
    - peer1.org2.example.com
```
，CLI启动的时候默认连接的是`peer0.org1.example.com`，并且启用了TLS。默认
以Admin@org1.example.com连接到Peer的。CLI启动的时候，执行`./scripts/script.sh`
脚本，这个脚本即`fabric/examples/e2e_cli/scripts/script.sh`，其完成了Fabric环境
的初始化和ChainCode的安装及运行。在文件映射配置上，`../chaincode/go/:/opt/gopath/src/github.com/hyperledger/fabric/examples/chaincode/go`，
也就是说需要安装的ChainCode都是在`fabric/examples/chaincode/go`
目录下，以后我们要开发自己的ChainCode，只需要把我们的代码复制到该目录即可。

##4、初始化Fabric环境
==============================

###4.1、启动Fabric环境的容器
==============================

```bash
docker-compose -f docker-compose-cli.yaml up -d
```
运行``docker ps``命令可以看启动的结果：

```textmate
CONTAINER ID        IMAGE                        COMMAND             CREATED             STATUS              PORTS                                                                       NAMES 
6f98f57714b5        hyperledger/fabric-tools     "/bin/bash"         8 seconds ago       Up 7 seconds                                                                                    cli 
6e7b3fd0e803        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 8 seconds        0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com 
9e67abfb982f        hyperledger/fabric-orderer   "orderer"           11 seconds ago      Up 8 seconds        0.0.0.0:7050->7050/tcp                                                      orderer.example.com 
908d7fe2a4c7        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 9 seconds        0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com 
6bb187ac10ec        hyperledger/fabric-peer      "peer node start"   11 seconds ago      Up 10 seconds       0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com 
150baa520ed0        hyperledger/fabric-peer      "peer node start"   12 seconds ago      Up 9 seconds        0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com
```
###4.2、创建Channel
==============================
要进入cli容器内部:
```bash
docker exec -it cli bash\
```
创建Channel的命令是peer channel create，我们前面创建2.4创建Channel的配置区块时，指定了Channel的名字是mychannel，
这里我们必须创建同样名字的Channel。

```bash
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile $ORDERER_CA
```
执行该命令后，系统会提示：
```text
[channelCmd] readBlock -> DEBU 020 Received block:0
```
系统会在cli内部的当前目录创建一个mychannel.block文件，其他节点要加入这个Channel
必须使用这个文件。

###4.3、各个Peer加入Channel
==============================
CLI默认连接的是peer0.org1，将这个Peer加入mychannel只需要运行如下命令:
```bash
peer channel join -b mychannel.block
```
系统返回消息：
```text
[channelCmd] executeJoin -> INFO 006 Peer joined the channel!
```
修改CLI的环境变量，使其指向另外的Peer。例如、将peer1.org1加入mychannel：
```bash
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer1.org1.example.com:7051

peer channel join -b mychannel.block
```
同样的方法将其他两个节点加入到channel中：
```bash
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel join -b mychannel.block

CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer1.org2.example.com:7051

peer channel join -b mychannel.block
```

###4.4、更新锚节点
==============================
对于Org1来说，peer0.org1是锚节点，需要连接上它并更新锚节点：
```bash
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls true --cafile $ORDERER_CA
```
另外对于Org2，peer0.org2是锚节点:
```bash
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx --tls true --cafile $ORDERER_CA
```
#5、链上代码的安装与运行
==============================
以上，整个Fabric网络和Channel都准备完毕，接下来我们来安装和运行ChainCode。
这个例子实现了a，b两个账户，相互之间可以转账。
##5.1、Install ChainCode安装链上代码
====================================
链上代码的安装需要在各个相关的Peer上进行，对于我们现在这种Fabric网络，
如果4个Peer都想对chaincode进行操作，那么就需要安装4次。

保持在CLI的命令行下，先切换到peer0.org1节点：
```bash
CORE_PEER_LOCALMSPID="Org1MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
```
使用peer chaincode install命令可以安装指定的ChainCode并对其命名：

```bash
peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
```
安装的过程其实就是对CLI中指定的代码进行编译打包，并把打包好的文件发送到Peer，等待接下来的实例化。

##5.2、Instantiate ChainCode实例化链上代码
==========================================
实例化链上代码主要是在Peer所在的机器上对前面安装好的链上代码进行包装，生成对应Channel的Docker镜像和Docker容器。
并且在实例化时我们可以指定背书策略。我们运行以下命令完成实例化：
```bash
peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n mycc -v 1.0 -c '{"Args":["init","a","100","b","200"]}' -P "OR      ('Org1MSP.member','Org2MSP.member')"
```
回到Ubuntu终端，使用docker ps可以看到有新的容器正在运行：

```text
CONTAINER ID        IMAGE                                 COMMAND                  CREATED              STATUS              PORTS                                                                       NAMES 
07791d4a99b7        dev-peer0.org1.example.com-mycc-1.0   "chaincode -peer.a..."   About a minute ago   Up About a minute                                                                               dev-peer0.org1.example.com-mycc-1.0 
6f98f57714b5        hyperledger/fabric-tools              "/bin/bash"              About an hour ago    Up About an hour                                                                                cli 
6e7b3fd0e803        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com
9e67abfb982f        hyperledger/fabric-orderer            "orderer"                About an hour ago    Up About an hour    0.0.0.0:7050->7050/tcp                                                      orderer.example.com 
908d7fe2a4c7        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com 
6bb187ac10ec        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com 
150baa520ed0        hyperledger/fabric-peer               "peer node start"        About an hour ago    Up About an hour    0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com
```
##5.3、在一个Peer上查询并发起交易
===============================
调用ChainCode的查询代码,在cli容器内执行:

```bash
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
```
返回结果：Query Result: 100.

把a账户的10元转给b。对应的代码：
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n mycc -c '{"Args":["invoke","a","b","10"]}'
```
##5.4、在另一个节点上查询交易  
==============================

在peer0.org2安装链上代码：      
```bash
CORE_PEER_LOCALMSPID="Org2MSP" 
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp 
CORE_PEER_ADDRESS=peer0.org2.example.com:7051

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02
```
运行查询命令：

```bash
peer chaincode query -C mychannel -n mycc -c '{"Args":["query","a"]}'
```
需要等待一定时间才返回结果：

Query Result: 90

因为peer0.org2也需要生成Docker镜像，创建对应的容器，才能通过容器返回结果。我们回到Ubuntu终端，执行`docker ps`，
可以看到又多了一个容器：

```text
CONTAINER ID        IMAGE                                 COMMAND                  CREATED             STATUS              PORTS                                                                       NAMES 
3e37aba50189        dev-peer0.org2.example.com-mycc-1.0   "chaincode -peer.a..."   2 minutes ago       Up 2 minutes                                                                                    dev-peer0.org2.example.com-mycc-1.0 
07791d4a99b7        dev-peer0.org1.example.com-mycc-1.0   "chaincode -peer.a..."   21 minutes ago      Up 21 minutes                                                                                   dev-peer0.org1.example.com-mycc-1.0 
6f98f57714b5        hyperledger/fabric-tools              "/bin/bash"              About an hour ago   Up About an hour                                                                                cli 
6e7b3fd0e803        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:10051->7051/tcp, 0.0.0.0:10052->7052/tcp, 0.0.0.0:10053->7053/tcp   peer1.org2.example.com
9e67abfb982f        hyperledger/fabric-orderer            "orderer"                About an hour ago   Up About an hour    0.0.0.0:7050->7050/tcp                                                      orderer.example.com 
908d7fe2a4c7        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:7051-7053->7051-7053/tcp                                            peer0.org1.example.com 
6bb187ac10ec        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:9051->7051/tcp, 0.0.0.0:9052->7052/tcp, 0.0.0.0:9053->7053/tcp      peer0.org2.example.com 
150baa520ed0        hyperledger/fabric-peer               "peer node start"        About an hour ago   Up About an hour    0.0.0.0:8051->7051/tcp, 0.0.0.0:8052->7052/tcp, 0.0.0.0:8053->7053/tcp      peer1.org1.example.com
```

   

   




