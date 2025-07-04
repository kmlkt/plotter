package pltt

import (
	"encoding/binary"
	"iter"
	"os"
	"sync"
	"time"
)

type ConcurrentFile struct {
	*os.File
	sync.RWMutex
}

var opened = make(map[string]*ConcurrentFile)

func open(key string) (*ConcurrentFile, error) {
	file, exists := opened[key]
	if exists {
		return file, nil
	}
	osfile, err := os.OpenFile("./"+key, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	// NB! we are not closing the file
	if err != nil {
		return nil, err
	}
	opened[key] = &ConcurrentFile{osfile, sync.RWMutex{}}
	return opened[key], nil
}

func write(key string, value float64) error {
	file, err := open(key)
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
	Value     float64
}

func read(key string) (iter.Seq[record], error) {
	file, err := open(key)
	if err != nil {
		return nil, err
	}

	return func(yield func(record) bool) {
		file.Lock()
		defer file.Unlock()
		for i := int64(1); true; i++ {
			file.Seek(i*-16, 2)
			var timestamp int64
			var value float64
			err = binary.Read(file, binary.BigEndian, &timestamp)
			if err != nil {
				return
			}
			err = binary.Read(file, binary.BigEndian, &value)
			if err != nil {
				return
			}
			if !yield(record{time.UnixMicro(timestamp), value}) {
				return
			}
		}
	}, nil
}

func since(s iter.Seq[record], t time.Time) iter.Seq[record] {
	return func(yield func(record) bool) {
		for r := range s {
			if r.Timestamp.Before(t) {
				return
			}
			if !yield(r) {
				return
			}
		}
	}
}
