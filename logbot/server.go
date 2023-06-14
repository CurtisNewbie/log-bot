package logbot

import (
	"github.com/curtisnewbie/gocommon/common"
)

func AfterServerBootstrapped(c common.ExecContext) error {
	logBotConfig := LoadLogBotConfig().Config
	for _, wc := range logBotConfig.WatchConfigs {
		go func(w WatchConfig) {
			c := common.EmptyExecContext()
			if e := WatchLogFile(c, w, logBotConfig.NodeName); e != nil {
				c.Log.Errorf("WatchLogFile, app: %v, file: %v, %v", w.App, w.File, e)
			}
		}(wc)
	}
	return nil
}
