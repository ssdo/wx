package wx

func (app *App) SetMenu(data interface{}) bool {
	return app.DoPost("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=", data)
}

func (app *App) GetMenu(out interface{}) bool {
	return app.Get("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=", out)
}

func (app *App) DeleteMenu() bool {
	return app.DoGet("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=")
}

func (app *App) AddConditionalMenu(data interface{}) (bool, string) {
	result := struct {
		Result
		Menuid  string
	}{}
	ok := app.Post("https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=", data, &result)
	if result.Errcode != 0 {
		app.logger.Error("failed to add conditional menu", "appName", app.name, "errcode", result.Errcode, "errmsg", result.Errmsg)
	}
	return ok && result.Errcode == 0, result.Menuid
}

func (app *App) DeleteConditionalMenu(menuid string) bool {
	return app.DoPost("https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token=", map[string]string{"menuid": menuid})
}

func (app *App) TestConditionalMenu(userId string, result interface{}) bool {
	return app.Post("https://api.weixin.qq.com/cgi-bin/menu/trymatch?access_token=", map[string]string{"user_id": userId}, result)
}
