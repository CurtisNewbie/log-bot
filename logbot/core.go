package logbot

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/server"
)

type LogPos struct {
	File     string
	Position int64
}

type AppLogFile struct {
	Id   int64
	App  string
	File string
}

func ListAppLogFiles(c common.ExecContext) ([]AppLogFile, error) {
	var files []AppLogFile
	e := mysql.GetConn().
		Table("app_log_file").
		Scan(&files).Error
	return files, e
}

func LastPos(c common.ExecContext, id int64) (int64, error) {
	return 0, nil
}

func recPos(c common.ExecContext, id int64, pos int64) error {
	// c.Log.Infof("ID: %s - Pos: %d", id, pos)
	return nil
}

func WatchLogFile(c common.ExecContext, appLog AppLogFile) error {
	f, err := os.Open(appLog.File)

	if err != nil {
		if !os.IsNotExist(err) { // is possible that the log file doesn't exist
			return fmt.Errorf("failed to open log file, %v", err)
		}
	}

	if f != nil {
		defer f.Close() // the log file is opened
	}

	pos, el := LastPos(c, appLog.Id)
	if el != nil {
		return fmt.Errorf("failed to find last pos, %v", el)
	}

	if f != nil && pos > 0 {
		fi, es := f.Stat()
		if es != nil {
			return es
		}

		// the file was truncated
		if pos > fi.Size() {
			pos = 0
		}

		// seek pos
		if pos > 0 {
			_, e := f.Seek(pos, io.SeekStart)
			if e != nil {
				return fmt.Errorf("failed to seek pos, %v", e)
			}
		}
	}

	// create reader for the file
	var rd *bufio.Reader
	if f != nil {
		rd = bufio.NewReader(f)
	}

	lastRead := time.Now()
	accum := 0

	// TODO: should be querying ElasticSearch for distributed environment, this should work for single node for now
	for {
		if rd == nil {
			time.Sleep(2 * time.Second) // wait for the file to be created

			f, err := os.Open(appLog.File)
			if err != nil {
				continue // the file is still not created
			}

			// new file, create reader and set pos = 0
			rd = bufio.NewReader(f)
			pos = 0
		}

		// check if the file is still valid
		if time.Since(lastRead) > 30*time.Second {
			_, es := f.Stat()
			if es != nil {
				f.Close()
				rd = nil
				f = nil
				// TODO: this doesn't seem to work
			}
		}

		line, err := rd.ReadString('\n')
		if err == nil {
			parseLine(c, line, appLog)
			pos += int64(len([]byte(line)))
			lastRead = time.Now()
			accum += 1

			if accum == 1000 {
				recPos(c, appLog.Id, pos)
				time.Sleep(500 * time.Millisecond)
				accum = 0
			}

			continue
		}

		// the file may be truncated or renamed
		if err == io.EOF {
			accum = 0
			time.Sleep(2 * time.Second)
			continue
		}

		if server.IsShuttingDown() {
			return nil
		}
	}
}

func parseLine(c common.ExecContext, line string, appLog AppLogFile) {
	c.Log.Infof("%s - %s\n", appLog.File, line)
}
