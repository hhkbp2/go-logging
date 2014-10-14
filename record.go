package logging

import (
    "fmt"
    "time"
)

// A LogRecord instance represents and event being logged.
// LogRecord instances are created every time something is logged. They
// contain all the information pretinent to the event being logged. The
// main information passed in is in Message and Args, which are combined
// using fmt.Sprintf() to create the message field of the record. The
// record also includes information such as when the record was created,
// the source line where the logging call was made, and any exception
// information to be logged.
type LogRecord struct {
    CreatedTime time.Time
    AscTime     string

    Name     string
    Level    LogLevelType
    PathName string
    FileName string
    LineNo   uint32
    FuncName string
    Format   string
    Args     []interface{}
    Message  string
}

// Initialize a logging record with interesting information.
func NewLogRecord(
    name string,
    level LogLevelType,
    pathName string,
    fileName string,
    lineNo uint32,
    funcName string,
    format string,
    args []interface{}) *LogRecord {

    return &LogRecord{
        CreatedTime: time.Now(),
        Name:        name,
        Level:       level,
        PathName:    pathName,
        FileName:    fileName,
        LineNo:      lineNo,
        FuncName:    funcName,
        Format:      format,
        Args:        args,
        Message:     "",
    }
}

// Return the string representation for this LogRecord.
func (self *LogRecord) String() string {
    return fmt.Sprintf("<LogRecord: %s, %s, %s, %s, \"%s\">",
        self.Name, self.Level, self.PathName, self.LineNo, self.Message)
}

// Return the message for this LogRecord.
// The message is composed of the Message and any user-supplied arguments.
func (self *LogRecord) GetMessage() string {
    if self.Args != nil {
        return fmt.Sprintf(self.Format, self.Args...)
    }
    return self.Format
}
