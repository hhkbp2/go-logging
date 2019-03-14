package logging

import (
	"bytes"
	"fmt"
	"github.com/hhkbp2/go-strftime"
	"regexp"
)

// Function type of extracting the corresponding LogRecord info for
// the attribute string.
type ExtractAttr func(record *LogRecord) string

// All predefined attribute string and their ExtractAttr functions.
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
	formatRe = initFormatRegexp()

	// Default format strings.
	defaultFormat     = "%(message)s"
	defaultDateFormat = "%Y-%m-%d %H:%M:%S %3n"
	defaultFormatter  = NewStandardFormatter(defaultFormat, defaultDateFormat)
)

// Initialize global regexp for attribute matching.
func initFormatRegexp() *regexp.Regexp {
	var buf bytes.Buffer
	buf.WriteString("(%(?:%")
	for attr, _ := range attrToFunc {
		buf.WriteString("|")
		buf.WriteString(regexp.QuoteMeta(attr[1:]))
	}
	buf.WriteString("))")
	re := buf.String()
	return regexp.MustCompile(re)
}

type GetFormatArgsFunc func(record *LogRecord) []interface{}

// Formatter interface is for converting a LogRecord to text.
// Formatters need to know how a LogRecord is constructed. They are responsible
// for converting a LogRecord to (usually) a string which can be interpreted
// by either a human or an external system.
type Formatter interface {
	Format(record *LogRecord) string
}

// The standard formatter. It allows a formatting string to be specified.
// If none is supplied, the default value of "%(message)s" is used.
//
// The formatter can be initialized with a format string which makes use of
// knowledge of the LogRecord attributes - e.g. the default value mentioned
// above makes use of the fact that the user's message and arguments are
// preformatted into a LogRecord's message attribute. Currently, the usefull
// attributes in a LogRecord are described by:
//
// %(name)s            Name of the logger(logging channel)
// %(levelno)d         Numeric logging level for the message
// %(levelname)s       Text logging level for the message
// %(pathname)s        Full pathname of the source file where the logging
//                     call was issued (is available)
// %(filename)s        Filename portion of pathname
// %(lineno)d          Source line number where the logging call was issued
// %(funcname)s        Function name
// %(created)d         Time when the LogRecord was created(time.Now()
//                     return value)
// %(asctime)s         Textual time when LogRecord was created
// %(message)s         The result of record.GetMessage(), computed just as the
//                     record is emitted
type StandardFormatter struct {
	format            string
	strFormat         string
	getFormatArgsFunc GetFormatArgsFunc
	toFormatTime      bool
	dateFormat        string
	dateFormatter     *strftime.Formatter
}

// Initialize the formatter with specified format strings.
// Allow for specialized date formatting with the dateFormat arguement.
func NewStandardFormatter(format string, dateFormat string) *StandardFormatter {
	toFormatTime := false
	size := 0
	f1 := func(match string) string {
		if match == "%%" {
			return "%"
		}
		if match == "%(asctime)s" {
			toFormatTime = true
		}
		size++
		return "%s"
	}
	strFormat := formatRe.ReplaceAllStringFunc(format, f1)
	funs := make([]ExtractAttr, 0, size)
	f2 := func(match string) string {
		extractFunc, ok := attrToFunc[match]
		if ok {
			funs = append(funs, extractFunc)
		}
		return match
	}
	formatRe.ReplaceAllStringFunc(format, f2)
	getFormatArgsFunc := func(record *LogRecord) []interface{} {
		result := make([]interface{}, 0, len(funs))
		for _, f := range funs {
			result = append(result, f(record))
		}
		return result
	}
	var dateFormatter *strftime.Formatter
	if toFormatTime {
		dateFormatter = strftime.NewFormatter(dateFormat)
	}
	return &StandardFormatter{
		format:            format,
		strFormat:         strFormat + "\n",
		getFormatArgsFunc: getFormatArgsFunc,
		toFormatTime:      toFormatTime,
		dateFormat:        dateFormat,
		dateFormatter:     dateFormatter,
	}
}

// FormatTime returns the creation time of the specified LogRecord as formatted text.
// This method should be called from Format() by a formatter which wants to
// make use of a formatted time. This method can be overridden in formatters
// to provide for any specific requirement, but the basic behaviour is as
// follows: the dateFormat is used with strftime.Format() to format
// the creation time of the record.
func (self *StandardFormatter) FormatTime(record *LogRecord) string {
	// Use the library go-strftime to format the time record.Created.
	return self.dateFormatter.Format(record.CreatedTime)
}

// Format the specified record as text.
// The record's attribute is used as the operand to a string formatting
// operation which yields the returned string. Before the formatting,
// a couple of preparatory steps are carried out. The message attribute of
// the record is computed using LogRecord.GetMessage(). If the formatting
// string uses the time, FormatTime() is called to format the event time.
func (self *StandardFormatter) Format(record *LogRecord) string {
	record.GetMessage()
	if self.toFormatTime {
		record.AscTime = self.FormatTime(record)
	}
	return self.FormatAll(record)
}

// Helper function using regexp to replace every valid format attribute string
// to the record's specific value.
func (self *StandardFormatter) FormatAll(record *LogRecord) string {
	return fmt.Sprintf(self.strFormat, self.getFormatArgsFunc(record)...)
}

// A formatter suitable for formatting a number of records.
type BufferingFormatter struct {
	lineFormatter Formatter
}

// Initialize the buffering formatter with specified line formatter.
func NewBufferingFormatter(lineFormatter Formatter) *BufferingFormatter {
	return &BufferingFormatter{
		lineFormatter: lineFormatter,
	}
}

// FormatHeader returns the header string for the specified records.
func (self *BufferingFormatter) FormatHeader(_ []*LogRecord) string {
	return ""
}

// FormatFooter returns the footer string for the specified records.
func (self *BufferingFormatter) FormatFooter(_ []*LogRecord) string {
	return ""
}

// Format the specified records and return the result as a a string.
func (self *BufferingFormatter) Format(records []*LogRecord) string {
	var buf bytes.Buffer
	if len(records) > 0 {
		buf.WriteString(self.FormatHeader(records))
		for _, record := range records {
			buf.WriteString(self.lineFormatter.Format(record))
		}
		buf.WriteString(self.FormatFooter(records))
	}
	return buf.String()
}
