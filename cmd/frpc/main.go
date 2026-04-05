package main

import (
	"github.com/xxl6097/go-frp-panel/cmd/frpc/service"
	"github.com/xxl6097/go-frp-panel/pkg"
	"github.com/xxl6097/go-frp-panel/pkg/utils"
)

func init() {
	if utils.IsMacOs() {
		pkg.AppName = "acfrpc"
		pkg.DisplayName = "acfrpc"
		pkg.Description = "acfrpc"
	}
}

//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest
func main() {
	service.Bootstrap()
}
