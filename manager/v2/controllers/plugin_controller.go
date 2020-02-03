package controllers

import (
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
	"github.com/brokercap/Bifrost/plugin"
	"github.com/brokercap/Bifrost/plugin/driver"
	"github.com/kataras/iris"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
)

func PluginController(app *mvc.Application) {
	app.Router.Get("/list", hero.Handler(pluginList))
	app.Router.Post("/reload", hero.Handler(pluginReload))
}

func pluginList() hero.Result {

	drivers := driver.Drivers()
	//因为plugin 加载so 插件，有可能会异常，所以这里需要你把异常的插件列表也加载进来并进行显示出来
	errorPluginMap := plugin.GetErrorPluginList()
	for name, v := range errorPluginMap{
		drivers[name] = v
	}

	return hero.View{
		Name: "plugin.list.html",
		Data: iris.Map{
			"Title": "Plugin List - Bifrost",
			"PluginAPIVersion": driver.GetApiVersion(),
			"Drivers": drivers,
		},
	}
}

func pluginReload() interface{} {
	err := plugin.LoadPlugin()
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	return datamodles.GenSuccessMsg("success")
}
