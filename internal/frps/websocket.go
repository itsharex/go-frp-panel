package frps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-frp-panel/pkg/comm"
	iface2 "github.com/xxl6097/go-frp-panel/pkg/comm/iface"
	"github.com/xxl6097/go-frp-panel/pkg/comm/ws"
)

func (this *frps) OnServerWebSocketMessageReceive(messageType int, payload []byte) {
	if payload != nil {
		//z.Debugf("来自frpc消息：%s", string(payload))
		var msg iface2.Message[any]
		err := json.Unmarshal(payload, &msg)
		if err != nil {
			z.Error(err)
			return
		}
		//z.Debugf("ws recv:%+v", msg)
		this.recvClientInfo(msg.SseID, msg.Action, msg.Data)
	}
}

func (this *frps) OnServerWebSocketDisconnect(session *iface2.WSSession) {
	if this.sseApi != nil {
		eve := iface2.SSEEvent{
			Event:   ws.DISCONNECT,
			Payload: session,
		}
		this.sseApi.Broadcast(eve)
	}
}

func (this *frps) OnServerWebSocketNewConnection(session *iface2.WSSession) {
}

func (this *frps) recvClientInfo(sseId, event string, data any) {
	if data == nil {
		z.Error("data is nil")
		//return
	}
	if this.sseApi != nil {
		eve := iface2.SSEEvent{
			Event:   event,
			Payload: data,
		}
		okk := this.sseApi.BroadcastTo(sseId, eve)
		if !okk {
			z.Errorf("Send error: %s sseId:%v", event, sseId)
		} else {
			z.Infof("Send success %s %v", event, sseId)
		}
	}
}

func (this *frps) apiClientCMD(w http.ResponseWriter, r *http.Request) {
	res, f := comm.Response(r)
	defer f(w)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		z.Error("body读取失败", err)
		res.Err(err)
		return
	}
	if body == nil {
		msg := "body is nil"
		z.Error(msg)
		res.Err(fmt.Errorf(msg))
		return
	}
	var msg iface2.Message[any]
	err = json.Unmarshal(body, &msg)
	if err != nil {
		z.Error("解析Json对象失败", err)
		res.Err(err)
		return
	}
	z.Debugf("发生给frpc:%s", string(body))
	if this.webSocketApi == nil {
		res.Error(fmt.Sprintf("webSocketApi is nil"))
		return
	}
	e := this.webSocketApi.SendByKey(msg.FrpId, msg.SecKey, websocket.TextMessage, body)
	if e != nil {
		z.Errorf("apiClientCMD error: %v", e)
		res.Err(e)
	} else {
		z.Infof("Send success %v", msg.FrpId)
		res.Ok("执行成功～")
	}
}
