package logging

import (
	"errors"
	"github.com/hhkbp2/go-strftime"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	Day  = 24 * time.Hour
	Week = 7 * Day
)

var (
	ErrorInvalidFormat = errors.New("invalid format")
)

// Handler for logging to a file, rotating the log file at certain timed
// intervals.
//
// if backupCount is > 0, when rollover is done, no more than backupCount
// files are kept - the oldest ones are deleted.
type TimedRotatingFileHandler struct {
	*BaseRotatingHandler
	when         string
	weekday      int
	interval     time.Duration
	rolloverTime time.Time
	backupCount  uint32
	suffix       string
	extMatch     string
	utc          bool
}

// Note: weekday index starts from 0(Monday) to 6(Sunday) in Python.
// But in Golang weekday index starts from 0(Sunday) to 6(Saturday).
// Here we stick to semantics of the original Python logging interface.
func NewTimedRotatingFileHandler(
	filepath string,
	mode int,
	bufferSize int,
	when string,
	interval uint32,
	backupCount uint32,
	utc bool) (*TimedRotatingFileHandler, error) {

	var timeInterval time.Duration
	var suffix, extMatch string
	var weekday int
	// Calculate the real rollover interval, which is just the number seconds
	// between rollovers. Also set the filename suffix used when a rollover
	// occurs. Current 'when' events supported:
	// S - Seconds
	// M - Minutes
	// H - Hours
	// D - Days
	// midnight - roll over at midnight
	// W{0-6} - roll over on a certain weekday; 0 - Monday
	// Case of the 'when' specifier is not important; lower or upper case
	// will work.
	when = strings.ToUpper(when)
	switch {
	case when == "S":
		timeInterval = time.Second
		suffix = "%Y-%m-%d_%H-%M-%S"
		extMatch = `^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}$`
	case when == "M":
		timeInterval = time.Minute
		suffix = "%Y-%m-%d_%H-%M"
		extMatch = `^\d{4}-\d{2}-\d{2}_\d{2}-\d{2}$`
	case when == "H":
		timeInterval = time.Hour
		suffix = "%Y-%m-%d_%H"
		extMatch = `^\d{4}-\d{2}-\d{2}_\d{2}$`
	case (when == "D") || (when == "MIDNIGHT"):
		timeInterval = Day
		suffix = "%Y-%m-%d"
		extMatch = `^\d{4}-\d{2}-\d{2}$`
	case strings.HasPrefix(when, "W"):
		timeInterval = Week
		if len(when) != 2 {
			return nil, ErrorInvalidFormat
		}
		dayChar := when[1]
		if (dayChar < '0') || (dayChar > '6') {
			return nil, ErrorInvalidFormat
		}
		// cast Python style index value to Golang style index value
		weekday = (int(dayChar-'0') + 1) % 7
		suffix = "%Y-%m-%d"
		extMatch = `^\d{4}-\d{2}-\d{2}$`
	default:
		return nil, ErrorInvalidFormat
	}
	timeInterval = time.Duration(int64(timeInterval) * int64(interval))
	baseHandler, err := NewBaseRotatingHandler(filepath, mode, bufferSize)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(baseHandler.GetFilePath())
	if err != nil {
		baseHandler.Close()
		return nil, err
	}
	object := &TimedRotatingFileHandler{
		BaseRotatingHandler: baseHandler,
		when:                when,
		weekday:             weekday,
		interval:            timeInterval,
		backupCount:         backupCount,
		suffix:              suffix,
		extMatch:            extMatch,
		utc:                 utc,
	}
	object.rolloverTime = object.computeRolloverTime(fileInfo.ModTime())
	return object, nil
}

