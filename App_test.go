package wx_test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ssdo/wx"
	"github.com/ssgo/config"
	"testing"
)

var testConfig = struct {
	TemplateId string
	Openid     string
	Apps       map[string]*wx.AppConfig
}{}

func init() {
	config.LoadConfig("test", &testConfig)
	wx.Init(wx.Config{
		Apps: testConfig.Apps,
	})
}

func TestAccessToken(t *testing.T) {
	app := wx.GetApp("test", nil)
	if app == nil {
		t.Fatal("app test not exists")
	} else {
		token := app.GetAccessToken()
		if token == "" {
			t.Fatal("failed to get access token")
		}
	}
}

//func TestTplMessage(t *testing.T) {
//	app := wx.GetApp("test", nil)
//	//app.SendTemplateMessage(testConfig.TemplateId, testConfig.Openid, map[string]string{
//	//	"name": "张三",
//	//}, "")
//
//	app.SendTemplateMessageWithColor(testConfig.TemplateId, testConfig.Openid, "#0099ff", map[string]string{
//		"name": "张三",
//	}, map[string]string{
//		"name": "#ff9900",
//	}, "https://www.baidu.com")
//}
