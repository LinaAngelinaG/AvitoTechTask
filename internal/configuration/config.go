package configuration

import (
	"AvitoTechTask/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIp string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"storage"`
}

var instance *Config
var once sync.Once

func GetConfig(logger *logging.Logger) *Config {
	once.Do(func() {
		logger.Info("create and get instance of configuration object")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			needDeb, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(needDeb)
			logger.Fatal(err)
		}
	})
	return instance
}
