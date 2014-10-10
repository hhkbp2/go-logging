package logging

import (
    "bytes"
    "fmt"
    "github.com/hhkbp2/go-strftime"
    "strings"
    "unicode/utf8"
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
    switch i {
    case 0:
        return record.Name
    case 1:
        return fmt.Sprintf("%d", record.Level)
    case 2:
        return GetLevelName(record.Level)
    case 3:
        return record.PathName
    case 4:
        return record.FileName
    case 5:
        return fmt.Sprintf("%d", record.LineNo)
    case 6:
        return record.CreatedTime.String()
    case 7:
        return record.AscTime
    case 8:
        return record.Message
    default:
        panic("unsupport format attribute")
    }
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
    format     string
    dateFormat string
}

func NewStandardFormatter(format string, dateFormat string) *StandardFormatter {
    return &StandardFormatter{
        format:     format,
        dateFormat: dateFormat,
    }
}

func (self *StandardFormatter) FormatTime(record *LogRecord) string {
    // Use the library go-strftime to format the time record.Created.
    return strftime.Format(self.dateFormat, record.CreatedTime)
}

func (self *StandardFormatter) Format(record *LogRecord) string {
    record.Message = record.GetMessage()
    if strings.Index(self.format, "%(asctime)s") == -1 {
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
        c, _ := utf8.DecodeRuneInString(format[i:])
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
                if str == attr {
                    buf.WriteString(AttrOfRecord(j, record))
                    i = attrEnd
                    break
                }
            }
        }
    }
    return buf.String()
}
