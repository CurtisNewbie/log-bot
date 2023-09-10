package logbot

import "github.com/curtisnewbie/miso/miso"

const (
	PROP_ENABLE_REMOVE_ERR_LOG_TASK = "miso.remove-error-log"
)

func init() {
	miso.SetDefProp(PROP_ENABLE_REMOVE_ERR_LOG_TASK, false)
}

type WatchConfig struct {
	App  string
	File string
	Type string
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
	miso.UnmarshalFromProp(&conf)
	return conf
}

func IsRmErrorLogTaskEnabled() bool {
	return miso.GetPropBool(PROP_ENABLE_REMOVE_ERR_LOG_TASK)
}
