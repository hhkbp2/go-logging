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
		"",
		false,
		[]interface{}{"message"})
)

func TestFormat_Name(t *testing.T) {
	formatter := NewStandardFormatter("%(name)s", "")
	require.Equal(t, testRecord.Name, formatter.Format(testRecord))
}

func TestFormat_LevelNo(t *testing.T) {
	formatter := NewStandardFormatter("%(levelno)d", "")
	levelNo := fmt.Sprintf("%d", testRecord.Level)
	require.Equal(t, levelNo, formatter.Format(testRecord))
}

func TestFormat_LevelName(t *testing.T) {
	formatter := NewStandardFormatter("%(levelname)s", "")
	levelName := GetLevelName(testRecord.Level)
	require.Equal(t, levelName, formatter.Format(testRecord))
}

func TestFormat_PathName(t *testing.T) {
	formatter := NewStandardFormatter("%(pathname)s", "")
	require.Equal(t, testRecord.PathName, formatter.Format(testRecord))
}

func TestFormat_FileName(t *testing.T) {
	formatter := NewStandardFormatter("%(filename)s", "")
	require.Equal(t, testRecord.FileName, formatter.Format(testRecord))
}

func TestFormat_LineNo(t *testing.T) {
	formatter := NewStandardFormatter("%(lineno)d", "")
	lineNo := fmt.Sprintf("%d", testRecord.LineNo)
	require.Equal(t, lineNo, formatter.Format(testRecord))
}

func TestFormat_FuncName(t *testing.T) {
	formatter := NewStandardFormatter("%(funcname)s", "")
	require.Equal(t, testRecord.FuncName, formatter.Format(testRecord))
}

func TestFormat_Created(t *testing.T) {
	formatter := NewStandardFormatter("%(created)d", "")
	created := fmt.Sprintf("%d", testRecord.CreatedTime.UnixNano())
	require.Equal(t, created, formatter.Format(testRecord))
}

func TestFormat_AscTime(t *testing.T) {
	formatter := NewStandardFormatter("%(asctime)s", defaultDateFormat)
	require.Equal(t, testRecord.AscTime, formatter.Format(testRecord))
}

func TestFormat_Message(t *testing.T) {
	formatter := defaultFormatter
	require.Equal(t, testRecord.GetMessage(), formatter.Format(testRecord))
}
