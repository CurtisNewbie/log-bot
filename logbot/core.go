package logbot

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/mysql"
	red "github.com/curtisnewbie/gocommon/redis"
	"github.com/curtisnewbie/gocommon/server"
	"github.com/go-redis/redis"
)

var (
	_logLinePat = regexp.MustCompile(`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[([\w ]+),([\w ]+)\] ([\w\.]+) +: +(.*)`)
)

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

func lastPos(c common.ExecContext, id int64) (int64, error) {
	cmd := red.GetRedis().Get(fmt.Sprintf("log-bot:pos:id:%v", id))
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return 0, nil
		}
		return 0, cmd.Err()
	}

	n, ea := strconv.Atoi(cmd.Val())
	if ea != nil {
		return 0, nil
	}
	if n < 0 {
		n = 0
	}
	return int64(n), nil
}

func recPos(c common.ExecContext, id int64, pos int64) error {
	posStr := strconv.FormatInt(pos, 10)
	cmd := red.GetRedis().Set(fmt.Sprintf("log-bot:pos:id:%v", id), posStr, 0)
	return cmd.Err()
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

	pos, el := lastPos(c, appLog.Id)
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
			c.Log.Infof("Log file '%v' seek to position %v", appLog.File, pos)
		}
	}

	// create reader for the file
	var rd *bufio.Reader
	if f != nil {
		rd = bufio.NewReader(f)
	}

	lastRead := time.Now()
	accum := 0 // lines read so far (will be reset when it reaches 1000)

	// TODO: should be querying ElasticSearch for distributed environment, this should work for single node for now
	for {
		if rd == nil {
			time.Sleep(2 * time.Second) // wait for the file to be created

			f, err = os.Open(appLog.File)
			if err != nil {
				f = nil
				continue // the file is still not created
			}
			c.Log.Infof("Opened %v", appLog.File)

			// new file, create reader and set pos = 0
			rd = bufio.NewReader(f)
			pos = 0
		}

		// check if the file is still valid
		if time.Since(lastRead) > 1*time.Minute {

			reopenFile := false

			fi, es := f.Stat()
			if es != nil {
				// if the file is deleted, es will still be nil
				reopenFile = true
			}

			if !reopenFile {
				// https://stackoverflow.com/questions/53184549/how-to-detect-deleted-file
				nlink := uint64(0)
				if sys := fi.Sys(); sys != nil {
					if stat, ok := sys.(*syscall.Stat_t); ok {
						nlink = uint64(stat.Nlink)
					}
				}
				if nlink < 1 { // no hard links, the underlying file is deleted already
					reopenFile = true
				}
			}

			if reopenFile {
				f.Close()
				rd = nil
				f = nil
				lastRead = time.Now()
				c.Log.Infof("Closed file '%v' fd", appLog.File)
				continue
			}
		}

		line, err := rd.ReadString('\n')
		if err == nil {
			lineLen := int64(len([]byte(line)))
			logLine, e := parseLine(c, line, appLog, pos+lineLen)
			if e == nil {
				if e := reportLine(c, logLine); e != nil {
					c.Log.Errorf("Failed to reportLine, logLine: %+v, %v", logLine, e)
				}
			} else {
				c.Log.Errorf("Failed to parse logLine, %v, line: '%v'", e, line)
			}

			pos += lineLen // move the position
			lastRead = time.Now()
			accum += 1

			if accum == 1000 {
				recPos(c, appLog.Id, pos) // record position every 1000 lines
				time.Sleep(500 * time.Millisecond)
				accum = 0
			}

			continue
		}

		// the file may be truncated or renamed
		if err == io.EOF {
			recPos(c, appLog.Id, pos) // record position every 1000 lines
			accum = 0
			time.Sleep(2 * time.Second)
			continue
		}

		if server.IsShuttingDown() {
			recPos(c, appLog.Id, pos)
			return nil
		}
	}
}

type logLine struct {
	Time    time.Time
	Level   string
	TraceId string
	SpanId  string
	Func    string
	Message string
}

func parseLogLine(c common.ExecContext, line string) (logLine, error) {
	matches := _logLinePat.FindStringSubmatch(line)
	if matches == nil {
		return logLine{}, fmt.Errorf("doesn't match pattern")
	}

	time, ep := time.Parse(`2006-01-02 15:04:05.000`, matches[1])
	if ep != nil {
		return logLine{}, fmt.Errorf("time format illegal, %v", ep)
	}
	return logLine{
		Time:    time,
		Level:   matches[2],
		TraceId: strings.TrimSpace(matches[3]),
		SpanId:  strings.TrimSpace(matches[4]),
		Func:    matches[5],
		Message: matches[6],
	}, nil
}

func parseLine(c common.ExecContext, line string, appLog AppLogFile, pos int64) (logLine, error) {
	logLine, err := parseLogLine(c, line)
	c.Log.Infof("id: %v, app: %v, pos: %v", appLog.App, appLog.Id, pos)
	return logLine, err
}

func reportLine(c common.ExecContext, line logLine) error {
	if line.Level != "ERROR" {
		return nil
	}

	c.Log.Errorf("Found %+v", line)
	// TODO
	return nil
}
