package logging

import (
	"github.com/hhkbp2/testify/require"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDictConfig_UnsupportVersion(t *testing.T) {
	conf := &Conf{
		Version: 2,
	}
	require.NotNil(t, DictConfig(conf))
}

func _testConfigLogger(t *testing.T) {
	logger1 := GetLogger("a.b")
	logger2 := GetLogger("a")
	message := "xxxxyyy"
	logger1.Info(message)
	logger1.Error(message)
	logger2.Debug(message)
	logger2.Error(message)
	Shutdown()
	// open the log file and check its content
	content, err := ioutil.ReadFile(testFileName)
	require.Nil(t, err)
	require.Equal(t, strings.Repeat(message+"\n", 2), string(content))
	// clean up the log file
	require.Nil(t, os.Remove(testFileName))
}

func TestApplyJsonConfigFile(t *testing.T) {
	file := "./config_example.json"
	require.Nil(t, ApplyJsonConfigFile(file))
	_testConfigLogger(t)
}

func TestApplyYAMLConfigFile(t *testing.T) {
	file := "./config_example.yml"
	require.Nil(t, ApplyYAMLConfigFile(file))
	_testConfigLogger(t)
}
