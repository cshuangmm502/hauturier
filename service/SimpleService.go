package service

import "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

func (t *ServiceSetup)QueryInfo(key string) (string,error) {
	req := channel.Request{ChaincodeID:t.ChaincodeID,Fcn:"query",Args:[][]byte{[]byte(key)}}
	response,err := t.Client.Query(req)
	if err!=nil {
		return "",err
	}
	return string(response.Payload),nil
}

func (t *ServiceSetup)SetInfo(key,value string) (string,error){
	eventID := "eventInfo"
	reg,notifier := registerEvent(t.Client,t.ChaincodeID,eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID:t.ChaincodeID,Fcn:"set",Args:[][]byte{[]byte(key),[]byte(value),[]byte(eventID)}}
	response,err :=t.Client.Execute(req)
	if err!=nil{
		return "",err
	}
	err = evenResult(notifier,eventID)
	if err!=nil{
		return "",err
	}
	return string(response.TransactionID),nil
}