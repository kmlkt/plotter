package pltt

import (
	"os"
	"path"
)

var (
	dirData  string
	dirRead  string
	dirWrite string
)

func init() {
	cwd, _ := os.Getwd()
	setFromEnv(&cwd, "DIR")
	dirData = path.Join(cwd, "data")
	dirRead = path.Join(cwd, "read")
	dirWrite = path.Join(cwd, "write")
}

func setFromEnv(s *string, key string) {
	v, exists := os.LookupEnv(key)
	if exists {
		*s = v
	}
}
