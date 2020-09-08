package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	//"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)
import "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
import "fmt"
import mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
import "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
import "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
import "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"

const ChaincodeVersion  = "1.0"

func SetSDK(ConfigFile string,initialized bool) (*fabsdk.FabricSDK, error){
	if initialized{
		return nil, fmt.Errorf("Fabric SDK已被实例化")
	}
	sdk,err := fabsdk.New(config.FromFile(ConfigFile))
	if err!=nil {
		return nil,fmt.Errorf("实例化SDK失败: %v",err)
	}
	fmt.Println("Fabric SDK实例化成功")
	return sdk,nil
}

func CreateChannel(sdk *fabsdk.FabricSDK,info *InitInfo) error {
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin),fabsdk.WithOrg(info.OrgName))
	fmt.Println("%+v",info)
	if clientContext == nil{
		return fmt.Errorf("根据指定的组织和用户创建资源管理客户端Context失败")
	}
	fmt.Println("%+v",clientContext)
	resMgmtClient,err := resmgmt.New(clientContext)
	if err !=nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v",err)
	}
	fmt.Println("tttttttttt")

	mspClient,err := mspclient.New(sdk.Context(),mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的OrgName创建Org MSP客户端实例失败: %v",err)
	}
	adminIdentity,err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err !=nil {
		return fmt.Errorf("获取指定id的签名标识失败: %v",err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID:info.ChannelID,ChannelConfigPath:info.ChannelConfig,SigningIdentities:[]msp.SigningIdentity{adminIdentity}}
	_,err = resMgmtClient.SaveChannel(req,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err !=nil {
		return fmt.Errorf("创建应用通道失败: %v",err)
	}
	fmt.Println("创建应用通道成功")

	info.OrgResMgmt = resMgmtClient

	err = info.OrgResMgmt.JoinChannel(info.ChannelID,resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err!=nil {
		return fmt.Errorf("加入通道失败%v",err)
	}
	fmt.Println("peers已加入通道")
	return nil
}


func InstallChaincode(sdk *fabsdk.FabricSDK,info *InitInfo) (*channel.Client,error) {
	fmt.Println("开始安装链码")
	ccpkg,err := gopackager.NewCCPackage(info.ChaincodePath,info.ChaincodeGoPath)
	if err!=nil {
		return nil,fmt.Errorf("创建链码包失败")
	}
	installCCReq := resmgmt.InstallCCRequest{Name:info.ChaincodeID,Path:info.ChaincodePath,Version:ChaincodeVersion,Package:ccpkg}
	_,err = info.OrgResMgmt.InstallCC(installCCReq,resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err !=nil {
		return nil,fmt.Errorf("安装链码失败")
	}
	fmt.Println("安装链码成功")
	fmt.Println("开始实例化链码")

	//指定背书策略,需要和configtx.yaml对应上?
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	instantialCCReq := resmgmt.InstantiateCCRequest{Name:info.ChaincodeID,Path:info.ChaincodePath,Version:ChaincodeVersion,Args:[][]byte{[]byte("hello"),[]byte("world")},Policy:ccPolicy}
	_,err = info.OrgResMgmt.InstantiateCC(info.ChannelID,instantialCCReq,resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err!= nil {
		return nil,fmt.Errorf("实例化链码失败")
	}
	fmt.Println("实例化链码成功")

	clientChannelContext := sdk.ChannelContext(info.ChannelID,fabsdk.WithUser(info.UserName),fabsdk.WithOrg(info.OrgName))
	channelClient,err := channel.New(clientChannelContext)
	if err!= nil {
		return nil,fmt.Errorf("创建通道管理客户端失败:%v",err)
	}
	fmt.Println("创建通道管理客户端成功,可利用该客户端查询或执行事务")
	return channelClient,nil
}
