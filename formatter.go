package logging

import (
    "bytes"
    "fmt"
    "github.com/hhkbp2/go-strftime"
    "regexp"
    "strings"
)

type ExtractAttr func(record *LogRecord) string

var (
    attrToFunc = map[string]ExtractAttr{
        "%(name)s": func(record *LogRecord) string {
            return record.Name
        },
        "%(levelno)d": func(record *LogRecord) string {
            return fmt.Sprintf("%d", record.Level)
        },
        "%(levelname)s": func(record *LogRecord) string {
            return GetLevelName(record.Level)
        },
        "%(pathname)s": func(record *LogRecord) string {
            return record.PathName
        },
        "%(filename)s": func(record *LogRecord) string {
            return record.FileName
        },
        "%(lineno)d": func(record *LogRecord) string {
            return fmt.Sprintf("%d", record.LineNo)
        },
        "%(funcname)s": func(record *LogRecord) string {
            return record.FuncName
        },
        "%(created)d": func(record *LogRecord) string {
            return fmt.Sprintf("%d", record.CreatedTime.UnixNano())
        },
        "%(asctime)s": func(record *LogRecord) string {
            return record.AscTime
        },
        "%(message)s": func(record *LogRecord) string {
            return record.Message
        },
    }
    formatRe *regexp.Regexp
)

func init() {
    var buf bytes.Buffer
    buf.WriteString("(%(?:%")
    for attr, _ := range attrToFunc {
        buf.WriteString("|")
        buf.WriteString(regexp.QuoteMeta(attr[1:]))
    }
    buf.WriteString("))")
    re := buf.String()
    formatRe = regexp.MustCompile(re)
}

var (
    defaultFormat     = "%(message)s"
    defaultDateFormat = "%Y-%m-%d %H:%M:%S %3n"
    defaultFormatter  = NewStandardFormatter(defaultFormat, defaultDateFormat)
)

type Formatter interface {
    Format(record *LogRecord) string
}

type StandardFormatter struct {
    format       string
    toFormatTime bool
    dateFormat   string
}

func NewStandardFormatter(format string, dateFormat string) *StandardFormatter {
    toFormatTime := false
    if strings.Index(format, "%(asctime)s") > 0 {
        toFormatTime = true
    }
    return &StandardFormatter{
        format:       format,
        toFormatTime: toFormatTime,
        dateFormat:   dateFormat,
    }
}

func (self *StandardFormatter) FormatTime(record *LogRecord) string {
    // Use the library go-strftime to format the time record.Created.
    return strftime.Format(self.dateFormat, record.CreatedTime)
}

func (self *StandardFormatter) Format(record *LogRecord) string {
    record.Message = record.GetMessage()
    if self.toFormatTime {
        record.AscTime = self.FormatTime(record)
    }
    return Format(self.format, record)
}

func repl(match string, record *LogRecord) string {
    if match == "%%" {
        return "%"
    }

    extractFunc, ok := attrToFunc[match]
    if ok {
        return extractFunc(record)
    }
    return match
}

func Format(format string, record *LogRecord) string {
    fn := func(match string) string {
        return repl(match, record)
    }
    return formatRe.ReplaceAllStringFunc(format, fn)
}
