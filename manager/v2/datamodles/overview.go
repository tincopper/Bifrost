package datamodles

import (
	"github.com/brokercap/Bifrost/config"
	"github.com/brokercap/Bifrost/plugin/driver"
	"runtime"
)

type OverView struct {
	DbCount              int
	ToServerCount        int
	PluginCount          int
	TableCount           int
	GoVersion            string
	BifrostVersion       string
	BifrostPluginVersion string
	StartTime            string
	GOOS                 string
	GOARCH               string
}

func NewOverView(dbCount int, tableCount int, toServerCount int,
	pluginCount int, startTime string) *OverView {
	return &OverView {
		DbCount:				dbCount,
		ToServerCount:			toServerCount,
		PluginCount:			pluginCount,
		TableCount:				tableCount,
		GoVersion:				runtime.Version(),
		BifrostVersion:			config.VERSION,
		BifrostPluginVersion:	driver.GetApiVersion(),
		StartTime:				startTime,
		GOOS:					runtime.GOOS,
		GOARCH:					runtime.GOARCH,
	}
}
