package logbot

import "github.com/curtisnewbie/gocommon/common"

type WatchConfig struct {
	App  string `yaml:"app"`
	File string `yaml:"file"`
	Type string `yaml:"type"`
}

type LogBotConfig struct {
	NodeName     string        `mapstructure:"node"`
	WatchConfigs []WatchConfig `mapstructure:"watch"`
}

type Config struct {
	Config LogBotConfig `mapstructure:"logbot"`
}

func LoadLogBotConfig() Config {
	var conf Config
	common.UnmarshalFromProp(&conf)
	return conf
}
