package conf

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"

	"igen/lib/logger"

	"igen/lib"
	"igen/lib/consul"
)

const (
	// EnvDev 开发环境
	EnvDev string = "dev"

	// EnvTest 测试环境
	EnvTest string = "test"

	// EnvProd 生产环境
	EnvProd string = "prod"
)

// App 实例化App信息
var App *Conf

// Conf 配置信息
type Conf struct {
	Env           string `toml:"env"` // 当前环境 dev, test, prod
	LocalAddress  string `toml:"localAddress"`
	PublicAddress string `toml:"publicAddress"`
	HTTPAddress   string `toml:"httpAddress"`
	HTTPPort      int    `toml:"httpPort"`
	ApiHost       string `toml:"apiHost"`
	LogLevel      string `toml:"logLevel"`
	LogFilename   string `toml:"logFilename"`

	AppSecrets map[string]string    `toml:"appSecret"` // 各个端对应的AppSecret
	Consul     consulConf           `toml:"consul"`
	Redis      map[string]redisConf `toml:"redis"` // redis配置

	useFileLogger bool
	redisPool     *redis.Pool // IM redis连接池
}

type consulConf struct {
	Address  string                `toml:"address"`
	Services []*consul.ServiceConf `toml:"services"`
}

// RedisConf Redis配置信息
type redisConf struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

func init() {
	App = &Conf{}
}

// Init 初始化配置
func Init() {
	if !IsTest() {
		parseFlag()
		dialConsul()
	}

	loadTomlConf()
	loadFlagConf()

	initLogger()
	initConsul() // initConsul again

	if IsDev() {
		logger.DebugAsJSON(App)
	}
}

func loadTomlConf() {
	r, err := parseToml(RealFilePath("conf/app.toml"))
	if err != nil {
		panic(err)
	}

	bakEnv := App.Env
	_, err = toml.DecodeReader(r, App)
	if err != nil {
		panic(err.Error())
	}
	if bakEnv != "" {
		App.Env = bakEnv
	}
}

// SetEnv set current env
func SetEnv(env string) {
	App.Env = env
}

// IsProd returns whether current env is Production
func IsProd() bool {
	return App.Env == EnvProd
}

// IsDev returns whether current env is Development
func IsDev() bool {
	return App.Env == EnvDev
}

// IsTest returns whether current env is Testing
func IsTest() bool {
	return App.Env == EnvTest
}

func initLogger() {
	if IsProd() || App.useFileLogger {
		logger.SetAppName(loggerAppName())
		logger.SetLevel(App.LogLevel)
		logger.SetWriter(logger.NewFileWriter(RealFilePath(App.LogFilename)))
	}
}

func loggerAppName() string {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	a := strings.Split(parts[0], "/")
	return strings.Join(a[0:len(a)-1], ".")
}

func dialConsul() {
	if App.Consul.Address == "" {
		App.Consul.Address = consulAddress()
	}
	consul.Dial(consul.Conf{Address: App.Consul.Address})
}

func initConsul() {
	dialConsul()

	var err error
	for _, ser := range App.Consul.Services {
		err = consul.Register(ser)
		if err != nil && IsProd() {
			panic(err)
		}
	}
}

// RedisPool returns the default redis pool
func RedisPool() *redis.Pool {
	if App.redisPool == nil {
		redisConf := getRedisConf()
		App.redisPool = lib.RedisPool(redisConf.Addr, redisConf.Password, redisConf.DB)
	}
	return App.redisPool
}

// getRedisConf 获取Redis配置信息
func getRedisConf() (conf redisConf) {
	if info, ok := App.Redis[App.Env]; ok {
		conf = info
	} else {
		panic(fmt.Sprintf("找不到redis在%s环境的配置信息", App.Env))
	}
	return
}
