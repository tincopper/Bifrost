package controllers

import (
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "github.com/kataras/iris/v12/sessions"
)

type TemplateHeader struct {
    Title string
}

func (TemplateHeader *TemplateHeader) setTile(title string) {
    TemplateHeader.Title = title
}

type UserController struct {
    Ctx iris.Context
    Session *sessions.Session
}

var loginStaticView = mvc.View{
    Name: "login.html",
    Data: iris.Map{"Title": "Login - Bifrost"},
}

func (user *UserController) GetLogin() mvc.Result {
    return loginStaticView
}

func UserLogin(ctx iris.Context) {
    data := TemplateHeader{Title: "Login - Bifrost"}
    ctx.View("login.html", data)
}

func UserDoLogin(ctx iris.Context) {

}

func UserLogout(ctx iris.Context) {

}

func UpdateUserController(ctx iris.Context) {

}

func DelUserController(ctx iris.Context) {

}

func ListUserController(ctx iris.Context) {

}