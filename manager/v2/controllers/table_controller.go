/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controllers

import (
	"encoding/json"
	"github.com/brokercap/Bifrost/manager/v2/common"
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
	pluginStorage "github.com/brokercap/Bifrost/plugin/storage"
	"github.com/brokercap/Bifrost/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
	"strings"
)

func TableController(app *mvc.Application) {
	app.Router.Post("/del", hero.Handler(tableDel))
	app.Router.Post("/add", hero.Handler(tableAdd))
	app.Router.Get("/toserverlist", hero.Handler(toServerList))
	app.Router.Post("/deltoserver", hero.Handler(delToServer))
	app.Router.Post("/addtoserver", hero.Handler(addToServer))
	app.Router.Post("/toserver/deal", hero.Handler(dealToServer))
}

type RequestBody struct {
	DbName     string
	TableName  string
	Schema     string
	ChannelId  string
	ToServerId string
	Index      string
}

func buildRequestBody(ctx iris.Context) *RequestBody {
	return &RequestBody{
		DbName:    ctx.FormValue("dbname"),
		TableName: ctx.FormValue("table_name"),
		Schema:    ctx.FormValue("schema_name"),
		ChannelId: ctx.FormValue("channelid"),
	}
}

func tableDel(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)

	err := server.DelTable(requestBody.DbName, schema, tableName)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	return datamodles.GenSuccessMsg("success")
}

func tableAdd(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	channelId := common.ParseIntDef(requestBody.ChannelId, 0)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)

	err := server.AddTable(requestBody.DbName, schema, tableName, channelId)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}

func toServerList(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)

	t1 := server.GetDBObj(requestBody.DbName)
	tableObj := t1.GetTable(schema, tableName)
	//tableObj := server.GetDBObj(dbname).GetTable(Schema,tablename)
	return tableObj.ToServerList
}

func dealToServer(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)
	toServerId := common.ParseIntDef(requestBody.ToServerId, 0)
	index := common.ParseIntDef(requestBody.Index, 0)

	toServerInfo := server.GetDBObj(requestBody.DbName).GetTable(schema, tableName).ToServerList[index]
	if toServerInfo.ToServerID == toServerId {
		toServerInfo.DealWaitError()
	}
	return datamodles.GenSuccessMsg("success")
}

func addToServer(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)

	toServerKey := ctx.FormValue("toserver_key")
	PluginName := ctx.FormValue("plugin_name")
	FieldListString := ctx.FormValue("fieldlist")
	MustBeSuccess := ctx.FormValue("mustbe")
	FilterQuery := ctx.FormValue("FilterQuery")
	FilterUpdate := ctx.FormValue("FilterUpdate")
	param := ctx.FormValue("param")

	var pluginParam map[string]interface{}
	err := json.Unmarshal([]byte(param), &pluginParam)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}

	if pluginStorage.GetToServerInfo(toServerKey) == nil {
		return datamodles.GenFailedMsg(toServerKey + "not exist")
	}

	fileList := make([]string, 0)
	if FieldListString != "" {
		for _, fieldName := range strings.Split(FieldListString, ",") {
			fileList = append(fileList, fieldName)
		}
	}

	var MustBeSuccessBool bool = false
	if MustBeSuccess == "1" || MustBeSuccess == "true" {
		MustBeSuccessBool = true
	}

	var FilterQueryBool bool = false
	if FilterQuery == "1" || FilterQuery == "true" {
		FilterQueryBool = true
	}

	var FilterUpdateBool bool = false
	if FilterUpdate == "1" || FilterUpdate == "true" {
		FilterUpdateBool = true
	}

	toServer := &server.ToServer{
		MustBeSuccess:  MustBeSuccessBool,
		FilterQuery:    FilterQueryBool,
		FilterUpdate:   FilterUpdateBool,
		ToServerKey:    toServerKey,
		PluginName:     PluginName,
		FieldList:      fileList,
		BinlogFileNum:  0,
		BinlogPosition: 0,
		PluginParam:    pluginParam,
	}
	dbObj := server.GetDBObj(requestBody.DbName)
	r, toServerId := dbObj.AddTableToServer(schema, tableName, toServer)
	if !r {
		return datamodles.GenFailedMsg("unkown error")
	}
	defer server.SaveDBConfigInfo()
	return datamodles.GenResultData(true, "success", toServerId)
}

func delToServer(ctx iris.Context) interface{} {
	requestBody := buildRequestBody(ctx)
	schema := common.TansferSchemaName(requestBody.Schema)
	tableName := common.TansferTableName(requestBody.TableName)
	toServerId := common.ParseIntDef(requestBody.ToServerId, 0)

	server.GetDBObj(requestBody.DbName).DelTableToServer(schema, tableName, toServerId)
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}
