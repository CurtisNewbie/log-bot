package logbot

import "github.com/curtisnewbie/gocommon/common"

const (
	PROP_ENABLE_REMOVE_ERR_LOG_TASK = "task.remove-error-log"
)

func init() {
	common.SetDefProp(PROP_ENABLE_REMOVE_ERR_LOG_TASK, false)
}

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

func IsRmErrorLogTaskEnabled() bool {
	return common.GetPropBool(PROP_ENABLE_REMOVE_ERR_LOG_TASK)
}
