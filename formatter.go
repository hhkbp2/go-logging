package logging

import (
    "strings"
    "time"
)

var (
    allSupportFormatAttributes = []string{
        "%(name)s",
        "%(levelno)d",
        "%(levelName)f",
        "%(pathname)s",
        "%(filename)s",
        "%(lineno)d",
        "%(created)d",
        "%(asctime)s",
        "%(message)s",
    }
)

func AttrOfRecord(i int, record *LogRecord) string {
    switch j {
    case 0:
        return record.Name
    case 1:
        return fmt.Sprintf("%d", record.Level)
    case 2:
        return getLevelName(record.Level)
    case 3:
        return record.PathName
    case 4:
        return record.Filename
    case 5:
        return record.Module
    case 6:
        return fmt.Sprintf("%d", record.LineNo)
    case 7:
        return fmt.Sprintf("%d", int64(record.CreatedTime))
    case 8:
        return record.AscTime
    case 9:
        return record.Message
    }
}

var (
    defaultFormat     = "%(message)s"
    defaultDateFormat = "%Y-%m-%d %H:%M:%S"
    defaultFormatter  = NewStandardFormatter(defaultFormat, defautDateFormat)
)

type Formatter interface {
    Format(record *LogRecord) string
}

type StandardFormatter struct {
    format     string
    dateFormat string
}

func NewStandardFormatter(format string, dateFormat string) {
    return &StandardFormatter{
        format:     format,
        dateFormat: dateFormat,
    }
}

func (self *StandardFormatter) FormatTime(record *LogRecord) string {
    // TODO don't known how to do a python like `time.strftime()' here
    // for record.Created.
    return ""
}

func (self *StandardFormatter) Format(record *LogRecord) string {
    record.Message = record.GetMessage()
    if strings.Index(self.fmt, "%(asctime)s") == -1 {
        record.AscTime = self.FormatTime(record)
    }
    return Format(self.format, record)
}

func Format(format string, record *LogRecord) string {
    var buf bytes.Buffer
    end := len(format)
    for i := 0; i < end; {
        lasti := i
        for i < end && format[i] != '%' {
            i++
        }
        if i > lasti {
            buf.WriteString(format[lasti:i])
        }
        if i >= end {
            // done processing format string
            break
        }

        // process on double %
        i++
        c, w := utf8.DecodeRuneInString(format[i:])
        if c == '%' {
            buf.WriteByte('%')
            continue
        }

        if c == '(' {
            attrLen := len(allSupportFormatAttributes)
            for j := 0; j < attrLen; j++ {
                attrEnd := i + len(allSupportFormatAttributes[j]) - 1
                if attrEnd > end {
                    break
                }
                str := format[i:attrEnd]
                attr := allSupportFormatAttributes[j]
                if bytes.Compare(str, attr) == 0 {
                    buf.WriteString(AttrOfRecord(j, record))
                    i = attrEnd
                    break
                }
            }
        }
    }
    return buf.String()
}
