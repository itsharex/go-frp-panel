package utils

import (
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/pkg/utils/util"
)

func ShowUpDirSize() {
	updir := zutil.AppHome("temp", "upgrade")
	total, used, free, err := util.GetDiskUsage(updir)
	z.Printf("Current Working Directory: %s %v\n", updir, err)
	z.Printf("Total space: %d bytes %v\n", total, ByteCountIEC(total))
	z.Printf("Used space: %d bytes %v\n\n", used, ByteCountIEC(used))
	z.Printf("Free space: %d bytes %v\n\n", free, ByteCountIEC(free))
}

func GetAppSpace() (string, string, string) {
	dir := zutil.AppHome()
	total, used, free, _ := util.GetDiskUsage(dir)
	//z.Printf("Current Working Directory: %s %v\n", dir, err)
	//z.Printf("Total space: %d bytes %v\n", total, ByteCountIEC(total))
	//z.Printf("Used space: %d bytes %v\n\n", used, ByteCountIEC(used))
	//z.Printf("Free space: %d bytes %v\n\n", free, ByteCountIEC(free))
	return ByteCountIEC(total), ByteCountIEC(used), ByteCountIEC(free)
}

func HasDiskSpace() bool {
	size := GetSelfSize()
	size *= 16
	dir := zutil.AppHome()
	//total, used, free, err := util.GetDiskUsage(dir)
	_, _, free, _ := util.GetDiskUsage(dir)
	//z.Printf("Current Working Directory: %s %v\n", dir, err)
	//z.Printf("Total space: %d bytes %v\n", total, ByteCountIEC(total))
	//z.Printf("Used space: %d bytes %v\n\n", used, ByteCountIEC(used))
	//z.Printf("Free space: %d bytes %v\n\n", free, ByteCountIEC(free))
	if free > size {
		return true
	}
	return false
}
