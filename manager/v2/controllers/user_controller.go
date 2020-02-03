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
    "github.com/brokercap/Bifrost/manager/v2/datamodles"
    "github.com/brokercap/Bifrost/server/user"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
)

type UserController struct {
    Ctx     iris.Context
}

func (c *UserController) PostUpdate() interface{} {
    userName := c.Ctx.FormValue("user_name")
    userPwd := c.Ctx.FormValue("password")
    userGroup := c.Ctx.FormValue("group")
    if userName == "" || userPwd == "" {
        return datamodles.GenFailedMsg("user_name and password must be not empty!")
    }
    err := user.UpdateUser(userName, userPwd, userGroup)
    if err != nil {
        return datamodles.GenFailedMsg("update failure, " + err.Error())
    }
    return datamodles.GenSuccessMsg("update success!")
}

func (c *UserController) PostDel() interface{} {
    userName := c.Ctx.FormValue("user_name")
    if userName == "" {
        return datamodles.GenFailedMsg("user_name must be not empty!")
    }
    err := user.DelUser(userName)
    if err != nil {
        return datamodles.GenFailedMsg("delete failure, " + err.Error())
    }
    return datamodles.GenSuccessMsg("delete success!")
}

func (c *UserController) GetList() {
    if c.Ctx.FormValue("format") == "json" {
        userList := getUserInfoList()
        _, _ = c.Ctx.JSON(userList)
    }
    // 重定向
    c.Ctx.Redirect("/user/list/page")
}

func (c *UserController) GetListPage() mvc.Result {
    return mvc.View{
        Name: "user.list.html",
        Data: iris.Map{"Title": "UserList-Bifrost",
            "UserList": getUserInfoList()},
    }
}

func getUserInfoList() []user.UserInfo {
    userList := user.GetUserList()
    //过滤密码,防止其他 monitor 用户查看到
    for k := range userList {
        userList[k].Password = ""
    }
    return userList
}
