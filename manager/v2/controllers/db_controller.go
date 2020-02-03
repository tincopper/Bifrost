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
	toserver "github.com/brokercap/Bifrost/plugin/storage"
	"github.com/brokercap/Bifrost/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func DbController(app *mvc.Application) {
	app.Router.Post("/add", hero.Handler(addDb))
	app.Router.Post("/update", hero.Handler(updateDb))
	app.Router.Post("/stop", hero.Handler(stopDb))
	app.Router.Post("/start", hero.Handler(startDb))
	app.Router.Post("/close", hero.Handler(closeDb))
	app.Router.Post("/del", hero.Handler(delDb))
	app.Router.Get("/list", hero.Handler(listDb))
	app.Router.Post("/check_uri", hero.Handler(checkDbConn))
	app.Router.Post("/checkposition", hero.Handler(checkDbPosition))
	app.Router.Get("/detail", hero.Handler(dbDetail))
	app.Router.Get("/tablelist", hero.Handler(dbTableList))
	app.Router.Get("/tablefields", hero.Handler(dbTableFields))
}

func dbTableFields(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	schemaName := ctx.FormValue("schemaName")
	tableName := ctx.FormValue("TableName")

	schemaName = common.TansferSchemaName(schemaName)
	tableName = common.TansferTableName(tableName)

	DBObj := server.GetDBObj(dbname)
	dbUri := DBObj.ConnectUri
	dbConn := common.DBConnect(dbUri)
	if dbConn == nil {
		return nil
	}
	defer dbConn.Close()
	tableFieldsList := common.GetSchemaTableFieldList(dbConn, schemaName, tableName)
	return tableFieldsList
}

func dbTableList(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	schemaName := ctx.FormValue("schemaName")
	dBObj := server.GetDBObj(dbname)
	dbUri := dBObj.ConnectUri
	dbConn := common.DBConnect(dbUri)
	if dbConn == nil {
		return nil
	}
	defer dbConn.Close()

	type ResultType struct {
		TableName   string
		ChannelName string
		AddStatus   bool
		TableType   string
	}
	var data []ResultType
	data = make([]ResultType, 0)
	TableList := common.GetSchemaTableList(dbConn, schemaName)
	TableList = append(TableList, common.TableListStruct{TableName: "AllTables", TableType: ""})
	var schemaName0, tableName0 string
	schemaName0 = common.TansferSchemaName(schemaName)

	for _, tableInfo := range TableList {
		tableName := tableInfo.TableName
		tableType := tableInfo.TableType
		tableName0 = common.TansferTableName(tableName)
		t := dBObj.GetTable(schemaName0, tableName0)
		if t == nil {
			data = append(data, ResultType{TableName: tableName, ChannelName: "", AddStatus: false, TableType: tableType})
		} else {
			t2 := dBObj.GetChannel(t.ChannelKey)
			if t2 == nil {
				data = append(data, ResultType{TableName: tableName, ChannelName: "", AddStatus: false, TableType: tableType})
			} else {
				data = append(data, ResultType{TableName: tableName, ChannelName: t2.Name, AddStatus: true, TableType: tableType})
			}
		}
	}
	return data
}

func dbDetail(ctx iris.Context) hero.Result {
	dbname := ctx.FormValue("dbname")
	dbUri := server.GetDBObj(dbname).ConnectUri
	dbConn := common.DBConnect(dbUri)
	if dbConn == nil {
		return nil
	}
	defer dbConn.Close()
	dataBaseList := common.GetSchemaList(dbConn)
	dataBaseList = append(dataBaseList, "AllDataBases")

	return hero.View{
		Name: "db.detail.html",
		Data: iris.Map{
			"Title":        dbname + " - Detail - Bifrost",
			"DbName":       dbname,
			"DataBaseList": dataBaseList,
			"ToServerList": toserver.GetToServerMap(),
			"ChannelList":  server.GetDBObj(dbname).ListChannel(),
		},
	}
}

