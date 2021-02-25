package vivo

import (
	"context"
	"fmt"
	sdk "github.com/holicc/push-sdk"
	"testing"
	"time"
)

func TestVivoNotify(t *testing.T) {
	miClient, err := NewVivoClient(sdk.Vivo{
		Platform: sdk.Platform{
			PushURL: "https://api-push.vivo.com.cn/message/send",
			AuthURL: "https://api-push.vivo.com.cn/message/auth",
		},
		AppSecret: "",
		AppKey:    "",
		AppId:     "",
	})
	if err != nil {
		t.Error(err)
	}
	req := MessageRequest{
		RegId:       "",
		Title:       "Hello World",
		Content:     "Hello WorldHello WorldHello WorldHello WorldHello WorldHello WorldHello World",
		SkipType:    1,
		SkipContent: "",
		RequestId:   time.Now().String(),
		NotifyType:  4,
		PushMode:    1,
		Extra:       nil,
	}
	notify, err := miClient.Notify(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(notify)
}
