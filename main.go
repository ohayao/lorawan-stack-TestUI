package main

import (
	"bytes"

	"github.com/ohayao/common/http"
	"github.com/ohayao/mvc"
)

func main() {
	mvc.Router.RegisterStatic("/", "./ui/dist", "index.html", nil)
	mvc.Router.RegisterController("v1", "", "api", &api{})
	mvc.Router.RegisterController("v2", "", "ttn", &ttn{})
	mvc.Run(3001)
}

// func test() {
// 	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/g2getstatus")
// 	h.Header("content-type", "application/json")
// 	res := h.PostBytes(bytes.NewBufferString(`{
// 		"apiKey": "nErrIBbGzdrS6iiHkUd0z9.4v18I8b7Pn5zD8Y8oOkZTsPgk2814f9IwJib1PMn",
// 		"bridgeId": "EB1X01875036",
// 		"lockId": "IGP108cc0502",
// 		"expiryDuration": "15"
// 	}`))
// 	fmt.Println(res, string(res))
// }

type api struct{}

func (*api) Post_GetStatus(ctx *mvc.HTTPContext) {
	body := ctx.GetBody()
	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/g2getstatus")
	h.Header("content-type", "application/json")
	ctx.HTML(200, string(h.PostBytes(bytes.NewBuffer(body))))
}
func (*api) Post_Lock(ctx *mvc.HTTPContext) {
	body := ctx.GetBody()
	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/g2lock")
	h.Header("content-type", "application/json")
	ctx.HTML(200, string(h.PostBytes(bytes.NewBuffer(body))))
}
func (*api) Post_Unlock(ctx *mvc.HTTPContext) {
	body := ctx.GetBody()
	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/g2unlock")
	h.Header("content-type", "application/json")
	ctx.HTML(200, string(h.PostBytes(bytes.NewBuffer(body))))
}
func (*api) Post_GetBridgeLogs(ctx *mvc.HTTPContext) {
	body := ctx.GetBody()
	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/getclientresponse")
	h.Header("content-type", "application/json")
	ctx.HTML(200, string(h.PostBytes(bytes.NewBuffer(body))))
}
func (*api) Post_AddLock(ctx *mvc.HTTPContext) {
	body := ctx.GetBody()
	h := http.New("https://p1k84eekmf.execute-api.ap-southeast-1.amazonaws.com/development/addlocks")
	h.Header("content-type", "application/json")
	ctx.HTML(200, string(h.PostBytes(bytes.NewBuffer(body))))
}
