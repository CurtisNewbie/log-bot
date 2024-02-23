package logbot

import (
	"time"

	"github.com/curtisnewbie/gocommon/auth"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	ResourceManageLogbot = "manage-logbot"
)

func BeforeServerBootstrapp(rail miso.Rail) error {
	common.LoadBuiltinPropagationKeys()
	auth.ExposeResourceInfo([]auth.Resource{
		{Name: "Manage LogBot", Code: ResourceManageLogbot},
	})

	miso.SubEventBus(ErrorLogEventBus, 2,
		func(rail miso.Rail, l LogLineEvent) error {
			return SaveErrorLog(rail, l)
		})

	miso.IPost("/log/error/list",
		func(inb *miso.Inbound, req ListErrorLogReq) (any, error) {
			return ListErrorLogs(inb.Rail(), req)
		}).
		Desc("List error logs").
		Resource(ResourceManageLogbot)

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
