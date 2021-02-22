package xiaomi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"push-sdk"
	http2 "push-sdk/http"
)

type XiaoMiMessageRequest struct {
	Payload               string          `json:"payload"`                 // 消息的内容。（注意：需要对payload字符串做urlencode处理）
	RestrictedPackageName string          `json:"restricted_package_name"` // App的包名
	PassThrough           int             `json:"pass_through"`            // 0 表示通知栏消息,1 表示透传消息
	Title                 string          `json:"title"`                   // 通知栏展示的通知的标题
	Description           string          `json:"description"`             // 通知栏展示的通知的描述
	RegistrationId        string          `json:"registration_id"`         // 根据registration_id，发送消息到指定设备上
	Extra                 *push_sdk.Extra `json:"extra"`
}

type XiaoMiMessageResponse struct {
	Result      string            `json:"result"`      // "ok" 表示成功
	Description string            `json:"description"` // 对发送消息失败原因的解释
	Data        map[string]string `json:"data"`        // 本身就是一个json字符串（其中id字段的值就是消息的Id）
	Code        int               `json:"code"`        // 0表示成功，非0表示失败
	Info        string            `json:"info"`        // 详细信息
	Reason      string            `json:"reason"`      // 错误原因
}

type XiaoMiClient struct {
	Mi     push_sdk.XiaoMi
	client *http2.HTTPClient
}

func NewXiaoMiClient(mi push_sdk.XiaoMi) (*XiaoMiClient, error) {
	if mi.AppPkgName == "" {
		return nil, errors.New("app pkg-name empty")
	}
	if mi.AppSecret == "" {
		return nil, errors.New("app secret empty")
	}
	return &XiaoMiClient{
		Mi:     mi,
		client: http2.NewHTTPClient(),
	}, nil
}

func (c *XiaoMiClient) Notify(ctx context.Context, body push_sdk.MessageRequest) (push_sdk.MessageResponse, error) {
	if e := body.Validate(); e != nil {
		return nil, e
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(ctx, &http2.PushRequest{
		Method: http.MethodPost,
		URL:    c.Mi.PushURL,
		Body:   data,
		Header: []http2.HTTPOption{
			http2.SetHeader("Authorization", fmt.Sprintf("key=%s", c.Mi.AppSecret)),
			http2.SetHeader("Content-Type", "application/x-www-form-urlencoded"),
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Status == 200 {
		var r XiaoMiMessageResponse
		err = json.Unmarshal(resp.Body, &r)
		if err != nil {
			return nil, err
		}
		return &r, nil
	}
	return nil, errors.New(fmt.Sprintf("response status %v %s", resp.Status, string(resp.Body)))
}

func (p *XiaoMiMessageRequest) Validate() error {
	if p.Title == "" {
		return errors.New("message title is empty")
	}
	if p.RegistrationId == "" {
		return errors.New("message registration id is empty")
	}
	if p.Payload == "" {
		return errors.New("message playload is empty")
	}
	if p.PassThrough != 0 && p.PassThrough != 1 {
		return errors.New("unknown message pass type")
	}

	return nil
}

func (p *XiaoMiMessageResponse) GetResult() string {
	return p.Result
}

func (p *XiaoMiMessageResponse) GetData() map[string]string {
	return p.Data
}
