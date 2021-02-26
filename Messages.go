package wx

func (app *App) SendTemplateMessage(templateId, to string, args map[string]string, url string) int {
	return app.SendTemplateMessageWithMiniProgram(templateId, to, "", args, nil, url, "", "")
}

func (app *App) SendTemplateMessageWithColor(templateId, to, textColor string, args map[string]string, argsColor map[string]string, url string) int {
	return app.SendTemplateMessageWithMiniProgram(templateId, to, textColor, args, argsColor, url, "", "")
}

func (app *App) SendTemplateMessageWithMiniProgram(templateId, to, textColor string, args map[string]string, argsColor map[string]string, url, mpId, mpPagePath string) int {
	token := app.GetAccessToken()
	data := make(map[string]map[string]string)
	for k, v := range args {
		data[k] = make(map[string]string)
		data[k]["value"] = v
		if argsColor != nil && argsColor[k] != "" {
			data[k]["color"] = argsColor[k]
		}
	}

	c := Config.httpClient
	postData := map[string]interface{}{
		"touser":      to,
		"template_id": templateId,
		"url":         url,
		"data":        data,
	}

	if textColor != "" {
		postData["color"] = textColor
	}

	if mpId != "" {
		postData["miniprogram"] = map[string]string{
			"appid":    mpId,
			"pagepath": mpPagePath,
		}
	}

	result := struct {
		Errcode int
		Errmsg  string
		Msgid   int
	}{}
	err := c.Post("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="+token, postData).To(&result)

	if err != nil {
		app.logger.Error("failed to send wx template message", "err", err.Error(), "templateId", templateId, "appName", app.name)
	}

	if result.Msgid == 0 || result.Errcode != 0 {
		app.logger.Error("failed to send wx template message", "templateId", templateId, "appName", app.name, "errcode", result.Errcode, "errmsg", result.Errmsg)
	}

	return result.Msgid
}
