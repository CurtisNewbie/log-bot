package logbot

import (
	"testing"

	"github.com/curtisnewbie/gocommon/common"
)

func TestParseLine(t *testing.T) {
	line := `2023-06-13 12:58:35.509 INFO  [                ,                ] consul.DeregisterService      : Deregistering current instance on Consul, service_id: 'goauth-8081'`
	logLine, err := parseLogLine(common.EmptyExecContext(), line)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-13 22:16:13.746 ERROR [v2geq7340pbfxcc9,k1gsschfgarpc7no] main.registerWebEndpoints.func2 : Oh on!`
	logLine, err = parseLogLine(common.EmptyExecContext(), line)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-14 09:50:30.500 DEBUG [ptqnta70npjjxfz8,114lkur90ui6ywqt] redis.TimedRLockRun.func1     : Released lock for key 'rcache:POST:/goauth/open/api/path/update'
`
	logLine, err = parseLogLine(common.EmptyExecContext(), line)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)

	line = `2023-06-16 16:29:48.811 INFO  [3pdac13hagg9v8cs,pskalvoqmaets17f] server.BootstrapServer        :



---------------------------------------------- goauth started (took: 59ms) --------------------------------------------

`

	logLine, err = parseLogLine(common.EmptyExecContext(), line)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", logLine)
}
