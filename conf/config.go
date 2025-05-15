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
	Server       Server       `yaml:"server"`
	Prom         Prom         `json:"prom"`
	StorageClass StorageClass `yaml:"storageClass"`
}

type Server struct {
	Addr string `yaml:"addr"`
}

type Prom struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
}

type StorageClass struct {
	Provisioner []string `yaml:"provisioner"`
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
