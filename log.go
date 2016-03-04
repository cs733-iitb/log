package log

//
// log.Log is an array on disk. It can be appended or truncated (all action is at the end).
// It keeps track of the last index to which data is written (-1 if the log is empty)
//

import (
	_ "fmt"
	"github.com/syndtr/goleveldb/leveldb"
	_ "github.com/syndtr/goleveldb/leveldb/storage"
	"strconv"
	"sync"
)

type Log struct {
	sync.Mutex
	ldb       *leveldb.DB
	lastIndex int64 // The highest index that has data.-1 if log is empty
}

// Create an indexed log at dbpath, or open one if it exists
func Open(dbpath string) (*Log, error) {
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return nil, err
	}

	var lg = &Log{ldb: db, lastIndex: -1}
	lastIndBytes, err := db.Get([]byte("lastIndex"), nil)
	if err == nil {
		lg.lastIndex = toIndex(lastIndBytes)
	}
	return lg, nil
}

// Append at the next available index, which is log.GetLastIndex()+1
func (lg *Log) Append(data []byte) error {
	lg.setLastIndex(lg.lastIndex + 1)

	tr, err := lg.ldb.OpenTransaction()
	ibytes := toBytes(lg.lastIndex)
	if err := tr.Put(ibytes, data, nil); err != nil {
		return err
	}
	if err := tr.Put([]byte("lastIndex"), ibytes, nil); err != nil {
		return err
	}
	err = tr.Commit()

	return err
}

// Remove all entries from (and including) fromIndex to the end
// GetLastIndex() will return fromIndex - 1
func (lg *Log) TruncateToEnd(fromIndex int64) error {
	tr, err := lg.ldb.OpenTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tr.Commit()
		} else {
			tr.Discard()
		}
	}()

	for i := fromIndex; i <= lg.lastIndex; i++ {
		err = tr.Delete(toBytes(i), nil)
		if err != nil {
			return err
		}
	}

	ibytes := toBytes(fromIndex - 1)
	if err := tr.Put([]byte("lastIndex"), ibytes, nil); err != nil {
		return err
	}
	lg.setLastIndex(fromIndex - 1)
	return nil
}

// Get data at index
func (lg *Log) Get(index int64) ([]byte, error) {
	return lg.ldb.Get(toBytes(index), nil)
}

func (lg *Log) Close() {
	if lg.ldb != nil {
		lg.ldb.Close()
	}
}

func (lg *Log) setLastIndex(i int64) {
	lg.Lock()
	lg.lastIndex = i
	lg.Unlock()
}

func (lg *Log) GetLastIndex() int64 {
	lg.Lock()
	i := lg.lastIndex
	lg.Unlock()
	return i
}

// int64 -> []byte
func toBytes(i int64) []byte {
	return []byte(strconv.FormatInt(i, 10)) // base 36
}

// []byte -> int64
func toIndex(key []byte) int64 {
	i, err := strconv.ParseInt(string(key), 10, 64) // base 36
	if err != nil {
		panic(err)
	}
	return i
}
