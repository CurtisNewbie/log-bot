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

func BeforeServerBootstrapp(c common.ExecContext) error {
	if e := bus.DeclareEventBus(ERROR_LOG_EVENT_BUS); e != nil {
		return e
	}

	if e := bus.SubscribeEventBus(ERROR_LOG_EVENT_BUS, 2, func(l LogLineEvent) error {
		ec := common.EmptyExecContext()
		return SaveErrorLog(ec, l)
	}); e != nil {
		return e
	}

	// List error logs endpoint
	server.IPost("/log/error/list", listErrorLogEp, goauth.PathDocExtra(goauth.PathDoc{Desc: "List error logs", Type: goauth.PT_PROTECTED, Code: RES_CODE}))

	// report resources and paths if enabled
	if goauth.IsEnabled() {
		server.OnServerBootstrapped(func(sc common.ExecContext) error {
			if e := goauth.AddResource(sc.Ctx, goauth.AddResourceReq{Name: RES_NAME, Code: RES_CODE}); e != nil {
				c.Log.Errorf("Failed to create goauth resource, %v", e)
			}
			return nil
		})
		goauth.ReportPathsOnBootstrapped()
	}

	if IsRmErrorLogTaskEnabled() {
		task.ScheduleNamedDistributedTask("0 0 0/1 * * ?", "RemoveErrorLogTask", func(ec common.ExecContext) error {
			gap := 7 * 24 * time.Hour // seven days ago
			return RemoveErrorLogsBefore(ec, time.Now().Add(-gap))
		})
	}

	return nil
}

func listErrorLogEp(c *gin.Context, ec common.ExecContext, req ListErrorLogReq) (ListErrorLogResp, error) {
	return ListErrorLogs(ec, req)
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
