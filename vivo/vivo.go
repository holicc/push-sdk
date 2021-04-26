package vivo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/holicc/push-sdk"
	"github.com/holicc/push-sdk/http"
	"strconv"
	"strings"
	"time"
)

type MessageRequest struct {
	RegId           string                 `json:"regId"`
	Title           string                 `json:"title"`
	Content         string                 `json:"content"`
	Expire          int64                  `json:"timeToLive"`
	SkipType        int                    `json:"skipType"`
	SkipContent     string                 `json:"skipContent"`
	RequestId       string                 `json:"requestId"`
	NotifyType      int                    `json:"notifyType"`
	PushMode        int                    `json:"pushMode"`
	ClientCustomMap map[string]interface{} `json:"clientCustomMap"`
	Extra           *SingleNotifyExtra     `json:"extra,omitempty"`
}

type MessageResponse struct {
	Result       int         `json:"result"`    // 0 表示成功，非0失败
	Desc         string      `json:"desc"`      // 文字描述接口调用情况
	RequestId    string      `json:"requestId"` // 请求ID
	InvalidUsers interface{} `json:"invalidUsers"`
	TaskId       string      `json:"taskId"` // 任务ID
}

type SingleNotifyExtra struct {
	CallBack      string `json:"callback,omitempty"`
	CallBackParam string `json:"callback.param,omitempty"`
}

type TokenInfo struct {
	AppId     string
	AppKey    string
	AppSecret string
	AuthURL   string

	Token      string
	CreateTime time.Time
}

type AuthTokenReq struct {
	AppId     string `json:"appId"`
	AppKey    string `json:"appKey"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
}

type AuthTokenResp struct {
	Result    int    `json:"result"`    // 0 成功，非0失败
	Desc      string `json:"desc"`      // 文字描述接口调用情况
	AuthToken string `json:"authToken"` // 默认有效一天
}

type client struct {
	vi     sdk.Vivo
	client *http.HTTPClient

	authClient *http.AuthClient
}

func NewVivoClient(vi sdk.Vivo) (*client, error) {
	if vi.AppId == "" {
		return nil, errors.New("app id empty")
	}
	if vi.AppKey == "" {
		return nil, errors.New("app key empty")
	}
	if vi.AppSecret == "" {
		return nil, errors.New("app secret empty")
	}

	return &client{
		vi:     vi,
		client: http.NewHTTPClient(),
		authClient: http.NewAuthClient(&TokenInfo{
			AppId:     vi.AppId,
			AppKey:    vi.AppKey,
			AppSecret: vi.AppSecret,
			AuthURL:   vi.AuthURL,
		}),
	}, nil
}

func (v *client) Notify(ctx context.Context, req sdk.MessageRequest) (sdk.MessageResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	token, err := v.authClient.GetAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(ctx, &http.PushRequest{
		Method: "POST",
		URL:    v.vi.PushURL,
		Body:   data,
		Header: []http.HTTPOption{
			http.SetHeader("Content-Type", "application/json"),
			http.SetHeader("authToken", token),
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, errors.New(fmt.Sprintf("notify request failed %s", string(resp.Body)))
	}

	var r MessageResponse
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (v *MessageRequest) Validate() error {
	if v.RegId == "" {
		return errors.New("reg id empty")
	}
	if v.Content == "" {
		return errors.New("content empty")
	}

	return nil
}

func (v *MessageRequest) GetRequestBody() ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *MessageResponse) GetResult() string {
	return strconv.Itoa(v.Result)
}

func (v *MessageResponse) GetData() map[string]string {
	return nil
}

func (t *TokenInfo) TokenRequest() ([]byte, error) {
	t.CreateTime = time.Now()
	timestamp := strconv.FormatInt(t.CreateTime.UTC().UnixNano()/(1e6), 10)
	authReq := &AuthTokenReq{
		AppId:     t.AppId,
		AppKey:    t.AppKey,
		Timestamp: timestamp,
		Sign:      generateSign(t.AppId, t.AppKey, timestamp, t.AppSecret),
	}

	data, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *TokenInfo) ParseResponse(data []byte) error {
	var token AuthTokenResp
	err := json.Unmarshal(data, &token)
	if err != nil {
		return err
	}
	if token.Result != 0 {
		return errors.New("get access token failed")
	}
	t.Token = token.AuthToken

	return nil
}

func (t *TokenInfo) SetHeader() []http.HTTPOption {
	return []http.HTTPOption{
		http.SetHeader("Content-Type", "application/json"),
	}
}

func (t *TokenInfo) GetAuthUrl() string {
	return t.AuthURL
}

func (t *TokenInfo) GetAuthMethod() string {
	return "POST"
}

func (t *TokenInfo) GetAccessToken() string {
	return t.Token
}

func (t *TokenInfo) IsValidate() bool {
	return time.Now().Sub(t.CreateTime).Hours() < 24.0
}

func generateSign(appId, appKey, timestamp, sec string) string {
	signStr := appId + appKey + timestamp + sec
	signStr = strings.Trim(signStr, "")
	hash := md5.New()
	hash.Write([]byte(signStr))
	return strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
}
