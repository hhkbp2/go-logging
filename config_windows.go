package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	// A map from string description to file modes.
	// The string descriptions are used in configuration file.
	FileModeNameToValues = map[string]int{
		"O_RDONLY": os.O_RDONLY,
		"O_WRONLY": os.O_WRONLY,
		"O_RDWR":   os.O_RDWR,
		"O_APPEND": os.O_APPEND,
		"O_CREATE": os.O_CREATE,
		"O_EXCL":   os.O_EXCL,
		"O_SYNC":   os.O_SYNC,
		"O_TRUNC":  os.O_TRUNC,
	}
)

// Apply all configuration in specified file.
func ApplyConfigFile(file string) error {
	ext := filepath.Ext(file)
	switch ext {
	case ".json":
		return ApplyJsonConfigFile(file)
	case ".yml":
		fallthrough
	case ".yaml":
		return ApplyYAMLConfigFile(file)
	default:
		return errors.New(fmt.Sprintf(
			"unknown format of the specified file: %s", file))
	}
}

// Apply all configuration in specified json file.
func ApplyJsonConfigFile(file string) error {
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(bin))
	decoder.UseNumber()
	var conf Conf
	if err = decoder.Decode(&conf); err != nil {
		return err
	}
	return DictConfig(&conf)
}

// Apply all configuration in specified yaml file.
func ApplyYAMLConfigFile(file string) error {
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var conf Conf
	if err = yaml.Unmarshal(bin, &conf); err != nil {
		return err
	}
	return DictConfig(&conf)
}

type ConfFilter struct {
	Name string `json:"name"`
}

type ConfFormatter struct {
	Format     *string `json:"format"`
	DateFormat *string `json:"datefmt"`
}

// A map represents configuration of various key and variable length.
// Access it just in raw key after parsed from config file.
type ConfMap map[string]interface{}

