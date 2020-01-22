package controllers

import (
    "fmt"
    "github.com/brokercap/Bifrost/manager/v2/datamodles"
    "github.com/brokercap/Bifrost/server/user"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "github.com/kataras/iris/v12/sessions"
)

type RootController struct {
    Ctx iris.Context
    Session *sessions.Session
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

var loginStaticView = mvc.View{
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

func (root *RootController) PostDologin() interface{}  {
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
    
    root.Session.Set("userId", userName)
    //sessionMgr.SetSessionVal(sessionID, "Group", groupName)
    return datamodles.GenSuccessMsg("success")
}

func (root *RootController) AnyLogout() {
    if root.isLoggedIn() {
        root.logout()
    }
    root.Ctx.Redirect("/login")
}
