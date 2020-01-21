package manager

import (
    "fmt"
    "github.com/brokercap/Bifrost/config"
    "github.com/brokercap/Bifrost/manager/v2/controllers"
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
)

// start manager web server
func StartV2(endpoint string) {
    app := irisInit()
    // 加载路由
    loadRout(app)
    // 启动
    startApp(app, endpoint)
}

func irisInit() *iris.Application {
    app := iris.New()
    app.Logger().SetLevel("debug")
    app.RegisterView(iris.HTML("./manager/template", ".html"))
    app.HandleDir("/", "./manager/public")
    app.HandleDir("/plugin", "./")
    // 错误处理
    app.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
        //_ = ctx.View("404.html")
        _, _ = ctx.HTML("<h1>404</h1>")
    })
    app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
        _ = ctx.View("index.html")
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

func loadRout(app *iris.Application) {
    //app.Get("/login", controllers.UserLogin)
    app.Get("/logout", controllers.UserLogout)
    app.Post("/dologin", controllers.UserDoLogin)
    app.UseGlobal(func(c context.Context) {
        fmt.Println("-----------")
        c.Next()
    })
    mvc.New(app.Party("/")).Register(sessionMgr).Handle(new(controllers.UserController))
    mvc.New(app.Party("/login")).Register(sessionMgr).Handle(new(controllers.UserController))
    mvc.New(app.Party("/login")).Register(sessionMgr).Handle(new(controllers.UserController))
    
    mvc.Configure(app.Party("/"), )
    mvc.Configure(app.Party("/user"), userRoute)
}

func userRoute(app *mvc.Application) {
    app.Router.Put("/update", controllers.UpdateUserController)
    app.Router.Delete("/del", controllers.DelUserController)
    app.Router.Get("/list", controllers.ListUserController)
}
