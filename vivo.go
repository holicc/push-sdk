package push_sdk

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"push-sdk/http"
	"strconv"
	"strings"
	"time"
)

type VivoMessageRequest struct {
	RegId       string             `json:"regId"`
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	SkipType    int                `json:"skipType"`
	SkipContent string             `json:"skipContent"`
	RequestId   string             `json:"requestId"`
	NotifyType  int                `json:"notifyType"`
	Extra       *SingleNotifyExtra `json:"extra,omitempty"`
}

type VivoMessageResponse struct {
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

type AuthTokenReq struct {
	AppId     string `json:"appId"`
	AppKey    string `json:"appKey"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
}

type AuthTokenResp struct {
	Result    int       `json:"result"`    // 0 成功，非0失败
	Desc      string    `json:"desc"`      // 文字描述接口调用情况
	AuthToken string    `json:"authToken"` // 默认有效一天
	Time      time.Time `json:"-"`
}

type VivoClient struct {
	Vivo   Vivo
	Token  *AuthTokenResp
	client *http.HTTPClient
}

func NewVivoClient(vi Vivo) (*VivoClient, error) {
	if vi.AppPkgName == "" {
		return nil, errors.New("app pkg-name empty")
	}
	if vi.AppId == "" {
		return nil, errors.New("app id empty")
	}
	if vi.AppKey == "" {
		return nil, errors.New("app key empty")
	}
	if vi.AppSecret == "" {
		return nil, errors.New("app secret empty")
	}

	client := http.NewHTTPClient()

	return &VivoClient{
		Vivo:   vi,
		client: client,
	}, nil
}

func getAccessToken(vi Vivo, client *http.HTTPClient) (*AuthTokenResp, error) {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano()/(1e6), 10)
	authReq := &AuthTokenReq{
		AppId:     vi.AppId,
		AppKey:    vi.AppKey,
		Timestamp: timestamp,
		Sign:      generateSign(vi.AppId, vi.AppKey, timestamp, vi.AppSecret),
	}

	data, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(context.Background(), &http.PushRequest{
		Method: "POST",
		URL:    vi.AuthURL,
		Body:   data,
		Header: nil,
	})
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, errors.New("request access token failed")
	}

	var token AuthTokenResp
	err = json.Unmarshal(resp.Body, &token)
	if err != nil {
		return nil, err
	}
	if token.Result != 0 {
		return nil, errors.New("get access token failed")
	}

	token.Time = time.Now()
	return &token, nil
}

func (v *VivoClient) Notify(ctx context.Context, req MessageRequest) (MessageResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if v.Token == nil || !v.Token.isValid() {
		token, err := getAccessToken(v.Vivo, v.client)
		if err != nil {
			return nil, err
		}
		v.Token = token
	}

	resp, err := v.client.Do(ctx, &http.PushRequest{
		Method: "POST",
		URL:    v.Vivo.PushURL,
		Body:   data,
		Header: []http.HTTPOption{
			http.SetHeader("Content-Type", "application/json"),
			http.SetHeader("authToken", v.Token.AuthToken),
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, errors.New(fmt.Sprintf("notify request failed %s", string(resp.Body)))
	}

	var r VivoMessageResponse
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (v *VivoMessageRequest) Validate() error {
	if v.RegId == "" {
		return errors.New("reg id empty")
	}
	if v.Content == "" {
		return errors.New("content empty")
	}

	return nil
}

func (v *VivoMessageResponse) GetResult() string {
	return v.Desc
}

func (v *VivoMessageResponse) GetData() map[string]string {
	return nil
}

func (r *AuthTokenResp) isValid() bool {
	return time.Now().Sub(r.Time).Hours() < 24.0
}

func generateSign(appId, appKey, timestamp, sec string) string {
	signStr := appId + appKey + timestamp + sec
	signStr = strings.Trim(signStr, "")
	hash := md5.New()
	hash.Write([]byte(signStr))
	return strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
}
