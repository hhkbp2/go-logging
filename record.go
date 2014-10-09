package logging

import (
    "fmt"
    "time"
)

type LogRecord struct {
    CreatedTime time.Time
    AscTime     string

    Name     string
    Level    LogLevelType
    PathName string
    FileName string
    LineNo   uint32
    Message  string
    Args     []interface{}
}

func NewLogRecord(
    name string,
    level LogLevelType,
    pathName string,
    fileName string,
    lineNo uint32,
    message string,
    args []interface{}) *LogRecord {

    return &LogRecord{
        CreatedTime: time.Now(),
        Name:        name,
        Level:       level,
        PathName:    pathName,
        FileName:    fileName,
        LineNo:      lineNo,
        Message:     message,
        Args:        args,
    }
}

func (self *LogRecord) String() string {
    return fmt.Sprintf("<LogRecord: %s, %s, %s, %s, \"%s\">",
        self.Name, self.Level, self.PathName, self.LineNo, self.Message)
}

func (self *LogRecord) GetMessage() string {
    if self.Args != nil {
        return fmt.Sprintf(self.Message, self.Args...)
    }
    return self.Message
}
