package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ohayao/common/http"
	"github.com/ohayao/mvc"
)

type ttn struct{}

type requestBody struct {
	ApiKey        string      `json:"apiKey"`
	ApplicationID string      `json:"applicationId"`
	DeviceID      string      `json:"deviceId"`
	Cmd           string      `json:"cmd,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

func cvt(body []byte) *requestBody {
	res := &requestBody{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil
	}
	return res
}

/*
igl-api
NNSXS.5IF2ERUGMBUQGWCTTLW3LPRVSEHVAZHCCOHIIUY.4KLC2NV6UCBD22LNYQ6BCKRSF5MILA4LLHM6BKDLQFRU32Z5UMTQ
body := `{
	"downlinks": [{
		"frm_payload": "MTIzNDU2Nzg5",
		"f_port": 2
	}]
}`
*/
const (
	baseURL = "http://120.79.166.115:1885/api/v3"
)

func (*ttn) Post_SendMessage(ctx *mvc.HTTPContext) {
	body := cvt(ctx.GetBody())
	if body == nil {
		ctx.JSON(500, `{}`)
		return
	}

	hp := http.New(fmt.Sprintf("%s/as/applications/%s/devices/%s/down/push", baseURL, body.ApplicationID, body.DeviceID))
	hp.Header("Authorization", fmt.Sprintf("Bearer %s", body.ApiKey))
	res := hp.PostBytes(bytes.NewBufferString(fmt.Sprintf(`{
		"downlinks": [{
			"frm_payload": "%s",
			"f_port": 2
		}]
	}`, body.Data)))
	ctx.JSON(200, res)
}

func (*ttn) Post_message__zanshibuyong(ctx *mvc.HTTPContext) {
	body := cvt(ctx.GetBody())
	if body == nil {
		ctx.JSON(500, `{}`)
		return
	}
	ch := make(chan string)
	data := make(chan *loraResponse)
	exit := make(chan bool)
	go getSSE(fmt.Sprintf("%s/events", baseURL), body, ch)
	go func() {
		for {
			isExit := false
			select {
			case str := <-ch:
				fmt.Println(str)
				if msg, err := Decode(str); err == nil {
					fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
					if msg.Result.Name == "as.up.data.forward" && msg.Result.Data.Type == "type.googleapis.com/ttn.lorawan.v3.ApplicationUp" {
						data <- msg
					}
				} else {
					fmt.Println("decode error", err)
				}
			case <-exit:
				isExit = true
			}
			if isExit {
				break
			}
		}
	}()
	for {
		isExit := false
		select {
		case msg := <-data:
			isExit = true
			fmt.Println("------------>>>>>", msg)
			ctx.JSON(200, msg)
		case <-time.After(time.Second * 60 * 5):
			isExit = true
		}
		if isExit {
			ctx.Text(200, "")
			go func() {
				exit <- true
			}()
			break
		}
	}
}

type webscoket struct {
	read  chan mvc.WebSocketData
	write chan mvc.WebSocketData
	close chan bool
}
type userLink struct {
	userID string
	ws     webscoket
}

var userWebsocket map[string][]userLink
var lock sync.Mutex

func addGroupUser(gid string, ul userLink) {
	lock.Lock()
	if userWebsocket[gid] == nil {
		userWebsocket[gid] = make([]userLink, 0)
	}
	userWebsocket[gid] = append(userWebsocket[gid], ul)
	lock.Unlock()
	go listen(gid, ul)
}
func listen(gid string, ul userLink) {
	for {
		select {
		case wd, ok := <-ul.ws.read:
			if !ok {
				goto Exit
			}
			recvData := string(wd.DataValue)
			if recvData == "close" {
				ul.ws.close <- true
				goto Exit
			} else {
				//fmt.Printf("收到客户端数据：%s\n", string(wd.DataValue))
				//replyData := fmt.Sprintf("%s说[%s]", ul.userID, recvData)
				//ul.ws.write <- mvc.WebSocketData{DataType: 1, DataValue: []byte(replyData)}
				//broadcast(gid, ul.userID, replyData)
			}
		case <-ul.ws.close:
			goto Exit
		case <-time.After(time.Second * 180):
			ul.ws.write <- mvc.WebSocketData{DataType: 1, DataValue: []byte("180秒没有读到数据（用户发送的数据），开始断开链接")}
			ul.ws.close <- true
			goto Exit
		}
	}
Exit:
	lock.Lock()
	defer lock.Unlock()
	for i, v := range userWebsocket[gid] {
		if v.userID == ul.userID {
			userWebsocket[gid] = append(userWebsocket[gid][:i], userWebsocket[gid][i+1:]...)
		}
	}
}
func broadcast(gid string, uid string, msg string) {
	lock.Lock()
	defer lock.Unlock()
	if k, ok := userWebsocket[gid]; ok {
		for _, v := range k {
			//if v.userID != uid {
			go func(_v userLink) {
				select {
				case _v.ws.write <- mvc.WebSocketData{DataType: 1, DataValue: []byte(msg)}:
					break
				case <-time.After(time.Second * 60 * 2):
					_v.ws.close <- true
				}
			}(v)
			//}
		}
	}
}

func (*ttn) Get_info(ctx *mvc.HTTPContext) {
	bs, _ := json.Marshal(userWebsocket)
	bs2, _ := json.Marshal(sse_conns)
	ctx.Text(200, string(bs)+"            "+string(bs2))
}

func (*ttn) Get_Websocket(ctx *mvc.HTTPContext) {
	if userWebsocket == nil {
		userWebsocket = make(map[string][]userLink)
	}
	body := &requestBody{
		ApiKey:        ctx.ParamURLString("token", ""),
		ApplicationID: ctx.ParamURLString("applicationId", ""),
		DeviceID:      ctx.ParamURLString("deviceId", ""),
	}
	gid := fmt.Sprintf("%s-%s", body.ApplicationID, body.DeviceID)
	uid := time.Now().Format("150405.000000")

	ws := webscoket{
		read:  make(chan mvc.WebSocketData, 1),
		write: make(chan mvc.WebSocketData, 1),
		close: make(chan bool, 1),
	}
	ul := userLink{
		userID: uid,
		ws:     ws,
	}
	addGroupUser(gid, ul)
	go connSSE(body, gid, uid)
	ctx.WebSocket(time.Second*30, true, ws.read, ws.write, ws.close)
}

func Decode(str string) (*loraResponse, error) {
	msg := &loraResponse{}
	err := json.Unmarshal([]byte(str), msg)
	return msg, err
}

type loraResponse struct {
	Result struct {
		Name        string    `json:"name,omitempty"`
		Time        time.Time `json:"time,omitempty"`
		UniqueID    string    `json:"unique_id,omitempty"`
		Identifiers []struct {
			DeviceIDs struct {
				DeviceID       string `json:"device_id,omitempty"`
				ApplicationIDs struct {
					ApplicationID string `json:"application_id,omitempty"`
				} `json:"application_ids"`
				DevEUI  string `json:"dev_eui,omitempty"`
				JoinEUI string `json:"join_eui,omitempty"`
				DevAddr string `json:"dev_addr,omitempty"`
			} `json:"device_ids"`
		} `json:"identifiers,omitempty"`
		Data struct {
			Type         string `json:"@type,omitempty"`
			EndDeviceIDs struct {
				DeviceID       string `json:"device_id,omitempty"`
				ApplicationIDs struct {
					ApplicationID string `json:"application_id,omitempty"`
				} `json:"application_ids"`
				DevEUI  string `json:"dev_eui,omitempty"`
				JoinEUI string `json:"join_eui,omitempty"`
				DevAddr string `json:"dev_addr,omitempty"`
			} `json:"end_device_ids,omitempty"`
			ReceivedAt    time.Time `json:"received_at,omitempty"`
			UplinkMessage struct {
				FPort      int    `json:"f_port,omitempty"`
				FCnet      int    `json:"f_cnt,omitempty"`
				FrmPayload string `json:"frm_payload,omitempty"`
			} `json:"uplink_message,omitempty"`
		} `json:"data,omitempty"`
	} `json:"result,omitempty"`
}
