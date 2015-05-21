package logging

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrorConfigVersionUnsupport    = errors.New("unsupport version")
	ErrorConfigInvalidType         = errors.New("invalid type")
	ErrorConfigUnknownLevel        = errors.New("unknown level")
	ErrorConfigUnknownFileMode     = errors.New("unknown file mode")
	ErrorConfigInvalidHandler      = errors.New("invalid handler")
	ErrorConfigInvalidFormatter    = errors.New("invalid formatter")
	ErrorConfigInvalidFilter       = errors.New("invalid filter")
	ErrorConfigIDMissing           = errors.New("id missing")
	ErrorConfigIDAlreadyExists     = errors.New("id already exists")
	ErrorConfigHandlerClassMissing = errors.New("handler class missing")
	ErrorConfigUnknownHandlerClass = errors.New("unknown handler class")
	ErrorConfigMapNoSuchKey        = errors.New("map has no such key")
	ErrorConfigMapValueTypeMistach = errors.New("value type mismatch")
	ErrorConfigNoSuchHandler       = errors.New("no such handler")
)

var (
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

type ConfFilter struct {
	Name string `json:"name"`
}

type ConfFormatter struct {
	Format     *string `json:"format"`
	DateFormat *string `json:"datefmt"`
}

type ConfMap map[string]interface{}

func (self ConfMap) GetBool(key string) (bool, error) {
	value, ok := self[key]
	if !ok {
		return false, ErrorConfigMapNoSuchKey
	}
	b, ok := value.(bool)
	if !ok {
		return false, ErrorConfigMapValueTypeMistach
	}
	return b, nil
}

func (self ConfMap) GetInt(key string) (int, error) {
	value, ok := self[key]
	if !ok {
		return 0, ErrorConfigMapNoSuchKey
	}
	i, ok := value.(int)
	if !ok {
		return 0, ErrorConfigMapValueTypeMistach
	}
	return i, nil
}

func (self ConfMap) GetUint16(key string) (uint16, error) {
	value, ok := self[key]
	if !ok {
		return 0, ErrorConfigMapNoSuchKey
	}
	i, ok := value.(uint16)
	if !ok {
		return 0, ErrorConfigMapValueTypeMistach
	}
	return i, nil
}

func (self ConfMap) GetUint32(key string) (uint32, error) {
	value, ok := self[key]
	if !ok {
		return 0, ErrorConfigMapNoSuchKey
	}
	i, ok := value.(uint32)
	if !ok {
		return 0, ErrorConfigMapValueTypeMistach
	}
	return i, nil
}

func (self ConfMap) GetUint64(key string) (uint64, error) {
	value, ok := self[key]
	if !ok {
		return 0, ErrorConfigMapNoSuchKey
	}
	i, ok := value.(uint64)
	if !ok {
		return 0, ErrorConfigMapValueTypeMistach
	}
	return i, nil
}

func (self ConfMap) GetString(key string) (string, error) {
	value, ok := self[key]
	if !ok {
		return "", ErrorConfigMapNoSuchKey
	}
	str, ok := value.(string)
	if !ok {
		return "", ErrorConfigMapValueTypeMistach
	}
	return str, nil
}

type ConfLogger struct {
	Level     string   `json:"level"`
	Propagate bool     `json:"propagate"`
	Filters   []string `json:"filters"`
	Handlers  []string `json:"handlers"`
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

func ApplyJsonConfigFile(file string) error {
	bin, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var conf Conf
	if err = json.Unmarshal(bin, &conf); err != nil {
		return err
	}
	return DictConfig(&conf)
}

func ApplyXMLConfigFile(file string) error {
	// TODO
	return nil
}

type SetLevelable interface {
	SetLevel(level LogLevelType) error
}

func ConfigLevel(m ConfMap, i SetLevelable) error {
	if arg, ok := m["level"]; ok {
		levelIn, ok := arg.(string)
		if !ok {
			return ErrorConfigInvalidType
		}
		levelIn = strings.ToUpper(levelIn)
		level, ok := nameToLevels[levelIn]
		if !ok {
			return ErrorConfigUnknownLevel
		}
		return i.SetLevel(level)
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
			return ErrorConfigInvalidType
		}
		for _, h := range handlersIn {
			name, ok := h.(string)
			if !ok {
				return ErrorConfigInvalidType
			}
			handler, ok := env.handlers[name]
			if !ok {
				return ErrorConfigInvalidHandler
			}
			i.AddHandler(handler)
		}
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
			return ErrorConfigInvalidType
		}
		formatter, ok := env.formatters[name]
		if !ok {
			return ErrorConfigInvalidFormatter
		}
		i.SetFormatter(formatter)
	}
	return nil
}

