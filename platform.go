package push_sdk

import "context"

type MessageRequest interface {
	Validate() error
}

type MessageResponse interface {
	GetResult() string
	GetData() map[string]string
}

type PushClient interface {
	Notify(ctx context.Context, req MessageRequest) (MessageResponse, error)
}

type PushConfig struct {
	XiaoMi `json:"xiaomi"`
	MeiZu  `json:"meizu"`
	Oppo   `json:"oppo"`
	Vivo   `json:"vivo"`
}

type Message struct {
	BusinessId    string            `json:"businessId"`    // 业务ID
	Title         string            `json:"title"`         // 标题，建议不超过10个汉字
	SubTitle      string            `json:"subTitle"`      // 副标题，建议不超过10个汉字
	Content       string            `json:"content"`       // 内容，建议不超过20个汉字
	Extra         map[string]string `json:"extra"`         // 自定义消息。只支持一维
	CallBack      string            `json:"callback"`      // 送达回执地址，供推送厂商调用，最大128字节
	CallbackParam string            `json:"callbackParam"` // 自定义回执参数
}

type Extra struct {
	NotifyEffect string `json:"extra.notify_effect"` // 预定义通知栏消息的点击行为
	IntentUri    string `json:"extra.intent_uri"`
}

type Platform struct {
	PushURL string
	AuthURL string
}

type XiaoMi struct {
	Platform
	AppPkgName string `json:"appPkgName"`
	AppSecret  string `json:"appSecret"`
}

type MeiZu struct {
	AppPkgName string `json:"appPkgName"`
	AppId      string `json:"appId"`
	AppSecret  string `json:"appSecret"`
}

type Oppo struct {
	AppPkgName   string `json:"appPkgName"`
	AppKey       string `json:"appKey"`
	MasterSecret string `json:"masterSecret"`
}

type Vivo struct {
	Platform
	AppPkgName string `json:"appPkgName"`
	AppId      string `json:"appId"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
}
