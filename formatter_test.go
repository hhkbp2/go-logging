package logging

import (
    "fmt"
    "github.com/hhkbp2/testify/require"
    "testing"
)

var (
    testRecord = NewLogRecord(
        "name",
        LevelInfo,
        "pathname",
        "filename",
        111,
        "funcname",
        "message",
        nil)
)

func TestFormat_Name(t *testing.T) {
    require.Equal(t, testRecord.Name, Format("%(name)s", testRecord))
}

func TestFormat_LevelNo(t *testing.T) {
    levelNo := fmt.Sprintf("%d", testRecord.Level)
    require.Equal(t, levelNo, Format("%(levelno)d", testRecord))
}

func TestFormat_LevelName(t *testing.T) {
    levelName := GetLevelName(testRecord.Level)
    require.Equal(t, levelName, Format("%(levelname)s", testRecord))
}

func TestFormat_PathName(t *testing.T) {
    require.Equal(t, testRecord.PathName, Format("%(pathname)s", testRecord))
}

func TestFormat_FileName(t *testing.T) {
    require.Equal(t, testRecord.FileName, Format("%(filename)s", testRecord))
}

func TestFormat_LineNo(t *testing.T) {
    lineNo := fmt.Sprintf("%d", testRecord.LineNo)
    require.Equal(t, lineNo, Format("%(lineno)d", testRecord))
}

func TestFormat_FuncName(t *testing.T) {
    require.Equal(t, testRecord.FuncName, Format("%(funcname)s", testRecord))
}

func TestFormat_Created(t *testing.T) {
    created := fmt.Sprintf("%d", testRecord.CreatedTime.UnixNano())
    require.Equal(t, created, Format("%(created)d", testRecord))
}

func TestFormat_AscTime(t *testing.T) {
    require.Equal(t, testRecord.AscTime, Format("%(asctime)s", testRecord))
}

func TestFormat_Message(t *testing.T) {
    require.Equal(t, testRecord.GetMessage(), Format(defaultFormat, testRecord))
}
