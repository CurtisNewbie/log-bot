package logbot

import (
	"github.com/curtisnewbie/gocommon/bus"
	"github.com/curtisnewbie/gocommon/common"
)

func BeforeServerBootstrapp(c common.ExecContext) error {
	if e := bus.DeclareEventBus(ERROR_LOG_EVENT_BUS); e != nil {
		return e
	}
	return bus.SubscribeEventBus(ERROR_LOG_EVENT_BUS, 2, func(l LogLineEvent) error {
		ec := common.EmptyExecContext()
		return SaveErrorLog(ec, l)
	})
}

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
