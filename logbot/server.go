package logbot

import (
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/goauth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/gin-gonic/gin"
)

const (
	RES_CODE = "manage-logbot"
	RES_NAME = "Manage LogBot"
)

func BeforeServerBootstrapp(rail miso.Rail) error {
	common.LoadBuiltinPropagationKeys()

	miso.SubEventBus(ERROR_LOG_EVENT_BUS, 2,
		func(rail miso.Rail, l LogLineEvent) error {
			return SaveErrorLog(rail, l)
		})

	// List error logs endpoint
	miso.IPost("/log/error/list",
		func(c *gin.Context, ec miso.Rail, req ListErrorLogReq) (any, error) {
			return ListErrorLogs(ec, req)
		}).
		Extra(goauth.PathDocExtra(goauth.PathDoc{Desc: "List error logs", Type: goauth.PT_PROTECTED, Code: RES_CODE})).
		Build()

	// report resources and paths if enabled
	goauth.ReportResourcesOnBootstrapped(rail, []goauth.AddResourceReq{
		{Name: RES_NAME, Code: RES_CODE},
	})

	if IsRmErrorLogTaskEnabled() {
		miso.ScheduleDistributedTask(miso.Job{
			Cron:            "0 0/1 * * ?",
			CronWithSeconds: false,
			Name:            "RemoveErrorLogTask",
			Run: func(ec miso.Rail) error {
				gap := 7 * 24 * time.Hour // seven days ago
				return RemoveErrorLogsBefore(ec, time.Now().Add(-gap))
			}})
	}

	return nil
}

func AfterServerBootstrapped(rail miso.Rail) error {
	logBotConfig := LoadLogBotConfig().Config
	for _, wc := range logBotConfig.WatchConfigs {
		go func(w WatchConfig, nextRail miso.Rail) {
			if e := WatchLogFile(nextRail, w, logBotConfig.NodeName); e != nil {
				nextRail.Errorf("WatchLogFile, app: %v, file: %v, %v", w.App, w.File, e)
			}
		}(wc, rail.NextSpan())
	}
	return nil
}