func checkDbPosition(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	type dbInfoStruct struct {
		BinlogFile            string
		BinlogPosition        int
		BinlogTimestamp       uint32
		CurrentBinlogFile     string
		CurrentBinlogPosition int
		NowTimestamp          uint32
		DelayedTime           uint32
	}
	dbObj := server.GetDbInfo(dbname)
	if dbObj == nil {
		return datamodles.GenFailedMsg(dbname + " not esxit")
	}

	dbUri := dbObj.ConnectUri
	dbInfo := &dbInfoStruct{}
	dbInfo.BinlogFile = dbObj.BinlogDumpFileName
	dbInfo.BinlogPosition = int(dbObj.BinlogDumpPosition)
	dbInfo.BinlogTimestamp = uint32(dbObj.BinlogDumpTimestamp)
	err := func(dbUri string) (e error) {
		e = nil
		defer func() {
			if err := recover(); err != nil {
				log.Println(string(debug.Stack()))
				e = fmt.Errorf(fmt.Sprint(err))
				return
			}
		}()
		dbconn := common.DBConnect(dbUri)
		if dbconn != nil {
			e = nil
		} else {
			e = fmt.Errorf("db conn ,uknow error")
		}
		defer dbconn.Close()
		MasterBinlogInfo := common.GetBinLogInfo(dbconn)
		if MasterBinlogInfo.File != "" {
			dbInfo.CurrentBinlogFile = MasterBinlogInfo.File
			dbInfo.CurrentBinlogPosition = MasterBinlogInfo.Position
		} else {
			e = fmt.Errorf("The binlog maybe not open,or no replication client privilege(s).you can show log more.")
		}
		return
	}(dbUri)
	dbInfo.NowTimestamp = uint32(time.Now().Unix())
	if dbInfo.BinlogTimestamp > 0 && (dbInfo.CurrentBinlogFile != dbInfo.BinlogFile || dbInfo.BinlogPosition != dbInfo.CurrentBinlogPosition) {
		dbInfo.DelayedTime = dbInfo.NowTimestamp - dbInfo.BinlogTimestamp
	}

	if err != nil {
		return datamodles.GenResultData(false, err.Error(), *dbInfo)
	} else {
		return datamodles.GenResultData(true, "success", *dbInfo)
	}
}

func checkDbConn(ctx iris.Context) interface{} {
	dbUri := ctx.FormValue("uri")
	type dbInfoStruct struct {
		BinlogFile     string
		BinlogPosition int
		ServerId       int
		BinlogFormat   string
	}
	dbInfo := &dbInfoStruct{}
	err := func(dbUri string) (e error) {
		e = nil
		defer func() {
			if err := recover(); err != nil {
				log.Println(string(debug.Stack()))
				e = fmt.Errorf(fmt.Sprint(err))
				return
			}
		}()
		dbconn := common.DBConnect(dbUri)
		if dbconn != nil {
			e = nil
		} else {
			e = fmt.Errorf("db conn ,uknow error")
		}
		defer dbconn.Close()
		masterBinlogInfo := common.GetBinLogInfo(dbconn)
		if masterBinlogInfo.File != "" {
			dbInfo.BinlogFile = masterBinlogInfo.File
			dbInfo.BinlogPosition = masterBinlogInfo.Position
			dbInfo.ServerId = common.GetServerId(dbconn)
			variablesMap := common.GetVariables(dbconn, "binlog_format")
			if _, ok := variablesMap["binlog_format"]; ok {
				dbInfo.BinlogFormat = variablesMap["binlog_format"]
			}
		} else {
			e = fmt.Errorf("The binlog maybe not open,or no replication client privilege(s).you can show log more.")
		}
		return
	}(dbUri)

	if err != nil {
		return datamodles.GenResultData(false, err.Error(), *dbInfo)
	} else {
		return datamodles.GenResultData(true, "success", *dbInfo)
	}
}

func listDb(ctx iris.Context) hero.Result {
	if ctx.FormValue("format") == "json" {
		return hero.Response{Object: server.GetListDb()}
	}
	return hero.View{
		Name: "db.list.html",
		Data: iris.Map{
			"Title":  "Bifrost",
			"DBList": server.GetListDb(),
		},
	}
}

