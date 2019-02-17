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

type RecordValueReceiver func(record *LogRecord) interface{}

func GetValueForField(fieldName string, record *LogRecord) interface{} {
	var value interface{}

	switch fieldName {
	case "name":
		value = record.Name
	case "levelno":
		value = record.Level
	case "levelname":
		value = record.Level
	case "pathname":
		value = record.PathName
	case "filename":
		value = record.FileName
	case "lineno":
		value = fmt.Sprintf("%d", record.LineNo)
	case "funcname":
		value = record.FuncName
	case "created":
		value = fmt.Sprintf("%d", record.CreatedTime.UnixNano())
	case "asctime":
		value = record.AscTime
	case "message":
		value = record.Message
	default:
		ctxField, ok := record.CtxFields[fieldName]
		if !ok {
			// no value for context field 'fieldName'
			// return default
			ctxField = "null"
		}
		value = ctxField
	}

	return value
}

// All predefined attribute string and their ExtractAttr functions.
var (
	formatRe = initFormatRegexp()

	// Default format strings.
	defaultFormat     = "%(message)s"
	defaultDateFormat = "%Y-%m-%d %H:%M:%S %3n"
	defaultFormatter  = NewStandardFormatter(defaultFormat, defaultDateFormat)
)

// Initialize global regexp for attribute matching.
func initFormatRegexp() *regexp.Regexp {
	re := `(%\([_\w]+\))(s|d)?`
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
	keywordRe := `\(([_\w]+)\)`
	keywordRegexp := regexp.MustCompile(keywordRe)
	var replaceFuncs []RecordValueReceiver
	f1 := func(match string) string {
		if match == "%%" {
			return "%"
		}
		if match == "%(asctime)s" {
			toFormatTime = true
		}
		size++
		matches := keywordRegexp.FindStringSubmatch(match)
		keyword := matches[1]

		valueReceiver := func(fieldName string) RecordValueReceiver {
			var extractFunc RecordValueReceiver
			extractFunc = func(record *LogRecord) interface{} {
				return GetValueForField(fieldName, record)
			}
			return extractFunc
		}(keyword)
		replaceFuncs = append(replaceFuncs, valueReceiver)

		if keyword == "levelno" {
			return "%d"
		}
		return "%s"
	}
	strFormat := formatRe.ReplaceAllStringFunc(format, f1)

	getFormatArgsFunc := func(record *LogRecord) []interface{} {
		result := make([]interface{}, 0, len(replaceFuncs))
		for _, f := range replaceFuncs {
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

// Return the creation time of the specified LogRecord as formatted text.
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

// Return the header string for the specified records.
func (self *BufferingFormatter) FormatHeader(_ []*LogRecord) string {
	return ""
}

// Return the footer string for the specified records.
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
