package controllers

import (
	"fmt"
	"github.com/brokercap/Bifrost/config"
	"github.com/brokercap/Bifrost/manager/v2/datamodles"
    "github.com/brokercap/Bifrost/plugin/storage"
    "github.com/brokercap/Bifrost/server"
	"github.com/brokercap/Bifrost/server/user"
    "github.com/brokercap/Bifrost/xdb/driver"
    "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"runtime"
)

type RootController struct {
	Ctx     iris.Context
	Session *sessions.Session
	StartTime string
}

func (root *RootController) getCurrentUserID() string {
	userID := root.Session.GetStringDefault("userId", "")
	return userID
}

func (root *RootController) isLoggedIn() bool {
	return root.getCurrentUserID() != ""
}

func (root *RootController) logout() {
	root.Session.Destroy()
}

var loginStaticView = mvc.View {
	Name: "login.html",
	Data: iris.Map{"Title": "Login - Bifrost"},
}

// curl GET /login
func (root *RootController) GetLogin() mvc.Result {
	if root.isLoggedIn() {
		root.logout()
	}
	return loginStaticView
}

func (root *RootController) PostDologin() interface{} {
	userName := root.Ctx.FormValue("user_name")
	password := root.Ctx.FormValue("password")

	fmt.Println("userName:", password, "password:", password)
	if userName == "" {
		return datamodles.GenFailedMsg("user no exist")
	}

	userInfo := user.GetUserInfo(userName)
	if userInfo.Password != password {
		return datamodles.GenSuccessMsg("password error")
	}

	groupName := userInfo.Group
	if groupName == "" {
		groupName = "monitor"
	}

	root.Session.Set("UserName", userName)
	root.Session.Set("Group", groupName)
	return datamodles.GenSuccessMsg("success")
}

func (root *RootController) AnyLogout() {
	if root.isLoggedIn() {
		root.logout()
	}
	root.Ctx.Redirect("/login")
}

func (root *RootController) Get() mvc.Result {
	return mvc.View{
		Name: "index.html",
		Data: iris.Map{"Title": "Bifrost-Index"},
	}
}

func (root *RootController) GetOverview() interface{} {
	dbList := server.GetListDb()
	dbCount := len(dbList)
	tableCount := 0
    for _, v := range dbList {
        tableCount += v.TableCount
    }

    pluginCount := len(driver.Drivers())
    toServerCount := len(storage.GetToServerMap())

    return datamodles.NewOverView(dbCount, tableCount, toServerCount, pluginCount, root.StartTime)
}

func (root *RootController) GetServerinfo() interface{} {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	return &datamodles.ServerInfo{SeftMemStats: memStats}
}

func (root *RootController) GetGetversion() interface{} {
	return datamodles.GenSuccessMsg(config.VERSION)
}