func delDb(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	result := server.DelDB(dbname)
	defer server.SaveDBConfigInfo()
	if !result {
		return datamodles.GenFailedMsg("failed")
	}
	return datamodles.GenSuccessMsg("success")
}

func closeDb(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	result := server.DbList[dbname].Close()
	if !result {
		return datamodles.GenFailedMsg("failed")
	}
	return datamodles.GenSuccessMsg("success")
}

func startDb(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	result := server.DbList[dbname].Start()
	if result != true {
		return datamodles.GenFailedMsg("failed")
	} else {
		defer server.SaveDBConfigInfo()
		return datamodles.GenSuccessMsg("success")
	}
}

func stopDb(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	server.DbList[dbname].Stop()
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}

func updateDb(ctx iris.Context) interface{} {
	dbname := strings.Trim(ctx.FormValue("dbname"), "")
	connUri := strings.Trim(ctx.FormValue("uri"), "")
	filename := strings.Trim(ctx.FormValue("filename"), "")
	maxFilename := strings.Trim(ctx.FormValue("max_filename"), "")
	serverIdStr := strings.Trim(ctx.FormValue("serverid"), "")
	positionStr := strings.Trim(ctx.FormValue("position"), "")
	maxPositionStr := strings.Trim(ctx.FormValue("max_position"), "")
	updateToServerStr := strings.Trim(ctx.FormValue("update_toserver"), "")

	var errMsg = ""
	position, err := strconv.Atoi(positionStr)
	if err != nil {
		errMsg = "position is err"
	}
	serverId, err := strconv.Atoi(serverIdStr)
	if err != nil {
		errMsg += "|serverid is err"
	}
	maxPosition, err := strconv.Atoi(maxPositionStr)
	if err != nil {
		//errMsg += "|maxPosition is err"
		maxPosition = 0
	}
	updateToServer, err := strconv.Atoi(updateToServerStr)
	if err != nil {
		updateToServer = 0
	}

	if errMsg != "" {
		return datamodles.GenFailedMsg(errMsg)
	}

	defer server.SaveDBConfigInfo()
	err = server.UpdateDB(dbname, connUri, filename, uint32(position), uint32(serverId),
		maxFilename, uint32(maxPosition), time.Now().Unix(), int8(updateToServer))
	if err == nil {
		return datamodles.GenSuccessMsg("success")
	} else {
		return datamodles.GenFailedMsg(err.Error())
	}
}

func addDb(ctx iris.Context) interface{} {
	dbname := strings.Trim(ctx.FormValue("dbname"), "")
	connUri := strings.Trim(ctx.FormValue("uri"), "")
	filename := strings.Trim(ctx.FormValue("filename"), "")
	maxFilename := strings.Trim(ctx.FormValue("max_filename"), "")
	serverIdStr := strings.Trim(ctx.FormValue("serverid"), "")
	positionStr := strings.Trim(ctx.FormValue("position"), "")
	maxPositionStr := strings.Trim(ctx.FormValue("max_position"), "")

	var errMsg = ""
	serverId, err := strconv.Atoi(serverIdStr)
	if err != nil {
		errMsg += "|serverId is err"
	}
	position, err := strconv.Atoi(positionStr)
	if err != nil {
		errMsg += "|position is err"
	}
	maxPosition, err := strconv.Atoi(maxPositionStr)
	if err != nil {
		//errMsg += "|maxPosition is err"
		maxPosition = 0
	}

	if errMsg != "" {
		return datamodles.GenFailedMsg(errMsg)
	}
	defer server.SaveDBConfigInfo()
	server.AddNewDB(dbname, connUri, filename, uint32(position), uint32(serverId),
		maxFilename, uint32(maxPosition), time.Now().Unix())
	c, _ := server.GetDBObj(dbname).AddChannel("default", 1)
	if c != nil {
		c.Start()
	}
	return datamodles.GenSuccessMsg("success")
}
