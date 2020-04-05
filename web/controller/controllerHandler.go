package controller

import (
	"net/http"
	"hauturier/service"
)

type Application struct {
	Fabric *service.ServiceSetup
}

func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	showView(w,r,"index.html",nil)
}

func (app *Application) SetInfoView(w http.ResponseWriter, r *http.Request) {
	showView(w,r,"setInfo.html",nil)
}

//根据指定的key,修改对应的val
func (app *Application) SetInfo(w http.ResponseWriter, r *http.Request) {
	//获取提交数据
	key := r.FormValue("key")
	val := r.FormValue("value")

	//调用业务层,反序列化
	transactionID,err := app.Fabric.SetInfo(key,val)

	//封装响应数据
	data := &struct {
		Flag bool
		Msg string
	}{
		Flag:true,
		Msg:"",
	}
	if err!=nil{
		data.Msg = err.Error()
	}else {
		data.Msg = "操作成功,交易ID:"+transactionID
	}

	//响应客户端
	showView(w,r,"setInfo.html",data)
}

func (app *Application) QueryInfo(w http.ResponseWriter,r *http.Request) {
	key := r.FormValue("key")

	result,err := app.Fabric.QueryInfo(key)
	data := &struct {
		Msg string
	}{
		Msg:"",
	}
	if err!=nil{
		data.Msg = "没有查询到对应信息"
	}else{
		data.Msg = "查询成功:"+result
	}

	showView(w,r,"queryResponse.html",data)
}