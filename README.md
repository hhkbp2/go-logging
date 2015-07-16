go-logging
==========

```go-logging``` is a Golang library that implements the Python-like logging facility. 

As we all know that logging is essientially significant for server side programming because in general logging the only way to report what happens inside the program. 

The [```logging```][python-logging-page] package of Python standard library is a popular logging facility among Pythoners. ```logging``` defines ```Logger``` as logging source, ```Handler``` as logging event destination, and supports Logger hierarchy and free combinations of both.  It is powerful and flexible,  in a similar style like [```Log4j```][log4j-page], which is a popular logging facility among Javaers.

When it comes to Golang, the standard release has a library called [```log```][golang-log-page] for logging. It's simple and good to log something into standard IO or a customized IO. In fact it's too simple to use in any **real** production enviroment, especially when compared to some other mature logging library. 

Due to the lack of a good logging facility, many people start to develop their own versions. For example in github there are dozens of logging repositories for Golang. I run into the same problem when I am writing some project in Golang. A powerful logging facility is needed to develop and debug it. I take a search on a few existing logging libraries for Golang but none of them seems to meet the requirement. So I decide to join the parade of "everyone is busy developing his own version", and then this library is created.

## Features

With an obivious intention to be a port of ```logging``` for Golang, ```go-logging``` has all the main features that ```logging``` package has:

1. It supports logging level, logging sources(logger) and destinations(handler) customization and flexible combinations of them
2. It supports logger hierarchy, optional filter on logger and handler, optional formatter on handler
3. It supports handlers that frequently-used in most real production enviroments, e.g. it could write log events to stdout, memory, file, syslog, udp/tcp socket, rpc(e.g., thrift. For the corresponding servers, please refer to the unit test) etc.
4. It could be configured throught handy config file in various format(e.g. yaml, json)

## Usage

Get this library using the standard go tool:

```bash
go get github.com/hhkbp2/go-logging
```

#### Example 1: Log to standard output

```go
package main

import (
	"github.com/hhkbp2/go-logging"
)

func main() {
	logger := logging.GetLogger("a.b")
	handler := logging.NewStdoutHandler()
	logger.AddHandler(handler)
	logger.Warnf("message: %s %d", "Hello", 2015)
}
```

The code above outputs as the following:

```text
message: Hello 2015
```

#### Example 2: Log to file

```go
package main

import (
	"github.com/hhkbp2/go-logging"
	"os"
	"time"
)

func main() {
	filePath := "./test.log"
	fileMode := os.O_APPEND
	bufferSize := 0
	bufferFlushTime := 30 * time.Second
	inputChanSize := 1
	// set the maximum size of every file to 100 M bytes
	fileMaxBytes := uint64(100 * 1024 * 1024)
	// keep 9 backup at most(including the current using one,
	// there could be 10 log file at most)
	backupCount := uint32(9)
	// create a handler(which represents a log message destination)
	handler := logging.MustNewRotatingFileHandler(
		filePath, fileMode, bufferSize, bufferFlushTime, inputChanSize,
		fileMaxBytes, backupCount)

	// the format for the whole log message
	format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
		"%(name)s %(message)s"
	// the format for the time part
	dateFormat := "%Y-%m-%d %H:%M:%S.%3n"
	// create a formatter(which controls how log messages are formatted)
	formatter := logging.NewStandardFormatter(format, dateFormat)
	// set formatter for handler
	handler.SetFormatter(formatter)

	// create a logger(which represents a log message source)
	logger := logging.GetLogger("a.b.c")
	logger.SetLevel(logging.LevelInfo)
	logger.AddHandler(handler)

	// ensure all log messages are flushed to disk before program exits.
	defer logging.Shutdown()

	logger.Infof("message: %s %d", "Hello", 2015)
}
```

Compile and run the code above, it would generate a log file "./test.log" under current working directory. The log file contains a single line:

```text
2015-04-04 14:20:33.714 INFO (main2.go:40) a.b.c message: Hello 2015
```

#### Example 3: Config Log via configuration file.

Write a configuration file ```config.yml``` as the following:

```go
formatters:
    f:
        format: "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) %(name)s %(message)s"
        datefmt: "%Y-%m-%d %H:%M:%S.%3n"
handlers:
    h:
        class: RotatingFileHandler
        filepath: "./test.log"
        mode: O_APPEND
        bufferSize: 0
        # 30 * 1000 ms -> 30 seconds
        bufferFlushTime: 30000
        inputChanSize: 1
        # 100 * 1024 * 1024 -> 100M
        maxBytes: 104857600
        backupCount: 9
        formatter: f
loggers:
    a.b.c:
        level: INFO
        handlers: [h]

```

and use it to config logging facility like:

```go
package main

import (
	"github.com/hhkbp2/go-logging"
)

func main() {
	config_file := "./config.yml"
	if err := logging.ApplyConfigFile(config_file); err != nil {
		panic(err.Error())
	}
	logger := logging.GetLogger("a.b.c")
	defer logging.Shutdown()
	logger.Infof("message: %s %d", "Hello", 2015)
}
```

It will write log as the same as the above example 2.

## Documentation

For docs, refer to: 

[https://godoc.org/github.com/hhkbp2/go-logging](https://godoc.org/github.com/hhkbp2/go-logging)

For much more details please refer to the documentation for [```logging```][python-logging-page].

[python-logging-page]: https://docs.python.org/2/library/logging.html

[log4j-page]: http://logging.apache.org/log4j/

[golang-log-page]: http://golang.org/pkg/log/

