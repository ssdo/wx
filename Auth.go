package wx

import (
	"fmt"
	"net/url"
)

func (app *App) Code2Openid(code string) string {
	r := struct {
		Access_token  string
		Expires_in    int64
		Refresh_token string
		Openid        string
		Scope         string
	}{}
	err := conf.httpClient.Get(fmt.Sprint("https://api.weixin.qq.com/sns/oauth2/access_token?appid=", app.id, "&secret=", app.secret, "&code=", code, "&grant_type=authorization_code")).To(&r)
	if err != nil {
		app.logger.Error(err.Error())
		return ""
	}
	app.authAccessToken = r.Access_token
	return r.Openid
}

func (app *App) MakeBaseAuthUrl(redirectUri, state string) string {
	return app.makeAuthUrl(redirectUri, state, "snsapi_base")
}

func (app *App) MakeUserAuthUrl(redirectUri, state string) string {
	return app.makeAuthUrl(redirectUri, state, "snsapi_userinfo")
}

func (app *App) makeAuthUrl(redirectUri, state, scope string) string {
	return fmt.Sprint(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=", app.id,
		"&redirect_uri=", url.PathEscape(redirectUri),
		"&response_type=code",
		"&scope=", scope,
		"&state=", state,
		"#wechat_redirect",
	)
}
