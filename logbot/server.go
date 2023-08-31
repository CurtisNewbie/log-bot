package logbot

import (
	"time"

	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/bus"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/server"
	"github.com/curtisnewbie/miso/task"
	"github.com/gin-gonic/gin"
)

const (
	RES_CODE = "manage-logbot"
	RES_NAME = "Manage LogBot"
)

func BeforeServerBootstrapp(c core.Rail) error {
	if e := bus.DeclareEventBus(ERROR_LOG_EVENT_BUS); e != nil {
		return e
	}

	bus.SubscribeEventBus(ERROR_LOG_EVENT_BUS, 2,
		func(rail core.Rail, l LogLineEvent) error {
			return SaveErrorLog(rail, l)
		})

	// List error logs endpoint
	server.IPost("/log/error/list", listErrorLogEp, goauth.PathDocExtra(goauth.PathDoc{Desc: "List error logs", Type: goauth.PT_PROTECTED, Code: RES_CODE}))

	// report resources and paths if enabled
	goauth.ReportResourcesOnBootstrapped(c, []goauth.AddResourceReq{
		{Name: RES_NAME, Code: RES_CODE},
	})

	if IsRmErrorLogTaskEnabled() {
		task.ScheduleNamedDistributedTask("0 0/1 * * ?", false, "RemoveErrorLogTask", func(ec core.Rail) error {
			gap := 7 * 24 * time.Hour // seven days ago
			return RemoveErrorLogsBefore(ec, time.Now().Add(-gap))
		})
	}

	return nil
}

func listErrorLogEp(c *gin.Context, ec core.Rail, req ListErrorLogReq) (ListErrorLogResp, error) {
	return ListErrorLogs(ec, req)
}

func AfterServerBootstrapped(rail core.Rail) error {
	logBotConfig := LoadLogBotConfig().Config
	for _, wc := range logBotConfig.WatchConfigs {
		go func(w WatchConfig, nextRail core.Rail) {
			if e := WatchLogFile(nextRail, w, logBotConfig.NodeName); e != nil {
				nextRail.Errorf("WatchLogFile, app: %v, file: %v, %v", w.App, w.File, e)
			}
		}(wc, rail.NextSpan())
	}
	return nil
}
