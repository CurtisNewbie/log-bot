package logbot

import (
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

type ErrorLog struct {
	Node     string
	App      string
	Time     util.ETime
	TraceId  string
	SpanId   string
	FuncName string
	Message  string
}

var (
	ReportLogPipeline = rabbit.NewEventPipeline[ErrorLog]("logbot:error-log:report:pipeline").
		LogPayload()
)

func InitPipeline(rail miso.Rail) {
	ReportLogPipeline.Listen(2, func(rail miso.Rail, t ErrorLog) error {
		SaveErrorLog(rail, LogLineEvent{
			App:     t.App,
			Node:    t.Node,
			Time:    t.Time,
			Level:   "ERROR",
			TraceId: t.TraceId,
			SpanId:  t.SpanId,
			Caller:  t.FuncName,
			Message: t.Message,
		})
		return nil
	})
}
