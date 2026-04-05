package frps

import (
	"os"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-frp-panel/pkg"
	"github.com/xxl6097/go-frp-panel/pkg/model"
	"github.com/xxl6097/go-frp-panel/pkg/utils"
	"github.com/xxl6097/go-service/pkg/ukey"
)

var cfgData *CfgModel
var cfgBytes []byte

type CfgModel struct {
	Frps v1.ServerConfig `json:"frps"`
	Data any             `json:"data"`
}

func (this *CfgModel) Bytes() []byte {
	return utils.ObjectToTomlText(this)
}

//func Test(c *CfgModel) {
//	cfgData = c
//}

func load() error {
	//defer glog.Flush()
	byteArray, err := ukey.Load()
	if err != nil {
		//z.Error(err)
		return err
	}
	cfgBytes = byteArray
	c := CfgModel{}
	err = utils.TomlTextToObject(byteArray, &c)
	if err != nil {
		z.Println("cfgBytes解析错误", err)
		return err
	}
	cfgData = &c
	//z.Printf("%d 配置加载成功：%+v\n", os.Getpid(), cfgData)
	pkg.Version()
	return nil
}

func GetCfgModel() *CfgModel {
	if cfgData == nil {
		err := load()
		if err != nil {
			return nil
		}
	}
	return cfgData
}

func GetCfgBuffer() []byte {
	cfg := GetCfgModel()
	if cfg == nil {
		return nil
	}
	return utils.ObjectToTomlText(cfg)
}
func SetCfgModel(temp *CfgModel) {
	cfgData = temp
}

//func PrintCfg() {
//	if cfgBytes != nil {
//		z.Println(string(cfgBytes))
//	}
//}

//func IsInit() error {
//	//z.Println("IsInit")
//	defer glog.Flush()
//	err := load()
//	if err != nil {
//		//z.Println(err)
//		return err
//	}
//	return nil
//}

//func Assert() {
//	//z.Println("Assert")
//	if IsInit() != nil {
//		if utils.IsMacOs() {
//			return
//		}
//		os.Exit(0)
//	}
//}

func (this *frps) InitClientsConfig() {
	host := os.Getenv("CLOUD_HOST")
	user := os.Getenv("CLOUD_USER")
	pass := os.Getenv("CLOUD_PASS")
	if host == "" || user == "" || pass == "" {
		return
	}
	this.cloudApi = &model.CloudApi{
		Addr: host,
		User: user,
		Pass: pass,
	}
	err := utils.Import(*this.cloudApi)
	if err != nil {
		z.Errorf("更新失败 %+v", err)
	}
}