func ConfigFilters(m ConfMap, i Filterer, env *ConfEnv) error {
	if arg, ok := m["filters"]; ok {
		filtersIn, ok := arg.([]interface{})
		if !ok {
			return ErrorConfigInvalidType
		}
		for _, f := range filtersIn {
			name, ok := f.(string)
			if !ok {
				return ErrorConfigInvalidType
			}
			filter, ok := env.filters[name]
			if !ok {
				return ErrorConfigInvalidFilter
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
				return ErrorConfigInvalidType
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
	if conf.Version != 1 {
		return ErrorConfigVersionUnsupport
	}
	// initialize all filters as specified
	for name, conf := range conf.Filters {
		// reject empty name
		if len(name) == 0 {
			return ErrorConfigIDMissing
		}
		// reject duplicate name
		if _, ok := env.filters[name]; ok {
			return ErrorConfigIDAlreadyExists
		}
		env.filters[name] = NewNameFilter(conf.Name)
	}
	// initialize all formatters as specified
	for name, conf := range conf.Formatters {
		if len(name) == 0 {
			return ErrorConfigIDMissing
		}
		if _, ok := env.formatters[name]; ok {
			return ErrorConfigIDAlreadyExists
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
			return ErrorConfigIDMissing
		}
		if _, ok := env.handlers[name]; ok {
			return ErrorConfigIDAlreadyExists
		}
		arg, ok := m["class"]
		if !ok {
			return ErrorConfigHandlerClassMissing
		}
		className, ok := arg.(string)
		if !ok {
			return ErrorConfigInvalidType
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
				return ErrorConfigUnknownLevel
			}
			handlerName, err := m.GetString("target")
			if err != nil {
				return err
			}
			target, ok := env.handlers[handlerName]
			if !ok {
				return ErrorConfigNoSuchHandler
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
				return ErrorConfigUnknownFileMode
			}
			handler, err = NewFileHandler(filename, mode)
			if err != nil {
				return err
			}
		case "RotatingFileHandler":
			filepath, err := m.GetString("filepath")
			if err != nil {
				return err
			}
			mode, err := m.GetInt("mode")
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
				filepath, mode, maxBytes, backupCount)
			if err != nil {
				return err
			}
		case "TimedRotatingFileHandler":
			filepath, err := m.GetString("filepath")
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
				filepath, when, interval, backupCount, utc)
			if err != nil {
				return err
			}
		case "SyslogHandler":
			priorityStr, err := m.GetString("priority")
			if err != nil {
				return err
			}
			priority, ok := SyslogNameToPriorities[priorityStr]
			if !ok {
				return ErrorConfigMapValueTypeMistach
			}
			tag, err := m.GetString("tag")
			if err != nil {
				return err
			}
			handler, err = NewSyslogHandler(priority, tag)
			if err != nil {
				return err
			}
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
		case "ThriftHandler":
			host, err := m.GetString("host")
			if err != nil {
				return err
			}
			port, err := m.GetUint16("port")
			if err != nil {
				return err
			}
			handler = NewThriftHandler(host, port)
		default:
			return ErrorConfigUnknownHandlerClass
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
			return ErrorConfigIDMissing
		}
		logger := GetLogger(name)
		if err := ConfigLogger(m, logger, false, env); err != nil {
			return err
		}
	}
	return nil
}
