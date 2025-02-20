package setting

import (
	"flag"
	"log" // 使用标准库的 log

	"github.com/go-ini/ini"
)

type App struct {
	AdminEmail    string
	AdminPassword string
	Host          string
	Port          string
	CorsUrl       string
}

var AppSetting = &App{}

type Database struct {
	Type         string
	User         string
	Password     string
	Host         string
	Port         string
	DBName       string
	SearchPath   string
	MaxIdleConns int
	MaxOpenConns int
}

var DatabaseSetting = &Database{}

type Log struct {
	LogSavePath   string
	LogFileName   string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	Compress      bool
	LogLevel      string
	TimeFormat    string
	ConsoleOutput bool
}

var LogSetting = &Log{}

type Rpc struct {
	ServerHost string
	ServerPort string
}

var RpcSetting = &Rpc{}

type JWT struct {
	Secret     string `ini:"secret"`
	ExpireTime int    `ini:"expire_time"`
}

var JWTSetting = &JWT{}

type RPC struct {
	LiveServiceAddr string `ini:"live_service_addr"`
}

var RPCSetting = &RPC{}

var (
	cfg        *ini.File
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "conf/app.ini", "path to config file")
}

func Init() {
	// 解析命令行参数
	flag.Parse()

	var err error
	log.Printf("settingh.Setup load config from: %s", configPath)
	cfg, err = ini.Load(configPath)
	if err != nil {
		log.Fatalf("setting.Setup failure, path: %s, error: %v", configPath, err)
		return
	}

	mapTo("app", AppSetting)
	mapTo("database", DatabaseSetting)
	mapTo("log", LogSetting)
	mapTo("rpc", RpcSetting)
	mapTo("jwt", JWTSetting)
	mapTo("rpc", RPCSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("settings.MapTo error, section: %s, error: %v", section, err)
	}
}
