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
	"fmt"
	"github.com/brokercap/Bifrost/manager/v2/common"
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
	"github.com/brokercap/Bifrost/server"
	"github.com/brokercap/Bifrost/server/count"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type FlowController struct {
	Ctx iris.Context
}

func (c *FlowController) GetIndex() mvc.Result {
	dbname := c.Ctx.FormValue("dbname")
	schema := c.Ctx.FormValue("Schema")
	tableName := c.Ctx.FormValue("table_name")
	channelId := c.Ctx.FormValue("channelid")

	return mvc.View{
		Name: "flow.html",
		Data: datamodles.NewFlowIndex("Flow-Bifrost", dbname, schema, tableName, channelId),
	}
}

func (c *FlowController) PostGet() interface{} {
	dbname := c.Ctx.FormValue("dbname")
	schema := c.Ctx.FormValue("Schema")
	tableName := c.Ctx.FormValue("table_name")
	channelId := c.Ctx.FormValue("channelid")
	flowType := c.Ctx.FormValue("type")

	if flowType == "" {
		flowType = "minute"
	}

	schema0 := common.TansferSchemaName(schema)
	tablename0 := common.TansferTableName(tableName)
	dbANdTableName := server.GetSchemaAndTableJoin(schema0, tablename0)

	var data []count.CountContent
	switch flowType {
	case "minute":
		data, _ = getFlowCount(&dbname, &dbANdTableName, &channelId, "Minute")
		break
	case "tenminute":
		data, _ = getFlowCount(&dbname, &dbANdTableName, &channelId, "TenMinute")
		break
	case "hour":
		data, _ = getFlowCount(&dbname, &dbANdTableName, &channelId, "Hour")
		break
	case "eighthour":
		data, _ = getFlowCount(&dbname, &dbANdTableName, &channelId, "EightHour")
		break
	case "day":
		data, _ = getFlowCount(&dbname, &dbANdTableName, &channelId, "Day")
		break
	default:
		data = make([]count.CountContent, 0)
		break
	}
	return data
}

func getFlowCount(dbname *string, dbANdTableName *string, channelId *string, FlowType string) ([]count.CountContent, error) {
	if *dbname == "" {
		return count.GetFlowAll(FlowType), nil
	}
	if *dbANdTableName != server.GetSchemaAndTableJoin("", "") {
		if *dbname == "" {
			return make([]count.CountContent, 0), fmt.Errorf("param error")
		}
		return count.GetFlowByTable(*dbname, *dbANdTableName, FlowType), nil
	}

	if *channelId != "" {
		return count.GetFlowByChannel(*dbname, *channelId, FlowType), nil
	}
	return count.GetFlowByDb(*dbname, FlowType), nil
}
