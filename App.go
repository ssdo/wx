package wx

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
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

type Result struct {
	Errcode int
	Errmsg  string
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
			} else {
				app.logger.Error("failed to get access token", "errcode", tokenResult.Errcode, "errmsg", tokenResult.Errmsg)
			}
		} else {
			app.logger.Error("failed to get access token", "err", r.Error.Error())
		}
	}
	app.accessToken = token
	return token
}

func (app *App) CheckSignature(signature string, timestamp string, nonce string) bool {
	args := []string{app.token, timestamp, nonce}
	sort.Strings(args)
	sign := u.Sha1String(strings.Join(args, ""))
	return sign == signature
}

func (app *App) DecodeMessage(signature, timestamp, nonce string, encryptMessage []byte, to interface{}) bool {
	xml1 := struct {
		ToUserName string
		Encrypt    string
	}{}
	err := xml.Unmarshal(encryptMessage, &xml1)
	if err != nil {
		app.logger.Error("[WX_MESSAGE] failed to decode message", "body", string(encryptMessage), "err", err.Error())
	}
	if err == nil { // && xml1.ToUserName == app.name
		args := []string{app.token, timestamp, nonce, xml1.Encrypt}
		sort.Strings(args)
		sign := u.Sha1String(strings.Join(args, ""))
		if sign == signature {
			aesKey, _ := base64.StdEncoding.DecodeString(app.aesKey + "=")
			buf := u.DecryptAesBytes(xml1.Encrypt, aesKey, nil)
			if len(buf) > 20 {
				var size int32
				_ = binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &size)
				if len(buf) >= 20+int(size) {
					err = xml.Unmarshal(buf[20:20+size], &to)
					if err != nil {
						app.logger.Error("[WX_MESSAGE] failed to decode message", "message", string(buf[20:20+size]))
					} else {
						return true
					}
				} else {
					app.logger.Error("[WX_MESSAGE] bad message size", "size", len(buf))
				}
			} else {
				app.logger.Error("[WX_MESSAGE] bad message size", "size", len(buf))
			}
		} else {
			app.logger.Error("[WX_MESSAGE] failed to verify signature", "sign", sign, "signature", signature)
		}
	}
	return false
}

func (app *App) Get(url string, out interface{}) bool {
	if token := app.GetAccessToken(); token != "" {
		c := conf.httpClient
		err := c.Get(url + token).To(out)

		if err == nil {
			return true
		}
		app.logger.Error("failed to call get", "err", err.Error(), "url", url, "appName", app.name)
	}

	return false
}

func (app *App) Post(url string, data interface{}, out interface{}) bool {
	if token := app.GetAccessToken(); token != "" {
		c := conf.httpClient
		r := c.Post(url+token, data)
		err := r.To(out)
		if err == nil {
			return true
		}
		app.logger.Error("failed to call post", "err", err.Error(), "url", url, "data", data, "appName", app.name)
	}

	return false
}

func (app *App) DoGet(url string) bool {
	result := Result{}
	if app.Get(url, &result) && result.Errcode == 0 {
		return true
	}
	app.logger.Error("failed to call do get", "appName", app.name, "errcode", result.Errcode, "errmsg", result.Errmsg)
	return false
}

func (app *App) DoPost(url string, data interface{}) bool {
	result := Result{}
	if app.Post(url, data, &result) && result.Errcode == 0 {
		return true
	}
	app.logger.Error("failed to call do post", "appName", app.name, "errcode", result.Errcode, "errmsg", result.Errmsg)
	return false
}
