package wx

import (
	"github.com/ssdo/utility"
	"github.com/ssgo/db"
	"github.com/ssgo/httpclient"
	"github.com/ssgo/redis"
	"time"
)

var inited = false

var OpenidLimiter *utility.Limiter

//// 用户表
//type OpenIdTable struct {
//	Table string // 表名
//	Id    string // id字段名
//	Name     string // 用户名字段名
//}

type AppConfig struct {
	Id     string
	Secret string
}

var Config = struct {
	httpClient          *httpclient.ClientPool
	Redis               *redis.Redis  // Redis连接池
	DB                  *db.DB        // 数据库连接池
	ApiTimeoutDuration  time.Duration // 调用微信API超时时间
	OpenidLimitDuration time.Duration // Openid限制器时间间隔
	OpenidLimitTimes    int           // Openid限制器时间单位内允许的次数
	//UserTable                UserTable                                               // 数据库用户表配置
	Apps map[string]*AppConfig
}{
	Redis:               nil,
	DB:                  nil,
	ApiTimeoutDuration:  30 * time.Second,
	OpenidLimitDuration: 5 * time.Minute,
	OpenidLimitTimes:    10000,
	//UserTable: UserTable{
	//	Table: "User",
	//	Id:    "id",
	//	//Name:     "name",
	//	Phone:    "phone",
	//	Password: "password",
	//	Salt:     "salt",
	//},
}

func Init() {
	if inited {
		return
	}
	inited = true

	if Config.Redis == nil {
		Config.Redis = redis.GetRedis("wx", nil)
	}
	if Config.DB == nil {
		Config.DB = db.GetDB("wx", nil)
	}
	OpenidLimiter = utility.NewLimiter("Openid", Config.OpenidLimitDuration, Config.OpenidLimitTimes, Config.Redis)

	Config.httpClient = httpclient.GetClient(Config.ApiTimeoutDuration)
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
