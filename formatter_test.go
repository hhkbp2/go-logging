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
		CtxFields{"ctxid": "123"},
		[]interface{}{"message"},
	)
)

func WithLineFeed(s string) string {
	return s + "\n"
}

func TestFormat_Name(t *testing.T) {
	formatter := NewStandardFormatter("%(name)s", "")
	require.Equal(t,
		WithLineFeed(testRecord.Name), formatter.Format(testRecord))
}

func TestFormat_LevelNo(t *testing.T) {
	formatter := NewStandardFormatter("%(levelno)d", "")
	levelNo := WithLineFeed(fmt.Sprintf("%d", testRecord.Level))
	require.Equal(t, levelNo, formatter.Format(testRecord))
}

func TestFormat_LevelName(t *testing.T) {
	formatter := NewStandardFormatter("%(levelname)s", "")
	levelName := WithLineFeed(GetLevelName(testRecord.Level))
	require.Equal(t, levelName, formatter.Format(testRecord))
}

func TestFormat_PathName(t *testing.T) {
	formatter := NewStandardFormatter("%(pathname)s", "")
	require.Equal(t,
		WithLineFeed(testRecord.PathName), formatter.Format(testRecord))
}

func TestFormat_FileName(t *testing.T) {
	formatter := NewStandardFormatter("%(filename)s", "")
	require.Equal(t,
		WithLineFeed(testRecord.FileName), formatter.Format(testRecord))
}

func TestFormat_LineNo(t *testing.T) {
	formatter := NewStandardFormatter("%(lineno)d", "")
	lineNo := WithLineFeed(fmt.Sprintf("%d", testRecord.LineNo))
	require.Equal(t, lineNo, formatter.Format(testRecord))
}

func TestFormat_FuncName(t *testing.T) {
	formatter := NewStandardFormatter("%(funcname)s", "")
	require.Equal(t,
		WithLineFeed(testRecord.FuncName), formatter.Format(testRecord))
}

func TestFormat_Created(t *testing.T) {
	formatter := NewStandardFormatter("%(created)d", "")
	created := WithLineFeed(
		fmt.Sprintf("%d", testRecord.CreatedTime.UnixNano()))
	require.Equal(t, created, formatter.Format(testRecord))
}

func TestFormat_AscTime(t *testing.T) {
	formatter := NewStandardFormatter("%(asctime)s", defaultDateFormat)
	testRecord.AscTime = formatter.FormatTime(testRecord)
	require.Equal(t,
		WithLineFeed(testRecord.AscTime), formatter.Format(testRecord))
}

func TestFormat_Message(t *testing.T) {
	formatter := defaultFormatter
	require.Equal(t,
		WithLineFeed(testRecord.GetMessage()), formatter.Format(testRecord))
}

func TestFormat_TXID(t *testing.T) {
	formatter := NewStandardFormatter("%(message)s TX_ID=%(ctxid)s", "")
	expectedStr := fmt.Sprintf("message TX_ID=123")
	require.Equal(t,
		WithLineFeed(expectedStr), formatter.Format(testRecord))
}
