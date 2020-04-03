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
