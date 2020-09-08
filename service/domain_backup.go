package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type ServiceSetup struct {
	ChaincodeID string
	Client *channel.Client
}

func registerEvent(client *channel.Client, chaincodeID string, eventID string) (fab.Registration,<-chan *fab.CCEvent){
	reg,notifier,err := client.RegisterChaincodeEvent(chaincodeID,eventID)
	if err!=nil{
		fmt.Println("注册链码事件失败")
	}
	return reg,notifier
}

func evenResult(notifier <-chan *fab.CCEvent,eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Println("接收到链码事件: %+v",ccEvent)
	case <-time.After(time.Second *20):
		return fmt.Errorf("不能根据指定的事件id接收到相应的链码事件(%s)",eventID)
	}
	return nil
}