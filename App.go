package wx

import (
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/u"
	"sort"
	"strings"
)

type App struct {
	name            string
	id              string
	secret          string
	token           string
	aesKey          string
	accessToken     string
	authAccessToken string
	logger          *log.Logger
}

func GetApp(name string, logger *log.Logger) *App {
	if logger == nil {
		logger = log.DefaultLogger
	}

	appConf := conf.Apps[name]
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
		token:  appConf.Token,
		aesKey: appConf.AesKey,
		logger: logger,
	}
}

func (app *App) GetAccessToken() string {
	if app.accessToken != "" {
		return app.accessToken
	}

	rd := conf.Redis
	token := rd.GET("WX_TOKEN_" + app.name).String()
	if token == "" {
		c := conf.httpClient
		r := c.Get(fmt.Sprint("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=", app.id, "&secret=", app.secret))
		if r.Error == nil {
			tokenResult := struct {
				Errcode      int
				Errmsg       string
				Access_token string
				Expires_in   int
			}{}

			_ = r.To(&tokenResult)
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
	app.accessToken = token
	return token
}

func (app *App) CheckSignature(signature string, timestamp int64, nonce int) bool {
	args := []string{app.token, u.String(timestamp), u.String(nonce)}
	sort.Strings(args)
	sign := u.Sha1String(strings.Join(args, ""))
	return sign == signature
}
