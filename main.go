package main

import (
	"hauturier/sdkInit"
	"hauturier/service"
	"hauturier/web"
	"hauturier/web/controller"
)
import "fmt"
import "os"

const (
	configFile = "config.yaml"
	initialized = false
)

func main(){
	initInfo := &sdkInit.InitInfo{
		ChannelID : "mychannel",
		ChannelConfig : os.Getenv("GOPATH") + "/src/github.com/hauturier.com/hauturier/channel-artifacts/channel.tx",
		OrgAdmin : "Admin",
		OrgName : "Org1",
		OrdererOrgName : "orderer.hauturier.com",

		ChaincodeID: "TestCC",
		ChaincodePath: "github.com/hauturier.com/hauturier/chaincode",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		UserName: "User1",
	}
	sdk,err := sdkInit.SetSDK(configFile,initialized)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk,initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient,err := sdkInit.InstallChaincode(sdk,initInfo)
	if err!= nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

	//req := channel.Request{ChaincodeID:"TestCC",Fcn:"query",Args:[][]byte{[]byte("hello")}}
	//response,err := channelClient.Query(req,channel.WithRetry(retry.DefaultChannelOpts))
	//if err!=nil {
	//	fmt.Println(err)
	//}else {
	//	fmt.Println(response.Payload)
	//}
	serviceSetup := service.ServiceSetup{
		ChaincodeID:"TestCC",
		Client:channelClient,
	}

	msg,err := serviceSetup.SetInfo("hello","motherfuck")
	if err!=nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}
	msg,err = serviceSetup.QueryInfo("hello")
	if err!=nil {
		fmt.Println(err)
	}else {
		fmt.Println(msg)
	}

	app := controller.Application{
		Fabric:&serviceSetup,
	}
	web.WebStart(&app)
}

