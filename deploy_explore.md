#安装部署fabric浏览器
=========================

##1、准备环境
=====================

- nodejs 8.11.x (Note that v9.x is not yet supported)
- PostgreSQL 9.5 or greater

###1.1 安装node 8.11.3
=====================

```bash
wget https://nodejs.org/dist/v8.11.3/node-v8.11.3-linux-x64.tar.xz
tar -xvf node-v8.11.3-linux-x64.tar.xz
```
将node的可执行文件夹写入环境变量`~/.profile`。

```bash
node -v
npm -v
```
###1.2 安装配置PostgreSQL

```bash
sudo apt-get install postgresql
```

##2、部署浏览器
=====================

###2.1 下载源码 
=====================
```bash
cd ~/go/src/github.com/hyperledger/
git clone https://github.com/hyperledger/blockchain-explorer.git
cd blockchain-explorer
checkout release-3.4
```
###2.2 初始化数据库
======================
```bash
sudo -u postgres psql

\i app/persistence/postgreSQL/db/explorerpg.sql
\i app/persistence/postgreSQL/db/updatepg.sql
```

查看数据库状态命令：
- `\l` 查看创建的数据库。 
- `\d` 查看创建的表。

##2.3 修改配置文件
========================
启动第一次实验的测试网络，根据其配置浏览器信息。配置文件位置为`blockchain-explorer/app/platform/fabric/config.json`

更改密钥信息位置：
```bash
sed -i "s|fabric-path/fabric-samples/first-network|/home/ubuntu/go/src/github.com/hyperledger/fabric-samples/balance-transfer/artifacts/channel
|g" config.json
```

##2.4 build浏览器
===================
```bash
cd blockchain-explorer
npm install
cd blockchain-explorer/app/test
npm install
npm run test
cd client/
npm install
npm test -- -u --coverage
npm run build
```

##2.5 启动网络
===============
启动fabric-sample/balance-transfer网络；
启动浏览器
```bash
node main
```