// Work out the rollover time based on the specified time.
func (self *TimedRotatingFileHandler) computeRolloverTime(
	currentTime time.Time) time.Time {

	result := currentTime.Add(self.interval)
	// If we are rolling over at midnight or weekly, then the interval is
	// already known.  What we need to figure out is WHEN the next interval is.
	// In other words, if you are rolling over at midnight, then
	// your base interval is 1 day, but you want to start that one day clock
	// at midnight, not now.
	// So, we have to fudge the rolloverTime value in order trigger the first
	// rollover at the right time.  After that, the regular interval will
	// take care of the rest.
	// Note that this code doesn't care about leap seconds.
	if (self.when == "MIDNIGHT") || strings.HasPrefix(self.when, "W") {
		var t time.Time
		if self.utc {
			t = currentTime.UTC()
		} else {
			t = currentTime.Local()
		}
		dayStartTime := time.Date(
			t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		result = currentTime.Add(Day - t.Sub(dayStartTime))
		// If we are rolling over on a certain day, add in the number of days
		// until the next rollover, but offset by 1 since we just calculated
		// the time until the next day starts.  There are three cases:
		// Case 1) The day to rollover is today; in this case, do nothing
		// Case 2) The day to rollover is further in the interval (i.e.,
		//         today is day 3 (Wednesday) and rollover is on day 6
		//         (Saturday). Days to next rollover is simply 6 - 3, or 3)
		// Case 3) The day to rollover is behind us in the interval (i.e.,
		//         today is day 5 (Friday) and rollover is on day 4 (Thursday).
		//         Days to rollover is 6 - 5 + 4 + 1, or 6.)  In this case,
		//         it's the number of days left in the current week (1) plus
		//         the number of days in the next week until the
		//         rollover day (5).
		// THe calculations described in 2) and 3) above need to
		// have a day added.  This is because the above time calculation
		// takes us to midnight on this day, i.e., the start of the next day.
		if strings.HasPrefix(self.when, "W") {
			weekday := int(t.Weekday())
			if weekday != self.weekday {
				var daysToWait int
				if weekday < self.weekday {
					daysToWait = self.weekday - weekday
				} else {
					daysToWait = 6 - weekday + self.weekday + 1
				}
				result = result.Add(time.Duration(int64(daysToWait) * int64(Day)))
				// NOTE: we skip the daylight savings time issues here
				// because time library in Golang doesn't support it.
			}
		}
	}
	return result
}

// Determine if rollover should occur.
func (self *TimedRotatingFileHandler) ShouldRollover(
	record *LogRecord) (bool, string) {

	overTime := time.Now().After(self.rolloverTime)
	return overTime, self.Format(record)
}

// Determine the files to delete when rolling over.
func (self *TimedRotatingFileHandler) getFilesToDelete() ([]string, error) {
	dirName, baseName := filepath.Split(self.GetFilePath())
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	prefix := baseName + "."
	pattern, err := regexp.Compile(self.extMatch)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, info := range fileInfos {
		fileName := info.Name()
		if strings.HasPrefix(fileName, prefix) {
			suffix := fileName[len(prefix):]
			if pattern.MatchString(suffix) {
				fileNames = append(fileNames, fileName)
			}
		}
	}
	// no need to sort fileNames since ioutil.ReadDir() returns sorted list.
	var result []string
	if uint32(len(fileNames)) < self.backupCount {
		return result, nil
	}
	result = fileNames[:uint32(len(fileNames))-self.backupCount]
	for i := 0; i < len(result); i++ {
		result[i] = filepath.Join(dirName, result[i])
	}
	return result, nil
}

// Do a rollover; in this case, a date/time stamp is appended to the filename
// when the rollover happens.  However, you want the file to be named for
// the start of the interval, not the current time.  If there is a backup
// count, then we have to get a list of matching filenames, sort them and
// remove the one with the oldest suffix.
func (self *TimedRotatingFileHandler) DoRollover() (err error) {
	self.Close()
	defer func() {
		if e := self.Open(); e != nil {
			if e != nil {
				err = e
			}
		}
	}()
	currentTime := time.Now()
	t := self.rolloverTime.Add(time.Duration(-int64(self.interval)))
	if self.utc {
		t = t.UTC()
	} else {
		t = t.Local()
	}
	baseFilename := self.GetFilePath()
	dfn := baseFilename + "." + strftime.Format(self.suffix, t)
	if FileExists(dfn) {
		if err := os.Remove(dfn); err != nil {
			return err
		}
	}
	if err := os.Rename(baseFilename, dfn); err != nil {
		return err
	}
	if self.backupCount > 0 {
		files, err := self.getFilesToDelete()
		if err != nil {
			return err
		}
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				return err
			}
		}
	}
	self.rolloverTime = self.computeRolloverTime(currentTime)
	return nil
}

// Emit a record.
func (self *TimedRotatingFileHandler) Emit(record *LogRecord) error {
	return self.RolloverEmit(self, record)
}

func (self *TimedRotatingFileHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}
