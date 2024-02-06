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

	"github.com/curtisnewbie/miso/miso"
	postbox "github.com/curtisnewbie/postbox/api"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	ErrorLogEventBus = "event.bus.logbot.log.error"
)

var (
	_goLogPat   = regexp.MustCompile(`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[([\w ]+),([\w ]+)\] ([\w\.]+) +: *((?s).*)`)
	_javaLogPat = regexp.MustCompile(`^([0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9:\.]+) +(\w+) +\[[\w \-]+,([\w ]*),([\w ]*),[\w ]*\] [\w\.]+ \-\-\- \[[\w\- ]+\] ([\w\-\.]+) +: *((?s).*)`)

	lineCache   = miso.NewTTLCache[string](time.Second*5, 1000)
	noopElseGet = func() (string, bool) { return "", false }
)

func init() {
	miso.SetDefProp("logbot.node", "default")
}

func lastPos(rail miso.Rail, app string, nodeName string) (int64, error) {
	cmd := miso.GetRedis().Get(fmt.Sprintf("log-bot:pos:%v:%v", nodeName, app))
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

func recPos(rail miso.Rail, app string, nodeName string, pos int64) error {
	posStr := strconv.FormatInt(pos, 10)
	cmd := miso.GetRedis().Set(fmt.Sprintf("log-bot:pos:%v:%v", nodeName, app), posStr, 0)
	return cmd.Err()
}

func WatchLogFile(rail miso.Rail, wc WatchConfig, nodeName string) error {
	rail.Infof("Watching log file '%v' for app '%v'", wc.File, wc.App)
	f, err := os.Open(wc.File)

	if err != nil {
		if !os.IsNotExist(err) { // is possible that the log file doesn't exist
			return fmt.Errorf("failed to open log file, %v", err)
		}
	}

	if f != nil {
		defer f.Close() // the log file is opened
	}

	pos, el := lastPos(rail, wc.App, nodeName)
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
			rail.Infof("Log file '%v' seek to position %v", wc.File, pos)
		}
	}

	// create reader for the file
	var rd *bufio.Reader
	if f != nil {
		rd = bufio.NewReader(f)
	}

	lastRead := time.Now()
	accum := 0 // lines read so far (will be reset when it reaches 1000)
	var prevBytesRead int64
	var prevLine string
	var prevLogLine *LogLine // a single log can contain multiple lines

	for {
		if rd == nil {
			time.Sleep(2 * time.Second) // wait for the file to be created

			f, err = os.Open(wc.File)
			if err != nil {
				f = nil
				continue // the file is still not created
			}
			rail.Infof("Opened %v", wc.File)

			// new file, create reader and set pos = 0
			rd = bufio.NewReader(f)
			pos = 0
		}

		// check if the file is still valid
		if time.Since(lastRead) > 15*time.Second {
			rail.Debug("Checking if the file is still valid, ", wc.File)

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

			lastRead = time.Now()

			if reopenFile {
				f.Close()
				rd = nil
				f = nil
				rail.Infof("Closed file '%v' fd", wc.File)
				continue
			}
		}

		line, err := rd.ReadString('\n')
		if err == nil {

			logLine, e := parseLogLine(rail, line, wc.Type)
			if e == nil {

				// prevLogLine == nil, won't happen unless it is the first log being parsed, or is really in incorrect format
				if prevLogLine != nil {

					// always report the previous log
					if e := reportLine(rail, *prevLogLine, nodeName, wc); e != nil {
						rail.Errorf("Failed to reportLine, logLine: %+v, %v", *prevLogLine, e)
					}

					// move the position only when we report the previous log
					pos += prevBytesRead
					recPos(rail, wc.App, nodeName, pos)
					rail.Debugf("app: %v, pos: %v", wc.App, pos)
				}

				prevBytesRead = int64(len([]byte(line)))
				prevLine = line
				prevLogLine = &logLine
			} else {
				// if current line is not parseable, it's part of previous line
				// we put them together and we parse again
				//
				// yes, we may will just parse it before we do reportLine, but for
				// 90% of the time, the log is single line
				// so it's better leave it here
				prevBytesRead += int64(len([]byte(line)))
				prevLine = prevLine + line
				if parsed, ep := parseLogLine(rail, prevLine, wc.Type); ep == nil {
					prevLogLine = &parsed
				}
			}

			lastRead = time.Now()
			accum += 1

			if accum == 250 {
				time.Sleep(250 * time.Millisecond)
				accum = 0
			}

			continue
		}

		if err == io.EOF {
			accum = 0
			time.Sleep(2 * time.Second)
			continue
		}

		if miso.IsShuttingDown() {
			return nil
		}
	}
}

type LogLineEvent struct {
	App     string
	Node    string
	Time    miso.ETime
	Level   string
	TraceId string
	SpanId  string
	Caller  string
	Message string
}

type LogLine struct {
	Time    miso.ETime
	Level   string
	TraceId string
	SpanId  string
	Caller  string
	Message string
}

