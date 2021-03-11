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
		AppKey:       "4df292cff85648c6879cf1923897c531",
		MasterSecret: "q0U2P2m246yVmPaOUN2zEqvm",
	})
	if err != nil {
		t.Error(err)
	}
	notify, err := miClient.Notify(context.Background(), &MessageRequest{
		TargetType:  2.0,
		TargetValue: "CN_fc8d3f2fd34d23a00e2d23d7be8a48f8",
		Notification: Notification{
			Title:           "信用卡激活享好礼",
			Content:         "尊敬的客户，您的京东金融联名卡尚未使用。现惠券等您拿。去激活\\u003e\\u003e",
			ClickActionType: 0,
			ActionParams:    "{\"jumpObject\":\"{\\\"jumpType\\\":2,\\\"jumpP\\\"{\\\\\\\"moduleCode\\\\\\\":\\\\\\\"creditCard\\\\\\\"}\\\",\\\"urlParam\\\":\\\"\\\"}\",\"msgType\":\"389\"}",
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(notify)
}
