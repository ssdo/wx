package wx

type QRResult struct {
	Result
	Ticket         string
	Expire_seconds int
	Url            string
	Qr             string
}

func (app *App) MakeTempQRWithId(expire int, scene int) (bool, QRResult) {
	result := QRResult{}
	ok := app.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=", map[string]interface{}{
		"expire_seconds": expire,
		"action_name":    "QR_SCENE",
		"action_info": map[string]map[string]int{
			"scene": {"scene_id": scene},
		},
	}, &result)
	if ok {
		result.Qr = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + result.Ticket
	}
	return ok, result
}

func (app *App) MakeTempQRWithString(expire int, scene string) (bool, QRResult) {
	result := QRResult{}
	ok := app.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=", map[string]interface{}{
		"expire_seconds": expire,
		"action_name":    "QR_STR_SCENE",
		"action_info": map[string]map[string]string{
			"scene": {"scene_str": scene},
		},
	}, &result)
	if ok {
		result.Qr = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + result.Ticket
	}
	return ok, result

}

func (app *App) MakeQRWithId(scene int) (bool, QRResult) {
	result := QRResult{}
	ok := app.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=", map[string]interface{}{
		"action_name": "QR_LIMIT_SCENE",
		"action_info": map[string]map[string]int{
			"scene": {"scene_id": scene},
		},
	}, &result)
	if ok {
		result.Qr = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + result.Ticket
	}
	return ok, result

}

func (app *App) MakeQRWithString(scene string) (bool, QRResult) {
	result := QRResult{}
	ok := app.Post("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=", map[string]interface{}{
		"action_name": "QR_LIMIT_STR_SCENE",
		"action_info": map[string]map[string]string{
			"scene": {"scene_str": scene},
		},
	}, &result)
	if ok {
		result.Qr = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + result.Ticket
	}
	return ok, result
}
