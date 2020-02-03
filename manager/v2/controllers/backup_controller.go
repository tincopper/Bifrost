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
	"github.com/brokercap/Bifrost/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/mvc"
	"io/ioutil"
	"time"
)

func BackupController(app *mvc.Application) {
	app.Router.Get("/export", hero.Handler(backupExport))
	app.Router.Post("/import", hero.Handler(backupImport))
}

func backupExport(ctx iris.Context) interface{} {
	b, err := server.GetSnapshotData()
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	fileName := "bifrost_" + time.Now().Format("2006-01-02 15:04:05") + ".json"
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
	return b
}

func backupImport(ctx iris.Context) interface{} {
	ctx.SetMaxRequestBodySize(32 << 20)
	file, _, err := ctx.FormFile("backup_file")
	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}
	fileContent, err := ioutil.ReadAll(file)
	file.Close()

	if err != nil {
		return datamodles.GenFailedMsg(err.Error())
	}

	server.DoRecoveryByBackupData(string(fileContent))
	return datamodles.GenSuccessMsg("success")
}
