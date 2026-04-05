package frps

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/gorilla/mux"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-frp-panel/internal/com/model"
	"github.com/xxl6097/go-frp-panel/pkg"
	"github.com/xxl6097/go-frp-panel/pkg/comm"
	"github.com/xxl6097/go-frp-panel/pkg/utils"
	"github.com/xxl6097/go-service/pkg/ukey"
	utils2 "github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logQueue = utils.NewLogQueue()

func init() {
	//glog.Hook(func(bytes []byte) {
	//	logQueue.AddMessage(string(bytes[2:]))
	//})
	z.GetLogConfig().Hook = func(entry zapcore.Entry) error {
		time := entry.Time.Format(time.DateTime)
		msg := entry.Message
		lineNum := entry.Caller.Line
		filepath.Base(entry.Caller.File)
		logs := fmt.Sprintf("%s %s:%d %s", time, filepath.Base(entry.Caller.File), lineNum, msg)
		logQueue.AddMessage(logs)
		return nil
	}
}

// /api/shutdown
func (this *frps) apiShutdown(w http.ResponseWriter, r *http.Request) {
	res := comm.GeneralResponse{Code: 0}
	defer func() {
		log.Infof("Http response [%s]: res: %+v", r.URL.Path, res)
		w.WriteHeader(res.Code)
		if len(res.Msg) > 0 {
			_, _ = w.Write([]byte(res.Msg))
		}
	}()

	log.Infof("Http request: [%s]", r.URL.Path)
	res.Msg = "ok"
}

func (this *frps) apiServerConfigSet(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	// 读取请求体
	tomlBytes, err := io.ReadAll(r.Body)
	if err != nil {
		res.Error(fmt.Sprintf("读取body失败%v", err))
		return
	}
	if utils2.IsPathExist(this.cfgFilePath) {
		err = utils.Write(this.cfgFilePath, tomlBytes)
		if err != nil {
			res.Error(err.Error())
		} else {
			res.Ok("配置更新成功～")
		}
		return
	}
	z.L().Debug("原始数据", zap.ByteString("tomlBytes", tomlBytes))
	newFrpsCfg := v1.ServerConfig{}
	err = utils.TomlTextToObject(tomlBytes, &newFrpsCfg)
	if err != nil {
		res.Error(fmt.Sprintf("配置失败：%v", err))
		return
	}
	z.L().Debug("对象数据", zap.Any("config", newFrpsCfg))
	cfg := GetCfgModel()
	cfg.Frps = newFrpsCfg
	z.L().Debug("cfg", zap.Any("cfg", newFrpsCfg))
	//下载和接收的最新文件 名称为上传文件的原始名称
	newBufferBytes, err := ukey.GenConfig(GetCfgBuffer(), false)
	if err != nil {
		res.Error(fmt.Sprintf("gen config err: %v", err))
		z.Error(res.Msg)
		return
	}
	if this.install != nil {
		err = this.install.UpgradeByBuffer(newBufferBytes)
		//err = this.install.Upgrade(r.Context(), signFilePath, "override")
		if err != nil {
			res.Error(fmt.Sprintf("更新失败～%v", err))
			return
		}
		res.Ok("配置更新成功～")
	}
}

//func (this *frps) apiServerConfigSet(w http.ResponseWriter, r *http.Request) {
//	res, f := comm.Response(r)
//	defer f(w)
//	// 读取请求体
//	tomlBytes, err := io.ReadAll(r.Body)
//	if err != nil {
//		res.Error(fmt.Sprintf("读取body失败%v", err))
//		return
//	}
//	if utils2.IsPathExist(this.cfgFilePath) {
//		err = utils.Write(this.cfgFilePath, tomlBytes)
//		if err != nil {
//			res.Error(err.Error())
//		} else {
//			res.Ok("配置更新成功～")
//		}
//		return
//	}
//	//z.Println(tomlBytes)
//	frpsCfg := v1.ServerConfig{}
//	err = utils.TomlTextToObject(tomlBytes, &frpsCfg)
//	if err != nil {
//		res.Error(fmt.Sprintf("配置失败：%v", err))
//		return
//	}
//	cfg := GetCfgModel()
//	cfg.Frps = frpsCfg
//	filePath, err := os.Executable()
//	if err != nil {
//		res.Error(fmt.Sprintf("%v", err))
//		return
//	}
//	//下载和接收的最新文件 名称为上传文件的原始名称
//	newBufferBytes, err := ukey.GenConfig(GetCfgBuffer(), false)
//	if err != nil {
//		res.Error(fmt.Sprintf("gen config err: %v", err))
//		z.Error(res.Msg)
//		return
//	}
//	signFilePath, err := utils.SignAndInstall(newBufferBytes, ukey.GetBuffer(), filePath)
//	if err != nil {
//		res.Error(err.Error())
//	} else {
//		//defer utils.Delete(signFilePath, "签名文件")
//		if this.install != nil {
//			err = this.install.Upgrade(r.Context(), signFilePath)
//			//err = this.install.Upgrade(r.Context(), signFilePath, "override")
//			if err != nil {
//				res.Error(fmt.Sprintf("更新失败～%v", err))
//				return
//			}
//			res.Ok("配置更新成功～")
//		}
//
//	}
//}

