package main

import (
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/utils"
)

func main() {
	newVersion := "v0.4.98"
	oldVersion := "v0.4.99"
	hasNewVersion := utils.CompareVersions(newVersion, oldVersion)
	z.Debug("计算结果:", hasNewVersion)
}
