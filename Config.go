package wx

import (
	"github.com/ssgo/db"
	"github.com/ssgo/httpclient"
	"github.com/ssgo/redis"
	"time"
)

var inited = false

//var OpenidLimiter *utility.Limiter

//// 用户表
//type OpenIdTable struct {
//	Table string // 表名
//	Id    string // id字段名
//	Name     string // 用户名字段名
//}

type AppConfig struct {
	Id     string
	Secret string
	Token  string
	AesKey string
}

type Config struct {
	httpClient          *httpclient.ClientPool
	Redis               *redis.Redis  // Redis连接池
	DB                  *db.DB        // 数据库连接池
	ApiTimeoutDuration  time.Duration // 调用微信API超时时间
	OpenidLimitDuration time.Duration // Openid限制器时间间隔
	OpenidLimitTimes    int           // Openid限制器时间单位内允许的次数
	Apps                map[string]*AppConfig
}

var conf = Config{}

func Init(config Config) {
	if inited {
		return
	}
	inited = true

	conf = config

	if conf.Redis == nil {
		conf.Redis = redis.GetRedis("wx", nil)
	}

	if conf.DB == nil {
		conf.DB = db.GetDB("wx", nil)
	}

	if conf.ApiTimeoutDuration == 0 {
		conf.ApiTimeoutDuration = 30 * time.Second
	}

	if conf.OpenidLimitDuration == 0 {
		conf.OpenidLimitDuration = 5 * time.Minute
	}
	if conf.OpenidLimitTimes == 0 {
		conf.OpenidLimitTimes = 10000
	}
	if conf.Apps == nil {
		conf.Apps = make(map[string]*AppConfig)
	}

	//OpenidLimiter = utility.NewLimiter("Openid", conf.OpenidLimitDuration, conf.OpenidLimitTimes, conf.Redis)

	conf.httpClient = httpclient.GetClient(conf.ApiTimeoutDuration)
}

var settedKey = []byte("?GQ$0K0GgLdO=f+~L68PLm$uhKr4'=tV")
var settedIv = []byte("VFs7@sK61cj^f?HZ")
var keysSetted = false

func SetEncryptKeys(key, iv []byte) {
	if !keysSetted {
		settedKey = key
		settedIv = iv
		keysSetted = true
	}
}
