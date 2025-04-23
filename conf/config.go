package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
	conf *Config
)

type Config struct {
}

type Server struct {
	Addr string
}

func GetConf() *Config {
	once.Do(func() {
		initConfig()
	})

	return conf
}

func initConfig() {
	prefix := "conf"
	env := getEnv()
	path := filepath.Join(prefix, filepath.Join(env, "conf.yaml"))
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf = new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}
}

func getEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "test"
	}

	return env
}
