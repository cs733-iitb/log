#log

log.Log is an array on disk. It can be appended or truncated (all action is at the end).
It keeps track of the last index to which data is written (-1 if the log is empty)

# Usage

```go
     import "github.com/cs733-iitb/log"

     lg := log.Open("mylog")
     defer lg.Close()

     lg.Append([]byte("foo"))
     lg.Append([]byte("bar"))
     lg.Append([]byte("baz"))

     bytes = log.Get(1) // should return "bar" in bytes
     i := log.GetLastIndex() // should return 2 as an int64 value

     log.TruncateToEnd(/*from*/ 1)
     i := log.GetLastIndex() // should return 0. One entry is left.

```

# Installation and Dependencies.

    go get github.com/cs733-iit/log
    go test -race github.com/cs733-iit/log

This library depends on the github.com/syndtr/leveldb

# Bugs

This library is a quick and dirty and inefficient solution. Treat it like a toy.

##Author: Sriram Srinivasan. sriram _at_ malhar.net
