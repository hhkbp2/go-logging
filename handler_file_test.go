package logging

import (
	"github.com/hhkbp2/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testFileName   = "./test.log"
	testFileMode   = os.O_TRUNC
	testBufferSize = 0
)

func TestFileHandler(t *testing.T) {
	handler, err := NewFileHandler(testFileName, testFileMode, testBufferSize)
	require.Nil(t, err)
	logger := GetLogger("file1")
	logger.AddHandler(handler)
	message := "test"
	logger.Errorf(message)
	logger.RemoveHandler(handler)
	handler.Close()
	// open the log file and check its content equals to message
	// then clean it up.
	content, err := ioutil.ReadFile(testFileName)
	require.Nil(t, err)
	require.Equal(t, message+"\n", string(content))
	err = os.Remove(testFileName)
	require.Nil(t, err)
}

func TestFileHandler_Asctime(t *testing.T) {
	handler, err := NewFileHandler(testFileName, testFileMode, testBufferSize)
	formatter := NewStandardFormatter(
		"%(asctime)s %(message)s",
		"%Y-%m-%d %H:%M:%S %3n")
	handler.SetFormatter(formatter)
	require.Nil(t, err)
	logger := GetLogger("file2")
	logger.AddHandler(handler)
	message := "test"
	logger.Errorf(message)
	logger.RemoveHandler(handler)
	handler.Close()
	// open the log file and check its content equals to message
	// then clean it up.
	content, err := ioutil.ReadFile(testFileName)
	require.Nil(t, err)
	require.Equal(t, 24+len(message)+1, len(content))
	err = os.Remove(testFileName)
	require.Nil(t, err)
}
