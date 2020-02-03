package controllers

import (
	"fmt"
	"github.com/brokercap/Bifrost/manager/v2/common"
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
	"github.com/brokercap/Bifrost/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

func ChannelController(app *mvc.Application) {
	app.Router.Get("/list", hero.Handler(channelList))
	app.Router.Get("/tablelist", hero.Handler(channelTableList))
	app.Router.Post("/add", hero.Handler(channelAdd))
	app.Router.Post("/stop", hero.Handler(channelStop))
	app.Router.Post("/start", hero.Handler(channelStart))
	app.Router.Post("/del", hero.Handler(channelDel))
	app.Router.Post("/close", hero.Handler(channelClose))
}

func channelList(ctx iris.Context) hero.Result {
	dbname := ctx.FormValue("dbname")
	if ctx.FormValue("format") == "json" {
		return hero.Response{
			Object: server.GetDBObj(dbname).ListChannel(),
		}
	}
	return hero.View{
		Name: "channel.list.html",
		Data: iris.Map{
			"Title":       dbname + " - Channel List - Bifrost",
			"DbName":      dbname,
			"ChannelList": server.GetDBObj(dbname).ListChannel(),
		},
	}
}

func channelTableList(ctx iris.Context) hero.Result {
	dbname := ctx.FormValue("dbname")
	channelID := common.ParseIntDef(ctx.FormValue("channelid"), 0)
	channelInfo := server.GetChannel(dbname, channelID)
	if channelInfo == nil {
		return hero.Response{
			Object: datamodles.GenFailedMsg("channel not exist"),
		}
	}

	db := server.GetDBObj(dbname)
	tableMap := db.GetTableByChannelKey(dbname, channelID)
	return hero.View{
		Name: "channel.table.list.html",
		Data: iris.Map{
			"Title":       dbname + " - Table List - Channel - Bifrost",
			"DbName":      dbname,
			"ChannelName": channelInfo.Name,
			"ChannelID":   channelID,
			"TableList":   tableMap,
		},
	}
}

func channelAdd(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	chname := ctx.FormValue("channel_name")
	cosumerCountStr := ctx.FormValue("cosumerCount")
	cosumerCount, err := strconv.Atoi(cosumerCountStr)
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	db := server.GetDBObj(dbname)
	if db == nil {
		return datamodles.GenFailedMsg(dbname + " not exist")
	}
	_, channelId := db.AddChannel(chname, cosumerCount)
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccess("success", channelId)
}

func channelStop(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	channelID := common.ParseIntDef(ctx.FormValue("channelid"), 0)

	ch := server.GetChannel(dbname, channelID)
	if ch == nil {
		return datamodles.GenFailedMsg("channel not exist")
	}
	ch.Stop()
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}

func channelStart(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	channelID := common.ParseIntDef(ctx.FormValue("channelid"), 0)

	ch := server.GetChannel(dbname, channelID)
	if ch == nil {
		return datamodles.GenFailedMsg("channel not exist")
	}
	ch.Start()
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}

func channelDel(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	channelID := common.ParseIntDef(ctx.FormValue("channelid"), 0)

	db := server.GetDBObj(dbname)
	TableMap := db.GetTableByChannelKey(dbname, channelID)
	n := len(TableMap)
	if len(TableMap) > 0 {
		return datamodles.GenFailedMsg("The channel bind table count:" + fmt.Sprint(n))
	}
	r := server.DelChannel(dbname, channelID)
	if r == true {
		defer server.SaveDBConfigInfo()
		return datamodles.GenSuccessMsg("success")
	} else {
		return datamodles.GenFailedMsg("channel or db not exist")
	}
}

func channelClose(ctx iris.Context) interface{} {
	dbname := ctx.FormValue("dbname")
	channelID := common.ParseIntDef(ctx.FormValue("channelid"), 0)

	ch := server.GetChannel(dbname, channelID)
	if ch == nil {
		return datamodles.GenFailedMsg("channel not exist")
	}
	ch.Close()
	defer server.SaveDBConfigInfo()
	return datamodles.GenSuccessMsg("success")
}
