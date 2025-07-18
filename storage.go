package pltt

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

func write(key string, value float64) error {
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
	Value     float64
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

func intervalSum(s iter.Seq[record], d time.Duration) iter.Seq[record] {
	return func(yield func(record) bool) {
		prev := int64(-1)
		ps := 0.0
		for r := range s {
			cur := r.Timestamp.UnixNano() / d.Nanoseconds()
			if prev == -1 {
				prev = cur
			}
			if cur == prev {
				ps += r.Value
			} else {
				for prev > cur {
					if !yield(record{time.UnixMicro((prev + 1) * d.Microseconds()), ps}) {
						return
					}
					if !yield(record{time.UnixMicro(prev * d.Microseconds()), ps}) {
						return
					}
					ps = 0
					prev--
				}
				ps = r.Value
			}
		}
		if ps != 0 {
			yield(record{time.UnixMicro((prev + 1) * d.Microseconds()), ps})
			yield(record{time.UnixMicro(prev * d.Microseconds()), ps})
		}
	}
}
