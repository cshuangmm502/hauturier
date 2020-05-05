package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
import "github.com/hyperledger/fabric/protos/peer"

type TestChaincoe struct {
	
}

func (t *TestChaincoe) Init(stub shim.ChaincodeStubInterface) peer.Response {
	arg :=stub.GetStringArgs()
	if len(arg)!=2 {
		return shim.Error("实例化参数仅能为2个")
	}
	err := stub.PutState(arg[0],[]byte(arg[1]))
	if err!=nil {
		return shim.Error("保存状态时出错")
	}
	return shim.Success(nil)
}

func (t *TestChaincoe) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function,args := stub.GetFunctionAndParameters()
	if function == "query"{
		return query(stub,args)
	}
	if function == "set"{
		return set(stub,args)
	}
	return shim.Error("非法函数操作")
}

func query(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args)!=1 {
		return shim.Error("查询参数只能为1")
	}
	result,err := stub.GetState(args[0])
	if err!=nil {
		return shim.Error("查询状态出错")
	}
	return shim.Success(result)
}

func set(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	if len(args)!=3 {
		return shim.Error("请输入正确的k-v键值对")
	}
	err := stub.PutState(args[0],[]byte(args[1]))
	if err!=nil {
		return shim.Error("保存状态时出错")
	}
	err = stub.SetEvent(args[2],[]byte{})
	if err!=nil {
		return shim.Error("保存状态时出错")
	}
	return shim.Success(nil)
}

func main(){
	err := shim.Start(new(TestChaincoe))
	if err != nil {
		fmt.Printf("启动链码失败",err)
	}
}
