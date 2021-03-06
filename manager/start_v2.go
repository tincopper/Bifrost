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
package manager

import (
    "github.com/brokercap/Bifrost/config"
    "github.com/brokercap/Bifrost/manager/v2/controllers"
    "github.com/brokercap/Bifrost/manager/v2/datamodles"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/mvc"
    "github.com/kataras/iris/v12/sessions"
    "time"
)

var (
    //app = iris.New()
    sessManager = sessions.New(sessions.Config{
        Cookie:  "xgo_cookie",
        Expires: 24 * time.Hour,
    })
    startTime = time.Now().Format("2006-01-03 15:04:05")
)

// start manager web server
func StartV2(endpoint string) {
    app := irisInit()
    // 加载路由
    loadRoute(app)
    // 启动
    startApp(app, endpoint)
}

func irisInit() *iris.Application {
    app := iris.New()
    app.Logger().SetLevel("debug")
    app.RegisterView(iris.HTML("./manager/template", ".html"))
    app.HandleDir("/", "./manager/public")
    app.HandleDir("/plugin", "./")
    app.OnAnyErrorCode(func(ctx iris.Context) {
        ctx.ViewData("Message", ctx.Values().
            GetStringDefault("message", "The page you're looking for doesn't exist"))
        ctx.View("shared/error.html")
    })
    // 错误处理
    app.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
        //_ = ctx.View("404.html")
        _, _ = ctx.HTML(ctx.Method() + " " + ctx.FullRequestURI() + "<h1>404</h1>")
    })
    app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
        _ = ctx.View("index.html")
    })
    // 权限处理
    app.Use(func(c context.Context) {
        if checkPermission(c) {
            c.Next()
        } else {
            c.StopExecution()
        }
    })
    return app
}

func startApp(app *iris.Application, endpoint string) {
    var err error
    if config.TLS {
        // 使用文件 安全传输层协议`TLS`
        err = app.Run(iris.TLS(endpoint, config.TLSServerKeyFile, config.TLSServerCrtFile), iris.WithOptimizations)
    } else {
        err = app.Run(iris.Addr(endpoint), iris.WithOptimizations)
    }
    if err != nil {
        app.Logger().Println("Manager Start Err:", err)
    }
}

func loadRoute(app *iris.Application) {
    mvc.New(app.Party("/")).Register(sessManager.Start, startTime).Handle(new(controllers.RootController))
    mvc.Configure(app.Party("/user")).Handle(new(controllers.UserController))
    mvc.Configure(app.Party("/flow")).Handle(new(controllers.FlowController))
    mvc.Configure(app.Party("/db"), controllers.DbController)
    mvc.Configure(app.Party("/backup"), controllers.BackupController)
    mvc.Configure(app.Party("/channel"), controllers.ChannelController)
    mvc.Configure(app.Party("/history"), controllers.HistoryController)
    mvc.Configure(app.Party("/plugin"), controllers.PluginController)
    mvc.Configure(app.Party("/table"), controllers.TableController)
}

func checkPermission(c context.Context) bool {
    if c.GetHeader("Authorization") != "" {
        username, password, ok := c.Request().BasicAuth()
        if !ok || username == "" {
            c.JSON(datamodles.GenFailedMsg("Author error"))
            return false
        }
        targetPwd := config.GetConfigVal("user", username)
        if targetPwd != password {
            c.JSON(datamodles.GenFailedMsg("password error"))
            return false
        }
        groupName := config.GetConfigVal("groups", username)
        if !checkAuthority(groupName, c) {
            return false
        }
    } else {
        username := sessManager.Start(c).GetString("UserName")
        if username != "" {
            groupName := sessManager.Start(c).GetString("Group")
            if !checkAuthority(groupName, c) {
                return false
            }
        } else {
            requestURI := c.Request().RequestURI
            if requestURI != "/login" && requestURI != "/dologin" && requestURI != "/logout" {
                c.Redirect("/login")
                return false
            }
        }
    }
    return true
}

func checkAuthority(groupName string, c context.Context) bool {
    if groupName != "administrator" && checkWriteRequest(c.FullRequestURI()) {
        c.JSON(datamodles.GenFailedMsg("user group : [ " + groupName + " ] no authority"))
        return false
    }
    return true
}
