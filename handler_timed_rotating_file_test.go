package logging

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hhkbp2/go-strftime"
	"github.com/hhkbp2/testify/require"
)

func cleanupLogFils(t *testing.T, basepath string) int {
	dirName, baseName := filepath.Split(basepath)
	fileInfos, err := ioutil.ReadDir(dirName)
	require.Nil(t, err)
	count := 0
	for _, info := range fileInfos {
		name := info.Name()
		if strings.HasPrefix(name, baseName) {
			require.Nil(t, os.Remove(filepath.Join(dirName, name)))
			count += 1
		}
	}
	return count
}

func checkFileContent(t *testing.T, file, content string) {
	c, err := ioutil.ReadFile(file)
	require.Nil(t, err)
	require.Equal(t, content, string(c))
}

func TestTimedRotatingFileHandler_WithBackup(t *testing.T) {
	defer Shutdown()
	cleanupLogFils(t, testFileName)
	when := "S"
	format := "%Y-%m-%d_%H-%M-%S"
	interval := 2
	handler, err := NewTimedRotatingFileHandler(
		testFileName,
		testFileMode,
		testBufferSize,
		testBufferFlushTime,
		testInputChanSize,
		when,
		uint32(interval),
		testRotateBackupCount,
		false)
	require.Nil(t, err)
	logger := GetLogger("trfile")
	logger.AddHandler(handler)
	message := "test"
	lastMessage := "last message"
	times := make([]time.Time, 0, testRotateBackupCount)
	for i := uint32(0); i < testRotateBackupCount+1; i++ {
		logger.Errorf(message)
		if i > 0 {
			times = append(times, time.Now())
		}
		time.Sleep(time.Duration(int64(time.Second) * int64(interval)))
	}
	logger.Errorf(lastMessage)
	logger.RemoveHandler(handler)
	handler.Close()
	for i := uint32(0); i < testRotateBackupCount; i++ {
		suffix := strftime.Format(format, times[i])
		checkFileContent(t, testFileName+"."+suffix, message+"\n")
	}
	checkFileContent(t, testFileName, lastMessage+"\n")
	require.Equal(t, 2, cleanupLogFils(t, testFileName))
}
