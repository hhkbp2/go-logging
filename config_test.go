package logging

import (
	"github.com/hhkbp2/testify/require"
	"testing"
)

func TestApplyJsonConfigFile(t *testing.T) {
	file := "./config_example.json"
	require.Nil(t, ApplyJsonConfigFile(file))
}

func TestDictConfig_UnsupportVersion(t *testing.T) {
	conf := &Conf{
		Version: 2,
	}
	require.NotNil(t, DictConfig(conf))
}