func parseLogLine(c miso.Rail, line string, typ string) (LogLine, error) {
	var pat *regexp.Regexp
	if typ == "java" {
		pat = _javaLogPat
	} else {
		pat = _goLogPat
	}
	matches := pat.FindStringSubmatch(line)
	if matches == nil {
		return LogLine{}, fmt.Errorf("doesn't match pattern")
	}

	time, ep := time.ParseInLocation(`2006-01-02 15:04:05.000`, matches[1], time.Local)
	if ep != nil {
		return LogLine{}, fmt.Errorf("time format illegal, %v", ep)
	}

	// only save the first 1000 characters
	msg := matches[6]
	msgRu := []rune(msg)
	if len(msgRu) > 1000 {
		msg = string(msgRu[:1001])
	}

	return LogLine{
		Time:    miso.ETime(time),
		Level:   matches[2],
		TraceId: strings.TrimSpace(matches[3]),
		SpanId:  strings.TrimSpace(matches[4]),
		Caller:  matches[5],
		Message: msg,
	}, nil
}

func reportLine(rail miso.Rail, line LogLine, node string, wc WatchConfig) error {
	if line.Level != "ERROR" {
		return nil
	}

	// prevent reporting same error message way too frequently
	if len(line.Message) < 2048 {
		_, ok := lineCache.Get(line.Message, noopElseGet)
		if ok {
			return nil
		}
		lineCache.Put(line.Message, "1") // race condition is fine
	}

	return miso.PubEventBus(rail,
		LogLineEvent{
			App:     wc.App,
			Node:    node,
			Time:    line.Time,
			Level:   line.Level,
			TraceId: line.TraceId,
			SpanId:  line.SpanId,
			Caller:  line.Caller,
			Message: line.Message,
		},
		ErrorLogEventBus,
	)
}

type SaveErrorLogCmd struct {
	Node    string
	App     string
	Caller  string
	TraceId string
	SpanId  string
	ErrMsg  string
	RTime   miso.ETime `gorm:"column:rtime"`
}

func SaveErrorLog(rail miso.Rail, evt LogLineEvent) error {
	el := SaveErrorLogCmd{
		Node:    evt.Node,
		App:     evt.App,
		Caller:  evt.Caller,
		TraceId: evt.TraceId,
		SpanId:  evt.SpanId,
		ErrMsg:  evt.Message,
		RTime:   evt.Time,
	}
	err := miso.GetMySQL().
		Table("error_log").
		Create(&el).
		Error

	if err == nil {
		adminUserNos := miso.GetPropStrSlice("admin.user-no")
		if len(adminUserNos) > 0 {
			if cerr := postbox.CreateNotification(rail, postbox.CreateNotificationReq{
				Title:           fmt.Sprintf("Logbot - %s has error", evt.App),
				Message:         fmt.Sprintf("%s [%s,%s] %s : %s", evt.Time.FormatClassic(), evt.TraceId, evt.SpanId, evt.Caller, evt.Message),
				ReceiverUserNos: adminUserNos,
			}); cerr != nil {
				rail.Errorf("failed to create platform notification, %v", cerr)
			}
		}
	}
	return err
}

type ListedErrorLog struct {
	Id      int64      `json:"id"`
	Node    string     `json:"node"`
	App     string     `json:"app"`
	Caller  string     `json:"caller"`
	TraceId string     `json:"traceId"`
	SpanId  string     `json:"spanId"`
	ErrMsg  string     `json:"errMsg"`
	RTime   miso.ETime `json:"rtime" gorm:"column:rtime"`
}

type ListErrorLogReq struct {
	App  string      `json:"app"`
	Page miso.Paging `json:"page"`
}

type ListErrorLogResp struct {
	Page    miso.Paging      `json:"page"`
	Payload []ListedErrorLog `json:"payload"`
}

func newListErrorLogsQry(rail miso.Rail, r ListErrorLogReq) *gorm.DB {
	t := miso.GetMySQL().
		Table("error_log")

	if r.App != "" {
		t = t.Where("app = ?", r.App)
	}

	return t
}

func ListErrorLogs(rail miso.Rail, r ListErrorLogReq) (ListErrorLogResp, error) {
	var listed []ListedErrorLog
	e := newListErrorLogsQry(rail, r).
		Offset(r.Page.GetOffset()).
		Limit(r.Page.GetLimit()).
		Order("rtime desc").
		Scan(&listed).Error

	if e != nil {
		return ListErrorLogResp{}, e
	}

	var total int
	e = newListErrorLogsQry(rail, r).
		Select("count(*)").
		Scan(&total).Error
	if e != nil {
		return ListErrorLogResp{}, e
	}

	return ListErrorLogResp{Page: r.Page.ToRespPage(total), Payload: listed}, nil
}

func RemoveErrorLogsBefore(rail miso.Rail, upperBound time.Time) error {
	rail.Infof("Remove error logs before %s", upperBound)
	return miso.GetMySQL().Exec("delete from error_log where rtime < ?", upperBound).Error
}
