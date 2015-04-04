go-logging
==========

```go-logging``` is a Golang library that implements the Python-like logging facility. 

As we all know that logging is essientially significant for server side programming because in general logging the only way to report what happens inside the program. 

The [```logging```][python-logging-page] package of Python standard library is a popular logging facility among Pythoners. ```logging``` defines ```Logger``` as logging source, ```Handler``` as logging event destination, and supports Logger hierarchy and free combinations of both.  It is powerful and flexible,  in a similar style like [```Log4j```][log4j-page], which is a popular logging facility among Javaers.

When it comes to Golang, the standard release has a library called [```log```][golang-log-page] for logging. It's simple and good to log something into standard IO or a customized IO. In fact it's too simple to use in any **real** production enviroment, especially when compared to some other mature logging library. 

Due to the lack of a good logging facility, many people start to develop their own versions. For example in github there are dozens of logging repositories for Golang. I run into the same problem when I am writing [```rafted```][rafted-github]. A powerful logging facility is needed to develop and debug it. I take a search on a few existing logging libraries for Golang but none of them seems to meet the requirement. So I decide to join the parade of "everyone is busy developing his own version", and then this library is created.

## Features

With an obivious intention to be a port of ```logging``` for Golang, ```go-logging``` has all the main features that ```logging``` package has:

1. It supports logging level, logging sources(logger) and destinations(handler) customization and flexible combinations of them
2. It supports logger hierarchy, optional filter on logger and handler, optional formatter on handler
3. It supports handlers that frequently-used in most real production enviroments, e.g. it could write log events to file, syslog, socket, rpc(the corresponding servers are also provided as bundled) etc.
4. It could be configured throught handy config file in various format

Please note that 3, 4 are under development and not fully ready at this moment.

## Usage

Get the code down using the standard go tool:

```bash
go get github.com/hhkbp2/go-logging
```

and write your code just like:

```go
package main
import (
       "github.com/hhkbp2/go-logging"
       "github.com/hhkbp2/go-logging/handlers"
)

func main() {
     logger := logging.GetLogger("a.b")
     handler := handlers.NewTerminalHandler()
     logger.AddHandler(handler)     
     logger.Warnf("test message")
}
```

For more examples please refer to the test cases in source code.

## Documentation

For docs, refer to: 

[https://gowalker.org/github.com/hhkbp2/go-logging](https://gowalker.org/github.com/hhkbp2/go-logging)  
[https://gowalker.org/github.com/hhkbp2/go-logging/handlers](https://gowalker.org/github.com/hhkbp2/go-logging/handlers)

For much more details please refer to the documentation for [```logging```][python-logging-page].

[python-logging-page]: https://docs.python.org/2/library/logging.html

[log4j-page]: http://logging.apache.org/log4j/

[golang-log-page]: http://golang.org/pkg/log/

[rafted-github]: http://github.com/hhkbp2/rafted

