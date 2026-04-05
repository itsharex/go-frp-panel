package frpc

import (
	httppkg "github.com/fatedier/frp/pkg/util/http"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-frp-panel/pkg"
	"github.com/xxl6097/go-frp-panel/pkg/utils"
	"github.com/xxl6097/go-service/pkg/gs"
	"go.uber.org/zap/zapcore"

	"fmt"
	"net/http"
	"path/filepath"
	"time"
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

func (this *frpc) adminHandlers(helper *httppkg.RouterRegisterHelper) {
	subRouter := helper.Router.NewRoute().Name("admin").Subrouter()
	subRouter.Use(helper.AuthMiddleware)
	staticPrefix := "/log/"
	baseDir := zutil.AppHome()
	subRouter.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, http.FileServer(http.Dir(baseDir))))

	subRouter.PathPrefix("/fserver/").Handler(http.StripPrefix("/fserver/", http.FileServer(http.Dir("/"))))
	subRouter.HandleFunc("/api/sse-stream", utils.SseHandler(logQueue))
	subRouter.HandleFunc("/api/files", this.upgrade.ApiFiles).Methods("PUT")

	subRouter.HandleFunc("/api/run", this.upgrade.ApiCMD).Methods("POST")
	subRouter.HandleFunc("/api/clear", this.upgrade.ApiClear).Methods("DELETE")
	// apis
	subRouter.HandleFunc("/api/version", this.upgrade.ApiVersion).Methods("GET")
	//subRouter.HandleFunc("/api/upgrade", this.upgrade.ApiUpdate).Methods("POST")
	//subRouter.HandleFunc("/api/upgrade", this.upgrade.ApiUpdate).Methods("PUT")
	//subRouter.HandleFunc("/api/checkversion", this.upgrade.ApiCheckVersion).Methods("GET")
	subRouter.HandleFunc("/api/restart", this.upgrade.ApiRestart).Methods("GET")
	subRouter.HandleFunc("/api/uninstall", this.upgrade.ApiUninstall).Methods("GET")

	subRouter.HandleFunc("/api/checkversion", gs.ApiCheckVersion(pkg.BinName)).Methods("GET")
	subRouter.HandleFunc("/api/upgrade", gs.ApiUpdate(this.install)).Methods("POST")
	subRouter.HandleFunc("/api/upgrade", gs.ApiUpdate(this.install)).Methods("PUT")

	subRouter.HandleFunc("/api/client/create", this.apiClientCreate).Methods("PUT")
	subRouter.HandleFunc("/api/client/create", this.apiClientCreate).Methods("POST")
	subRouter.HandleFunc("/api/client/upload", this.apiClientCreate).Methods("POST")
	subRouter.HandleFunc("/api/client/delete", this.apiClientDelete).Methods("DELETE")
	subRouter.HandleFunc("/api/client/status", this.apiClientStatus).Methods("GET")
	subRouter.HandleFunc("/api/client/list", this.apiClientList).Methods("GET")
	subRouter.HandleFunc("/api/client/config/get", this.apiClientConfigGet).Methods("GET")
	subRouter.HandleFunc("/api/client/config/set", this.apiClientConfigSet).Methods("POST")

	subRouter.HandleFunc("/api/proxy/ports", this.apiProxyPorts).Methods("GET")
	subRouter.HandleFunc("/api/proxy/ips", this.apiProxyLocalIps).Methods("GET")
	subRouter.HandleFunc("/api/proxy/port/check", this.apiProxyPortCheck).Methods("GET")
	subRouter.HandleFunc("/api/proxy/remote/ports", this.apiProxyRemotePorts).Methods("GET")
	subRouter.HandleFunc("/api/proxy/tcp/add", this.apiProxyTCPAdd).Methods("PUT")

	subRouter.HandleFunc("/api/proxy/github/api", this.apiProxyGithubApi).Methods("PUT")

	subRouter.HandleFunc("/api/client/config/import", this.apiClientConfigImport).Methods("POST")
	subRouter.HandleFunc("/api/client/config/export", this.apiClientConfigExport).Methods("POST")

}
