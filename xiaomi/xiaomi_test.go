package xiaomi

import (
	"context"
	"fmt"
	sdk "github.com/holicc/push-sdk"
	"testing"
)

func TestXiaomiNotify(t *testing.T) {
	miClient, err := NewXiaoMiClient(sdk.XiaoMi{
		Platform: sdk.Platform{
			PushURL: "https://api.xmpush.xiaomi.com/v4/message/regid",
			AuthURL: "",
		},
		AppPkgName: "",
		AppSecret:  "",
	})
	if err != nil {
		t.Error(err)
	}
	notify, err := miClient.Notify(context.Background(), &MessageRequest{
		Payload:               "",
		RestrictedPackageName: "",
		PassThrough:           0,
		NotifyType:            1,
		Title:                 "",
		Description:           "",
		RegistrationId:        "",
		Extra:                 nil,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(notify)
}
