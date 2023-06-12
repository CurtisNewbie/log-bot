package logbot

import (
	"fmt"

	"github.com/curtisnewbie/gocommon/common"
)

func AfterServerBootstrapped(c common.ExecContext) error {
	logFiles, e := ListAppLogFiles(c)
	if e != nil {
		return fmt.Errorf("failed to list app log files, %v", e)
	}
	for _, f := range logFiles {
		go func(f AppLogFile) {
			c := common.EmptyExecContext()
			if e := WatchLogFile(c, f); e != nil {
				c.Log.Errorf("failed to register watch task, app: %v, file: %v, %v", f.App, f.File, e)
			}
		}(f)
	}
	return nil
}
