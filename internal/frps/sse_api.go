package frps

import (
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/xxl6097/glog/pkg/z"
	iface2 "github.com/xxl6097/go-frp-panel/pkg/comm/iface"
	"github.com/xxl6097/go-frp-panel/pkg/comm/ws"
)

func (this *frps) OnSseDisconnect(client *iface2.SSEClient) {
}

func (this *frps) OnSseNewConnection(client *iface2.SSEClient) {
	if strings.HasPrefix(client.SseId, ws.CLIENT_LIST) {
		//列表
		if this.webSocketApi != nil {
			list := this.webSocketApi.GetListSize()
			if list != nil && len(list) > 0 {
				this.sseApi.BroadcastByType(ws.CLIENT_LIST, iface2.SSEEvent{
					Event: ws.CLIENT_LIST,
					Payload: map[string]interface{}{
						"id":   client.SseId,
						"data": list,
					},
				})
			}
		}
	} else if strings.HasPrefix(client.SseId, ws.CLIENT_DETAIL) {
		//详情
		if this.sseApi != nil {
			this.clientRefresh(client)
			eve := iface2.SSEEvent{
				Event:   ws.SSE_CONNECT,
				Payload: client,
			}
			//z.Debugf("---->%+v", client)
			err := this.sseApi.Send(client, eve)
			if err != nil {
				z.Errorf("Send error: %s", err)
			} else {
				//z.Infof("Send success %v", client.SseId)
			}
		}

	}
}

func (this *frps) clientRefresh(client *iface2.SSEClient) {
	if this.webSocketApi != nil && client != nil {
		msg := iface2.Message[any]{
			Action: ws.CLIENT_REFRESH,
			SseID:  client.SseId,
			SecKey: client.SecKey,
			FrpId:  client.FrpID,
		}
		b, e := json.Marshal(msg)
		if e != nil {
			z.Errorf("getConfigs error: %v", e)
			return
		}
		e = this.webSocketApi.SendByKey(client.FrpID, client.SecKey, websocket.TextMessage, b)
		if e != nil {
			z.Errorf("getConfigs error: %v", e)
		}
	}
}
