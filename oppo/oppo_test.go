package oppo

import (
	"context"
	"fmt"
	sdk "github.com/holicc/push-sdk"
	"testing"
)

func TestOPPONotify(t *testing.T) {
	miClient, err := NewOppoClient(sdk.Oppo{
		Platform: sdk.Platform{
			PushURL: "https://api.push.oppomobile.com/server/v1/message/notification/unicast",
			AuthURL: "https://api.push.oppomobile.com/server/v1/auth",
		},
		AppKey:       "",
		MasterSecret: "",
	})
	if err != nil {
		t.Error(err)
	}
	notify, err := miClient.Notify(context.Background(), &MessageRequest{
		TargetType:  2.0,
		TargetValue: "",
		Notification: Notification{
			Title:             "Hello World",
			SubTitle:          "hello",
			Content:           "hello world hello world hello world hello world hello world hello world hello world hello world ",
			ClickActionType:   0,
			ClickActionUrl:    "",
			CallBackUrl:       "",
			CallBackParameter: "",
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(notify)
}
