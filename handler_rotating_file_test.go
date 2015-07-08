package logging

import (
	"github.com/hhkbp2/testify/require"
	"os"
	"strings"
	"testing"
)

var (
	testRotateMaxBytes    = uint64(5 * 1000) // 5k bytes
	testRotateBackupCount = uint32(1)
)

func checkFileSize(t *testing.T, filename string) {
	fileInfo, err := os.Stat(filename)
	require.Nil(t, err)
	require.True(t, fileInfo.Size() > 0)
	require.True(t, uint64(fileInfo.Size()) <= testRotateMaxBytes)
}

func removeFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	require.Nil(t, err)
}

func TestRotatingFileHandler_TruncateWithBackup(t *testing.T) {
	defer Shutdown()
	handler, err := NewRotatingFileHandler(
		testFileName, testFileMode, testRotateMaxBytes, testRotateBackupCount)
	require.Nil(t, err)
	logger := GetLogger("rfile1")
	logger.AddHandler(handler)
	// every message is 99 byte, and \n
	message := strings.Repeat("abcdefghij", 9) + "rstuvwxyz"
	size := uint64(len(message) + 1)
	total := testRotateMaxBytes * (uint64(testRotateBackupCount) + 2) / size
	for i := uint64(0); i < total; i++ {
		logger.Errorf(message)
	}
	logger.RemoveHandler(handler)
	handler.Close()
	checkFileSize(t, testFileName)
	checkFileSize(t, testFileName+".1")
	removeFile(t, testFileName)
	removeFile(t, testFileName+".1")
}

func TestRotatingFileHandler_AppendWithoutBackup(t *testing.T) {
	defer Shutdown()
	// clean up the existing log file
	if FileExists(testFileName) {
		require.Nil(t, os.Remove(testFileName))
	}
	backupCount := uint32(0)
	handler, err := NewRotatingFileHandler(
		testFileName, os.O_APPEND, testRotateMaxBytes, backupCount)
	require.Nil(t, err)
	logger := GetLogger("rfile2")
	logger.AddHandler(handler)
	message := strings.Repeat("abcdefghij", 9) + "rstuvwxyz"
	size := uint64(len(message) + 1)
	totalSize := testRotateMaxBytes * (uint64(testRotateBackupCount) + 2)
	times := totalSize / size
	for i := uint64(0); i < times; i++ {
		logger.Errorf(message)
	}
	logger.RemoveHandler(handler)
	handler.Close()
	fileInfo, err := os.Stat(testFileName)
	require.Nil(t, err)
	require.True(t, uint64(fileInfo.Size()) > testRotateMaxBytes)
	require.Equal(t, totalSize, uint64(fileInfo.Size()))
	removeFile(t, testFileName)
}

func BenchmarkRotatingFileHandler(b *testing.B) {
	b.StopTimer()
	defer Shutdown()
	if FileExists(testFileName) {
		os.Remove(testFileName)
	}
	rotateMaxBytes := uint64(1024 * 1024 * 1024 * 1) // 1G
	backupCount := uint32(0)
	handler, err := NewRotatingFileHandler(
		testFileName, os.O_APPEND, rotateMaxBytes, backupCount)
	if err != nil {
		panic("fail to get handler")
	}
	logger := GetLogger("rfileBen")
	logger.AddHandler(handler)
	message := strings.Repeat("abcdefghij", 10)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.Errorf(message)
	}
}