func (self ConfMap) GetBool(key string) (bool, error) {
	value, ok := self[key]
	if !ok {
		return false, errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	if v, ok := value.(bool); ok {
		return v, nil
	}
	if str, ok := value.(string); ok {
		switch {
		case strings.EqualFold(str, "true"):
			return true, nil
		case strings.EqualFold(str, "false"):
			return false, nil
		}
	}
	if n, ok := value.(json.Number); ok {
		v, err := n.Int64()
		if err != nil {
			return false, err
		}
		if v != 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, errors.New(fmt.Sprintf(
		"value: %#v of key: %s should be of type bool not type %s",
		value, key, reflect.TypeOf(value)))
}

func (self ConfMap) GetInt(key string) (int, error) {
	value, ok := self[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	if v, ok := value.(int); ok {
		return v, nil
	}
	if str, ok := value.(string); ok {
		return strconv.Atoi(str)
	}
	if n, ok := value.(json.Number); ok {
		v, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return int(v), nil
	}
	return 0, errors.New(fmt.Sprintf(
		"value: %#v of key: %s should be of type int not type %s",
		value, key, reflect.TypeOf(value)))
}

func (self ConfMap) GetUint16(key string) (uint16, error) {
	value, ok := self[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	if v, ok := value.(int); ok {
		return uint16(v), nil
	}
	if str, ok := value.(string); ok {
		v, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		return uint16(v), nil
	}
	if n, ok := value.(json.Number); ok {
		v, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return uint16(v), nil
	}
	return 0, errors.New(fmt.Sprintf(
		"value: %#v of key: %s should be of type uint16 not type %s",
		value, key, reflect.TypeOf(value)))
}

func (self ConfMap) GetUint32(key string) (uint32, error) {
	value, ok := self[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	if v, ok := value.(int); ok {
		return uint32(v), nil
	}
	if str, ok := value.(string); ok {
		v, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	}
	if n, ok := value.(json.Number); ok {
		v, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return uint32(v), nil
	}
	return 0, errors.New(fmt.Sprintf(
		"value: %#v of key: %s should be of type uint32 not type %s",
		value, key, reflect.TypeOf(value)))
}

func (self ConfMap) GetUint64(key string) (uint64, error) {
	value, ok := self[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	if v, ok := value.(int); ok {
		return uint64(v), nil
	}
	if str, ok := value.(string); ok {
		v, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		return uint64(v), nil
	}
	if n, ok := value.(json.Number); ok {
		v, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return uint64(v), nil
	}
	return 0, errors.New(fmt.Sprintf(
		"value: %#v of key: %s should be of type uint64 not type %s",
		value, key, reflect.TypeOf(value)))
}

func (self ConfMap) GetString(key string) (string, error) {
	value, ok := self[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("no config for key: %s", key))
	}
	str, ok := value.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf(
			"value: %#v of key: %s should be of type string", value, key))
	}
	return str, nil
}

type Conf struct {
	Version    int                      `json:"version"`
	Root       ConfMap                  `json:"root"`
	Loggers    map[string]ConfMap       `json:"loggers"`
	Handlers   map[string]ConfMap       `json:"handlers"`
	Formatters map[string]ConfFormatter `json:"formatters"`
	Filters    map[string]ConfFilter    `json:"filters"`
}

type ConfEnv struct {
	handlers   map[string]Handler
	formatters map[string]Formatter
	filters    map[string]Filter
}

func NewConfigEnv() *ConfEnv {
	return &ConfEnv{
		handlers:   make(map[string]Handler),
		formatters: make(map[string]Formatter),
		filters:    make(map[string]Filter),
	}
}

type SetLevelable interface {
	SetLevel(level LogLevelType) error
}

func ConfigLevel(m ConfMap, i SetLevelable) error {
	if arg, ok := m["level"]; ok {
		levelIn, ok := arg.(string)
		if !ok {
			return errors.New(fmt.Sprintf(
				"level value: %#v should be of type string", arg))
		}
		levelIn = strings.ToUpper(levelIn)
		level, ok := nameToLevels[levelIn]
		if !ok {
			return errors.New(fmt.Sprintf("unknown level: %s", levelIn))
		}
		return i.SetLevel(level)
	}
	return nil
}

type SetFormatterable interface {
	SetFormatter(formatter Formatter)
}

func ConfigFormatters(m ConfMap, i SetFormatterable, env *ConfEnv) error {
	if arg, ok := m["formatter"]; ok {
		name, ok := arg.(string)
		if !ok {
			return errors.New(fmt.Sprintf(
				"formatter value: %#v should be of type string", arg))
		}
		formatter, ok := env.formatters[name]
		if !ok {
			return errors.New(fmt.Sprintf("unknown formatter: %s", name))
		}
		i.SetFormatter(formatter)
	}
	return nil
}

type AddHandlerable interface {
	AddHandler(handler Handler)
}

func ConfigHandlers(m ConfMap, i AddHandlerable, env *ConfEnv) error {
	if arg, ok := m["handlers"]; ok {
		handlersIn, ok := arg.([]interface{})
		if !ok {
			return errors.New(fmt.Sprintf(
				"handlers value: %#v should be of type string slice", arg))
		}
		for _, h := range handlersIn {
			name, ok := h.(string)
			if !ok {
				return errors.New(fmt.Sprintf(
					"%#v in handlers should be of type string", h))
			}
			handler, ok := env.handlers[name]
			if !ok {
				return errors.New(fmt.Sprintf("unknown handler: %s", name))
			}
			i.AddHandler(handler)
		}
	}
	return nil
}

func ConfigFilters(m ConfMap, i Filterer, env *ConfEnv) error {
	if arg, ok := m["filters"]; ok {
		filtersIn, ok := arg.([]interface{})
		if !ok {
			return errors.New(fmt.Sprintf(
				"filters value: %#v should be of type string slice", arg))
		}
		for _, f := range filtersIn {
			name, ok := f.(string)
			if !ok {
				return errors.New(fmt.Sprintf(
					"%#v in filters should be of type string", f))
			}
			filter, ok := env.filters[name]
			if !ok {
				return errors.New(fmt.Sprintf("unknown filter: %s", name))
			}
			i.AddFilter(filter)
		}
	}
	return nil
}

func ConfigLogger(m ConfMap, logger Logger, isRoot bool, env *ConfEnv) error {
	if err := ConfigLevel(m, logger); err != nil {
		return err
	}
	// propagate setting will not be applicable to root logger
	if !isRoot {
		if arg, ok := m["propagate"]; ok {
			propagate, ok := arg.(bool)
			if !ok {
				return errors.New(fmt.Sprintf(
					"propagate value: %#v should be of type bool", arg))
			}
			logger.SetPropagate(propagate)
		}
	}
	if err := ConfigHandlers(m, logger, env); err != nil {
		return err
	}
	return ConfigFilters(m, logger, env)
}

func DictConfig(conf *Conf) error {
	env := NewConfigEnv()
	// check version for compatibility.  Currently only version 1 is supported.
	if (conf.Version != 0) && (conf.Version != 1) {
		return errors.New(fmt.Sprintf("unsupport version: %d", conf.Version))
	}
	// initialize all filters as specified
	for name, conf := range conf.Filters {
		// reject empty name
		if len(name) == 0 {
			return errors.New("filter should have non-empty ID")
		}
		// reject duplicate name
		if _, ok := env.filters[name]; ok {
			return errors.New(fmt.Sprintf("filter id: %s already exists", name))
		}
		env.filters[name] = NewNameFilter(conf.Name)
	}
	// initialize all formatters as specified
	for name, conf := range conf.Formatters {
		if len(name) == 0 {
			return errors.New("formatter should have non-empty ID")
		}
		if _, ok := env.formatters[name]; ok {
			return errors.New(fmt.Sprintf(
				"formatter id: %s already exists", name))
		}
		var format, dateFormat string
		if conf.Format != nil {
			format = *conf.Format
		} else {
			format = defaultFormat
		}
		if conf.DateFormat != nil {
			dateFormat = *conf.DateFormat
		} else {
			dateFormat = defaultDateFormat
		}
		env.formatters[name] = NewStandardFormatter(format, dateFormat)
	}
	// initialize all handlers as specified
	for name, m := range conf.Handlers {
		if len(name) == 0 {
			return errors.New("handler should have non-empty ID")
		}
		if _, ok := env.handlers[name]; ok {
			return errors.New(fmt.Sprintf(
				"handler id: %s already exists", name))
		}
		arg, ok := m["class"]
		if !ok {
			return errors.New(fmt.Sprintf(
				"handler id: %s should specify class", name))
		}
		className, ok := arg.(string)
		if !ok {
			return errors.New(fmt.Sprintf(
				"handler id: %s class should be of type string"))
		}
		var handler Handler
		switch className {
		case "NullHandler":
			handler = NewNullHandler()
		case "MemoryHandler":
			capacity, err := m.GetUint64("capacity")
			if err != nil {
				return err
			}
			levelStr, err := m.GetString("level")
			if err != nil {
				return err
			}
			levelStr = strings.ToUpper(levelStr)
			level, ok := nameToLevels[levelStr]
			if !ok {
				return errors.New(fmt.Sprintf("unknown level: %s", levelStr))
			}
			handlerName, err := m.GetString("target")
			if err != nil {
				return err
			}
			target, ok := env.handlers[handlerName]
			if !ok {
				return errors.New(fmt.Sprintf(
					"target handler id: %s not exists", handlerName))
			}
			handler = NewMemoryHandler(capacity, level, target)
		case "StdoutHandler":
			handler = NewStdoutHandler()
		case "FileHandler":
			filename, err := m.GetString("filename")
			if err != nil {
				return err
			}
			modeStr, err := m.GetString("mode")
			if err != nil {
				return err
			}
			mode, ok := FileModeNameToValues[modeStr]
			if !ok {
				return errors.New(fmt.Sprintf("unknown file mode: %s", modeStr))
			}
			bufferSize, err := m.GetInt("bufferSize")
			if err != nil {
				return err
			}
			handler, err = NewFileHandler(filename, mode, bufferSize)
			if err != nil {
				return err
			}
		case "RotatingFileHandler":
			filepath, err := m.GetString("filepath")
			if err != nil {
				return err
			}
			modeStr, err := m.GetString("mode")
			if err != nil {
				return err
			}
			mode, ok := FileModeNameToValues[modeStr]
			if !ok {
				return errors.New(fmt.Sprintf("unknown file mode: %s", modeStr))
			}
			bufferSize, err := m.GetInt("bufferSize")
			if err != nil {
				return err
			}
			bufferFlushTimeMS, err := m.GetInt("bufferFlushTime")
			if err != nil {
				return err
			}
			bufferFlushTime := time.Millisecond * time.Duration(bufferFlushTimeMS)
			inputChanSize, err := m.GetInt("inputChanSize")
			if err != nil {
				return err
			}
			maxBytes, err := m.GetUint64("maxBytes")
			if err != nil {
				return err
			}
			backupCount, err := m.GetUint32("backupCount")
			if err != nil {
				return err
			}
			handler, err = NewRotatingFileHandler(
				filepath,
				mode,
				bufferSize,
				bufferFlushTime,
				inputChanSize,
				maxBytes,
				backupCount)
			if err != nil {
				return err
			}
		case "TimedRotatingFileHandler":
			filepath, err := m.GetString("filepath")
			if err != nil {
				return err
			}
			modeStr, err := m.GetString("mode")
			if err != nil {
				return err
			}
			mode, ok := FileModeNameToValues[modeStr]
			if !ok {
				return errors.New(fmt.Sprintf("unknown file mode: %s", modeStr))
			}
			bufferSize, err := m.GetInt("bufferSize")
			if err != nil {
				return err
			}
			when, err := m.GetString("when")
			if err != nil {
				return err
			}
			interval, err := m.GetUint32("interval")
			if err != nil {
				return err
			}
			backupCount, err := m.GetUint32("backupCount")
			if err != nil {
				return err
			}
			utc, err := m.GetBool("utc")
			if err != nil {
				return err
			}
			handler, err = NewTimedRotatingFileHandler(
				filepath, mode, bufferSize, when, interval, backupCount, utc)
			if err != nil {
				return err
			}
		case "DatagramHandler":
			host, err := m.GetString("host")
			if err != nil {
				return err
			}
			port, err := m.GetUint16("port")
			if err != nil {
				return err
			}
			handler = NewDatagramHandler(host, port)
		case "SocketHandler":
			host, err := m.GetString("host")
			if err != nil {
				return err
			}
			port, err := m.GetUint16("port")
			if err != nil {
				return err
			}
			handler = NewSocketHandler(host, port)
		default:
			return errors.New(fmt.Sprintf("unsupported class name: %s", className))
		}
		if err := ConfigLevel(m, handler); err != nil {
			return err
		}
		if err := ConfigFormatters(m, handler, env); err != nil {
			return err
		}
		if err := ConfigFilters(m, handler, env); err != nil {
			return err
		}
		env.handlers[name] = handler
	}
	// set root logger
	if len(conf.Root) > 0 {
		if err := ConfigLogger(conf.Root, root, true, env); err != nil {
			return err
		}
	}
	// initialize all loggers as specified
	for name, m := range conf.Loggers {
		if len(name) == 0 {
			return errors.New("logger should have non-empty ID")
		}
		logger := GetLogger(name)
		if err := ConfigLogger(m, logger, false, env); err != nil {
			return err
		}
	}
	return nil
}
