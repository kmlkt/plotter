package plotter

import (
	"encoding/binary"
	"iter"
	"os"
	"path"
	"sync"
	"time"
)

func init() {
	os.MkdirAll(dirData, 0750)
}

type ConcurrentFile struct {
	*os.File
	sync.RWMutex
}

var opened = make(map[string]*ConcurrentFile)

func open(tableName string) (*ConcurrentFile, error) {
	file, exists := opened[tableName]
	if exists {
		return file, nil
	}
	osfile, err := os.OpenFile(path.Join(dirData, tableName), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	// NB! we are not closing the file
	if err != nil {
		return nil, err
	}
	opened[tableName] = &ConcurrentFile{osfile, sync.RWMutex{}}
	return opened[tableName], nil
}

func write(key string, value decimal) error {
	tn, err := readOrCreateTableName(key, permWrite)
	if err != nil {
		return err
	}
	file, err := open(tn)
	file.Lock()
	defer file.Unlock()
	if err != nil {
		return err
	}
	timestamp := time.Now().UnixMicro()
	err = binary.Write(file, binary.BigEndian, timestamp)
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, value)
	if err != nil {
		return err
	}
	return nil
}

type record struct {
	Timestamp time.Time
	Value     decimal
}

func read(key string) (iter.Seq[record], error) {
	tn, err := readOrCreateTableName(key, permRead)
	if err != nil {
		return nil, err
	}
	file, err := open(tn)
	if err != nil {
		return nil, err
	}

	return func(yield func(record) bool) {
		file.Lock()
		defer file.Unlock()
		for i := int64(1); true; i++ {
			file.Seek(i*-16, 2)
			var timestamp int64
			var value decimal
			err = binary.Read(file, binary.BigEndian, &timestamp)
			if err != nil {
				return
			}
			err = binary.Read(file, binary.BigEndian, &value)
			if err != nil {
				return
			}
			if !yield(record{time.UnixMicro(timestamp).UTC(), value}) {
				return
			}
		}
	}, nil
}
