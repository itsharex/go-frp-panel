package frps

import (
	"os"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-frp-panel/pkg"
	"github.com/xxl6097/go-service/pkg/github"
)

func (this *frps) checkFrpc() {
	z.L().Info("checkFrpc请求")
	github.Api().Request(pkg.GithubUser, pkg.GithubRepo)
}

func (this *frps) check() {
	z.L().Info("开始检测客户端...")
	for {
		this.checkFrpc()
		time.Sleep(time.Hour * 8)
	}
}

func (this *frps) CheckVersion() {
	checks := os.Getenv("CHECK_CLIENTS")
	if checks != "" {
		return
	}
	go this.check()
}
