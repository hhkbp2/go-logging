package handlers

import (
    "github.com/hhkbp2/go-logging"
    "github.com/hhkbp2/testify/require"
    "io/ioutil"
    "os"
    "testing"
)

var (
    testFileName = "./test.log"
    testFileMode = os.O_TRUNC
)

func TestFileHandler(t *testing.T) {
    handler, err := NewFileHandler(testFileName, testFileMode)
    require.Nil(t, err)
    logger := logging.GetLogger("a")
    logger.AddHandler(handler)
    message := "test"
    logger.Error(message)
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
