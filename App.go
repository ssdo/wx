package wx

import (
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/u"
)

type App struct {
	name   string
	id     string
	secret string
	logger *log.Logger
}

func GetApp(name string, logger *log.Logger) *App {
	if logger == nil {
		logger = log.DefaultLogger
	}

	appConf := Config.Apps[name]
	if appConf == nil {
		return nil
	}

	secret := u.DecryptAes(appConf.Secret, settedKey, settedIv)
	if secret == "" {
		secret = appConf.Secret
		logger.Warning("wx secret is not encrypted")
	}

	return &App{
		name:   name,
		id:     appConf.Id,
		secret: secret,
		logger: logger,
	}
}

func (app *App) GetAccessToken() string {
	rd := Config.Redis
	token := rd.GET("WX_TOKEN_" + app.name).String()
	if token == "" {
		c := Config.httpClient
		r := c.Get(fmt.Sprint("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=", app.id, "&secret=", app.secret))
		if r.Error == nil {
			tokenResult := struct {
				Errcode      int
				Errmsg       string
				Access_token string
				Expires_in   int
			}{}

			r.To(&tokenResult)
			fmt.Println(">>>>TOKEN", u.JsonP(tokenResult))
			if tokenResult.Errcode == 0 && tokenResult.Access_token != "" {
				token = tokenResult.Access_token
				tokenResult.Expires_in -= 60
				if tokenResult.Expires_in <= 0 {
					tokenResult.Expires_in = 7200 - 60
				}
				if token != "" {
					// 缓存此token
					rd.SETEX("WX_TOKEN_"+app.name, tokenResult.Expires_in, token)
				}
			}
		}
	}
	return token
}
