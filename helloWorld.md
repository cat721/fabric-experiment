# 智能合约的开发部署

##1 第一个智能合约 hello world

首先，我们来学习如何书写一个最简单的智能合约。在fabric中，智能合约对于链的操作有两种形式，读数据和写数据，我们下面学习一下如何对链进行操作。

```go
package main
 
import (
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)
 
type Helloworld struct {
 
}
 
func (t * Helloworld) Init(stub shim.ChaincodeStubInterface) peer.Response{
 
    args:= stub.GetStringArgs() //获得输入的参数
 
    err := stub.PutState(args[0],[]byte(args[1])) //将输入的参数存入链中
 
    if err != nil {
        shim.Error(err.Error()) //如果失败返回错误信息
    }
 
    return shim.Success(nil) //如果成功返回nil
}
 
func (t *Helloworld) Invoke (stub shim.ChaincodeStubInterface) peer.Response{
 
    fn, args := stub.GetFunctionAndParameters() //获取调用的函数和函数参数
 
    if fn =="set" {
        return t.set(stub, args)
    }else if fn == "get"{
        return t.get(stub , args)
    }
    return shim.Error("Invoke fn error")
}
 
func (t *Helloworld) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
    err := stub.PutState(args[0],[]byte(args[1])) //将数据存入链上
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(nil)
}
 
func (t *Helloworld) get (stub shim.ChaincodeStubInterface, args [] string) peer.Response{
 
    value, err := stub.GetState(args[0]) //获得数据
 
    if err != nil {
        return shim.Error(err.Error())
    }
 
    return shim.Success(value) //成功，返回获取的值
}
 
func main(){
    err := shim.Start(new(Helloworld))
    if err != nil {
        fmt.Println("start error")
    }
}
```

首先,智能合约中应该至少有两个共有函数`Init`和`Invoke`，分别负责在智能合约初始化数据以及存放和链进行交互的函数。
在上述的例子中分别实现了对于链上信息的写入和读取。fabric中的chaincode是用go写成的，并且关于更多对于链的操作，点击[这里](https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/interfaces.go)。

##2 智能合约的部署

1. 将智能合约的代码拷贝到`fabric/examples/e2e_cli/examples/chaincode/go/`。

```bash
mkdir ~/go/src/github.com/hyperledger/fabric/examples/e2e_cli/examples/chaincode/go/hello

cp hello.go ~/go/src/github.com/hyperledger/fabric/examples/e2e_cli/examples/chaincode/go/hello
```

2. 启动fabric 网络。

```bash
cd ~/go/src/github.com/hyperledger/fabric/examples/e2e_cli/examples

./network_setup.sh up
```

3. 安装chaincode

```bash
docker exec -it cli bash

peer chaincode install -n mycc -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/hello
```

4. 实例化chaincode
```bash
    peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile $ORDERER_CA -n mycc -v 0 -c '{"Args":["str","HelloWorld"]}' -C mychannel
```

5.查询chaincode 
```bash
peer chaincode query -n mycc -c '{"Args":["get","str"]}' -C myc
```

6.发起交易
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile $ORDERER_CA -C mychannel -n mycc -c '{"Args":["set","str","newHelloWorld"]}'
```

7.再次查询
```bash
'{"Args":["set","str","newHelloWorld"]}'
```