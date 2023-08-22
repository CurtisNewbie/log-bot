package logbot

import (
	"time"

	"github.com/curtisnewbie/gocommon/bus"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/curtisnewbie/gocommon/task"
	"github.com/gin-gonic/gin"
)

const (
	RES_CODE = "manage-logbot"
	RES_NAME = "Manage LogBot"
)

func BeforeServerBootstrapp(c common.Rail) error {
	if e := bus.DeclareEventBus(ERROR_LOG_EVENT_BUS); e != nil {
		return e
	}

	bus.SubscribeEventBus(ERROR_LOG_EVENT_BUS, 2,
		func(rail common.Rail, l LogLineEvent) error {
			return SaveErrorLog(rail, l)
		})

	// List error logs endpoint
	server.IPost("/log/error/list", listErrorLogEp, goauth.PathDocExtra(goauth.PathDoc{Desc: "List error logs", Type: goauth.PT_PROTECTED, Code: RES_CODE}))

	// report resources and paths if enabled
	if goauth.IsEnabled() {
		server.PreServerBootstrap(func(sc common.Rail) error {
			if e := goauth.AddResource(sc.Ctx, goauth.AddResourceReq{Name: RES_NAME, Code: RES_CODE}); e != nil {
				c.Errorf("Failed to create goauth resource, %v", e)
			}
			return nil
		})
		goauth.ReportPathsOnBootstrapped()
	}

	if IsRmErrorLogTaskEnabled() {
		task.ScheduleNamedDistributedTask("0 0/1 * * ?", false, "RemoveErrorLogTask", func(ec common.Rail) error {
			gap := 7 * 24 * time.Hour // seven days ago
			return RemoveErrorLogsBefore(ec, time.Now().Add(-gap))
		})
	}

	return nil
}

func listErrorLogEp(c *gin.Context, ec common.Rail, req ListErrorLogReq) (ListErrorLogResp, error) {
	return ListErrorLogs(ec, req)
}

func AfterServerBootstrapped(rail common.Rail) error {
	logBotConfig := LoadLogBotConfig().Config
	for _, wc := range logBotConfig.WatchConfigs {
		go func(w WatchConfig, nextRail common.Rail) {
			if e := WatchLogFile(nextRail, w, logBotConfig.NodeName); e != nil {
				nextRail.Errorf("WatchLogFile, app: %v, file: %v, %v", w.App, w.File, e)
			}
		}(wc, rail.NextSpan())
	}
	return nil
}
