go-logging
==========

```go-logging``` is a Golang library that implements the Python-like logging facility. 

The [```logging```][python-logging-page] package of Python standard library is a popular logging facility among Pythoners. ```logging``` defines ```Logger``` as logging source, ```Handler``` as logging event destination, and supports Logger hierarchy and free combinations of both.  It is powerful and flexible,  in a similar style like [```Log4j```][log4j-page], which is a popular logging facility among Javaers.

When it comes to Golang, the standard release has a library called [```log```][golang-log-page] for logging. It's simple and good to log something into standard IO or a customized IO. In fact it's too simple to use in any **real** production enviroment, especially when compared to some other mature logging library. It's somewhat confusing Golang has such a naive Log library but targets itself as a server side programming language. As we all know that logging is essientially significant for server side programming because in general logging the only way to report what happens inside the program.

Due to the lack of a good logging facility, many people start to develop their own versions. For example there are dozens of logging repositories for Golang. I run into the same problem when I am writing [```rafted```][rafted-github]. A powerful logging facility is needed to develop and debug it. I take a search on a few existing logging libraries for Golang but none of them seems to meet the requirement. So I decide to join the parade of "everyone is busy developing his own version", and then this library is created.

## Features

With an obivious intention to be a port of ```logging``` for Golang, ```go-logging``` has all the main features that ```logging``` package has:

1. It supports logging level, logging sources(Logger) and destinations(Handler) customization and flexible combinations of them
2. It supports Logger hierarchy
3. Optional Filter on Logger and Handler
4. It could be configured throught handy config file in various format
5. It support handlers that frequently-used in most real production enviroments, e.g. it could write log events to file, socket, rpc(the corresponding servers are also provided as bundled) etc.

Please note that 4, 5 are under development and not fully ready at this moment.

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
     logger.Warn("test message")
}
```

For more examples please refer to the test cases in source code.

## Documentation

For docs, run:

```bash
go doc github.com/hhkbp2/go-logging
```

For much more details please refer to the documentation for [```logging```][python-logging-page].

[python-logging-page]: https://docs.python.org/2/library/logging.html

[log4j-page]: http://logging.apache.org/log4j/

[golang-log-page]: http://golang.org/pkg/log/

[rafted-github]: http://github.com/hhkbp2/rafted

