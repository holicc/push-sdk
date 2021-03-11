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
		AppPkgName: "cn.tf.mobilebank",
		AppSecret:  "SGzgdZ4bt4RgRB7BkieWRA==",
	})
	if err != nil {
		t.Error(err)
	}
	notify, err := miClient.Notify(context.Background(), &MessageRequest{
		Payload:               "{\"jumpObject\":\"{\\\"jumpType\\\":2,\\\"jumpP\\\"{\\\\\\\"moduleCode\\\\\\\":\\\\\\\"creditCard\\\\\\\"}\\\",\\\"urlParam\\\":\\\"\\\"}\",\"msgType\":\"389\"}",
		RestrictedPackageName: "cn.tf.mobilebank",
		PassThrough:           0,
		NotifyType:            1,
		Title:                 "信用卡激活享好礼",
		Description:           "尊敬的客户，您的京东金融联名卡尚未使用。现惠券等您拿。去激活\\u003e\\u003e",
		RegistrationId:        "xNNdlH/OWZgEvkhihVduxuBF4zY+oL55v8QoGD+v46vxqqFNCHMGJLSy0Ycb+6oD",
		Extra:                 nil,
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(notify)
}
