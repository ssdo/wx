package wx

type Tag struct {
	Id    int
	Name  string
	Count int
}

type Userinfo struct {
	Subscribe       uint
	Openid          string
	Nickname        string
	Sex             uint
	Language        string
	City            string
	Province        string
	Country         string
	Headimgurl      string
	Subscribe_time  uint
	Unionid         string
	Remark          string
	Groupid         int
	Tagid_list      []int
	Subscribe_scene string
	Qr_scene        int
	Qr_scene_str    string
}

func (app *App) GetTags() (bool, []Tag) {
	r := struct {
		Tags []Tag
	}{}
	ok := app.Get("https://api.weixin.qq.com/cgi-bin/tags/get?access_token=", &r)
	return ok, r.Tags
}

func (app *App) SetTag(openids []string, tagId int) bool {
	return app.DoPost("https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=", map[string]interface{}{
		"openid_list": openids,
		"tagid":       tagId,
	})
}

func (app *App) RemoveTag(openids []string, tagId int) bool {
	return app.DoPost("https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=", map[string]interface{}{
		"openid_list": openids,
		"tagid":       tagId,
	})
}

func (app *App) GetUserInfo(openid string) (bool, Userinfo) {
	r := Userinfo{}
	ok := app.Get("https://api.weixin.qq.com/cgi-bin/user/info?openid="+openid+"&lang=zh_CN&access_token=", &r)
	return ok, r
}
