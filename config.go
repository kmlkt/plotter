package plotter

import (
	"os"
	"path"
)

var (
	port     string = "80"
	dirData  string
	dirRead  string
	dirWrite string
)

func init() {
	setFromEnv(&port, "PORT")
	cwd, _ := os.Getwd()
	dirStorage := path.Join(cwd, "storage")
	setFromEnv(&dirStorage, "DIR")
	dirData = path.Join(dirStorage, "data")
	dirRead = path.Join(dirStorage, "read")
	dirWrite = path.Join(dirStorage, "write")
}

func setFromEnv(s *string, key string) {
	v, exists := os.LookupEnv(key)
	if exists {
		*s = v
	}
}
