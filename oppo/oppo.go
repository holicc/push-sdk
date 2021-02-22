package oppo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"push-sdk"
	"push-sdk/http"
	"strconv"
	"strings"
	"time"
)

type MessageRequest struct {
	Message   *SingleMessage `json:"message"`
	AuthToken string         `json:"auth_token"`
}

type SingleMessage struct {
	TargetType   int          `json:"target_type"`
	TargetValue  string       `json:"target_value"`
	Notification Notification `json:"notification"`
}

type Notification struct {
	Title             string `json:"title"`
	SubTitle          string `json:"sub_title"`
	Content           string `json:"content"`
	ClickActionType   int    `json:"click_action_type"`
	ClickActionUrl    string `json:"click_action_url"`
	CallBackUrl       string `json:"call_back_url"`
	CallBackParameter string `json:"call_back_parameter"`
}

type Response struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *PushMessageData `json:"data"`
}

type PushMessageData struct {
	BroadcastMessageId string `json:"message_id"`
	SingleMessageId    string `json:"messageId"`
	Status             string `json:"status"`
	TaskId             string `json:"task_id"`
}

type TokenResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AuthToken  string    `json:"auth_token"`
		CreateTime time.Time `json:"create_time"`
	} `json:"data"`
}

type TokenInfo struct {
	AppKey       string
	MasterSecret string
	AuthURL      string

	Token      string
	CreateTime time.Time
}

type client struct {
	httpclient *http.HTTPClient
	authClient *http.AuthClient

	op push_sdk.Oppo
}

func NewOppoClient(op push_sdk.Oppo) (*client, error) {
	if op.AppPkgName == "" {
		return nil, errors.New("app pkg-name empty")
	}
	if op.AppKey == "" {
		return nil, errors.New("app key empty")
	}
	if op.MasterSecret == "" {
		return nil, errors.New("master secret empty")
	}

	return &client{
		httpclient: http.NewHTTPClient(),
		authClient: http.NewAuthClient(&TokenInfo{
			AppKey:       op.AppKey,
			MasterSecret: op.MasterSecret,
			AuthURL:      op.AuthURL,
		}),
		op: op,
	}, nil
}

func (o *client) Notify(ctx context.Context, req push_sdk.MessageRequest) (push_sdk.MessageResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	v.Add("message", string(data))

	token, err := o.authClient.GetAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := o.httpclient.Do(ctx, &http.PushRequest{
		Method: "POST",
		URL:    o.op.PushURL,
		Body:   []byte(v.Encode()),
		Header: []http.HTTPOption{
			http.SetHeader("Content-Type", "application/x-www-form-urlencoded"),
			http.SetHeader("auth_token", token),
		},
	})

	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, errors.New(fmt.Sprintf("notify request failed %s", string(resp.Body)))
	}

	var r Response
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (m *Response) GetResult() string {
	return m.Message
}

func (m *Response) GetData() map[string]string {
	return nil
}

func (o *MessageRequest) Validate() error {
	return nil
}

func (t *TokenInfo) TokenRequest() ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano()/(1e6), 10)

	signStr := t.AppKey + timestamp + t.MasterSecret
	signStr = strings.Trim(signStr, "")
	hash := md5.New()
	hash.Write([]byte(signStr))
	sign := strings.ToLower(hex.EncodeToString(hash.Sum(nil)))

	req := map[string]string{
		"app_key":   t.AppKey,
		"timestamp": timestamp,
		"sign":      sign,
	}
	data := url.Values{}
	for key, value := range req {
		data.Add(key, value)
	}

	return []byte(data.Encode()), nil
}

func (t *TokenInfo) ParseResponse(d []byte) error {
	var token TokenResp
	err := json.Unmarshal(d, &token)
	if err != nil {
		return err
	}

	if token.Code == 0 {
		t.Token = token.Data.AuthToken
		t.CreateTime = token.Data.CreateTime
		return nil
	}

	return errors.New("response parse success but can not obtain token")
}

func (t *TokenInfo) SetHeader() []http.HTTPOption {
	return []http.HTTPOption{
		http.SetHeader("Content-Type", "application/x-www-form-urlencoded"),
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
