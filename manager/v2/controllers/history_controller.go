package controllers

import (
	"encoding/json"
	"github.com/brokercap/Bifrost/manager/v2/common"
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
	"github.com/brokercap/Bifrost/server"
	"github.com/brokercap/Bifrost/server/history"
	"github.com/kataras/iris"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
)

func HistoryController(app *mvc.Application) {
	app.Router.Get("/list", hero.Handler(list))
	app.Router.Post("/add", hero.Handler(add))
	app.Router.Post("/stop", hero.Handler(stop))
	app.Router.Post("/del", hero.Handler(del))
	app.Router.Post("/start", hero.Handler(start))
	app.Router.Post("/checkwhere", hero.Handler(checkWhere))
}

func list(ctx iris.Context) hero.Result {
	dbname := ctx.FormValue("dbname")
	tableName := ctx.FormValue("table_name")
	schema := ctx.FormValue("schema_name")
	status := ctx.FormValue("status")

	historyStatus := getHistoryStatus(status)
	historyList := history.GetHistoryList(dbname, schema, tableName, historyStatus)

	if ctx.FormValue("format") == "json" {
		return hero.Response{Object: historyList}
	}

	return hero.View{
		Name: "history.list.html",
		Data: iris.Map{
			"Title":       "History List - Bifrost",
			"DbName":      dbname,
			"TableName":   tableName,
			"SchemaName":  schema,
			"HistoryList": historyList,
			"DbList":      server.GetListDb(),
			"StatusList":  allHistoryStatus(),
			"Status":      status,
		},
	}
}

func allHistoryStatus() []history.HisotryStatus {
	return []history.HisotryStatus{
		history.HISTORY_STATUS_ALL,
		history.HISTORY_STATUS_CLOSE,
		history.HISTORY_STATUS_RUNNING,
		history.HISTORY_STATUS_HALFWAY,
		history.HISTORY_STATUS_OVER,
		history.HISTORY_STATUS_KILLED}
}

func getHistoryStatus(status string) history.HisotryStatus {
	var historyStatus history.HisotryStatus
	switch status {
	case "close":
		historyStatus = history.HISTORY_STATUS_CLOSE
		break
	case "running":
		historyStatus = history.HISTORY_STATUS_RUNNING
		break
	case "over":
		historyStatus = history.HISTORY_STATUS_OVER
		break
	case "halfway":
		historyStatus = history.HISTORY_STATUS_HALFWAY
		break
	case "killed":
		historyStatus = history.HISTORY_STATUS_KILLED
		break
	default:
		historyStatus = history.HISTORY_STATUS_ALL
		break
	}
	return historyStatus
}

func add(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	tableName := ctx.FormValue("table_name")
	schema := ctx.FormValue("schema_name")
	property := ctx.FormValue("property")
	toServerIDList := ctx.FormValue("ToserverIds")

	if common.TansferTableName(tableName) == "*" {
		return datamodles.GenResultData(false, "不能给 AllTables 添加全量任务!", 0)
	}

	var historyProperty history.HistoryProperty
	err := json.Unmarshal([]byte(property), &historyProperty)
	if err != nil {
		return datamodles.GenResultData(false, err.Error(), 0)
	}

	var toServerIds []int
	err = json.Unmarshal([]byte(toServerIDList), &toServerIds)
	if err != nil {
		return datamodles.GenResultData(false, err.Error(), 0)
	}
	if len(toServerIds) == 0 {
		return datamodles.GenResultData(false, "toServerIds error", 0)
	}

	err = history.CheckWhere(dbname, schema, tableName, historyProperty.Where)
	if err != nil {
		return datamodles.GenResultData(false, err.Error(), 0)
	}

	id, err := history.AddHistory(dbname, schema, tableName, historyProperty, toServerIds)
	if err != nil {
		return datamodles.GenResultData(false, err.Error(), 0)
	} else {
		return datamodles.GenResultData(true, "success", id)
	}
}

func stop(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	id := common.ParseIntDef(ctx.FormValue("id"), 0)
	if id == 0 {
		return datamodles.GenFailedMsg("id error not be int")
	}

	err := history.KillHistory(dbname, id)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	return datamodles.GenSuccessMsg("success")
}

func del(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	id := common.ParseIntDef(ctx.FormValue("id"), 0)
	if id == 0 {
		return datamodles.GenFailedMsg("id error not be int")
	}

	b := history.DelHistory(dbname, id)
	if b == false {
		return datamodles.GenFailedMsg("del error")
	}
	return datamodles.GenSuccessMsg("success")
}

func start(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	id := common.ParseIntDef(ctx.FormValue("id"), 0)
	if id == 0 {
		return datamodles.GenFailedMsg("id error not be int")
	}

	err := history.Start(dbname, id)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	return datamodles.GenSuccessMsg("success")
}

func checkWhere(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	tablename := ctx.FormValue("table_name")
	schema := ctx.FormValue("schema_name")
	property := ctx.FormValue("property")

	var err error
	var historyProperty history.HistoryProperty
	err = json.Unmarshal([]byte(property), &historyProperty)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	if historyProperty.Where == "" {
		return datamodles.GenSuccessMsg("success")
	}

	err = history.CheckWhere(dbname, schema, tablename, historyProperty.Where)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	return datamodles.GenSuccessMsg("success")
}