// /api/restart
func (this *frps) apiRestart(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	res.Msg = "restart sucess"
	if res.Code == 0 && this.install != nil {
		go func() {
			time.Sleep(time.Second)
			err := this.install.Restart()
			if err != nil {
				z.Error("重启失败")
			}
			z.Error("重启ok")
		}()
	}
}

func (this *frps) apiPanelinfo(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	res.Sucess("获取成功", map[string]interface{}{
		"appName":     pkg.AppName,
		"gitRevision": pkg.GitRevision,
		"gitBranch":   pkg.GitBranch,
		"goVersion":   pkg.GoVersion,
		"displayName": pkg.DisplayName,
		"description": pkg.Description,
		"appVersion":  pkg.AppVersion,
		"buildTime":   pkg.BuildTime,
	})
}

// /api/server/config/get
func (this *frps) apiServerConfigGet(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	if utils2.IsPathExist(this.cfgFilePath) {
		content, err := utils.Read(this.cfgFilePath)
		if err != nil {
			res.Error(err.Error())
		} else {
			res.Raw = content
		}
		return
	}
	frpsToml := GetCfgModel().Frps
	z.L().Info("获取Frps配置", zap.Any("config", frpsToml))
	res.Raw = utils.ObjectToTomlText(frpsToml)
}

// /api/proxy/:type
func (svr *frps) apiProxyByType(w http.ResponseWriter, r *http.Request) {
	res := comm.GeneralResponse{Code: 200}
	params := mux.Vars(r)
	proxyType := params["type"]

	defer func() {
		log.Infof("Http response [%s]: code [%d]", r.URL.Path, res.Code)
		w.WriteHeader(res.Code)
		if len(res.Msg) > 0 {
			_, _ = w.Write([]byte(res.Msg))
		}
	}()
	log.Infof("Http request: [%s]", r.URL.Path)

	res.Msg = proxyType
}

func (this *frps) apiBindInfo(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	port := this.cfg.BindPort
	bindPort := os.Getenv("BIND_PORT")
	if bindPort != "" {
		n, e := strconv.Atoi(bindPort)
		if e != nil {
			res.Error(e.Error())
			z.Errorf("bind port [%s] err %v", bindPort, e)
			return
		}
		//bindPort = strconv.Itoa(n)
		port = n
		z.Debugf("bind port [%s] %v", bindPort, e)
	}
	data := map[string]interface{}{
		"bindPort": port,
	}
	res.Any(data)
}

func (this *frps) apiEnv(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	name := r.URL.Query().Get("name")
	res.Raw = []byte(fmt.Sprintf("%s：%s", name, os.Getenv(name)))
}

func (this *frps) GetUserAll() ([]model.User, error) {
	binpath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(filepath.Dir(binpath), "user", "*.json"))
	if err != nil {
		return nil, err
	}

	var obj map[string]int
	if this.webSocketApi != nil {
		obj = this.webSocketApi.GetListSize()
	}
	var users []model.User
	for _, file := range files {
		user, e := model.Read(file)
		if e == nil {
			if obj != nil && len(obj) > 0 {
				n, ok := obj[user.ID]
				if ok {
					user.Count = n
				}
			}
			users = append(users, *user)
		}
	}
	return users, nil
}

func (this *frps) GetUser(id string) (*model.User, error) {
	binpath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	uPath := filepath.Join(filepath.Dir(binpath), "user", fmt.Sprintf("%s.json", id))

	if !utils2.FileExists(uPath) {
		return nil, fmt.Errorf("%s not exists", uPath)
	}
	return model.Read(uPath)
}
