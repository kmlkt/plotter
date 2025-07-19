package plotter

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path"
)

func init() {
	os.MkdirAll(dirRead, 0750)
	os.MkdirAll(dirWrite, 0750)
}

type permission byte

const (
	permRead permission = iota
	permWrite
)

func (p permission) String() string {
	if p == permRead {
		return "read"
	} else {
		return "write"
	}
}

func readOrCreateTableName(key string, perm permission) (string, error) {
	tn, err := readTableName(key, perm)
	if errors.Is(err, errorKeyNotFound) {
		tn, err = createKeyPair(key, key)
		if errors.Is(err, errorKeyExists) {
			return "", formatError(errorKeyNoPermission, key, perm)
		}
	}
	return tn, err
}

func readTableName(key string, perm permission) (string, error) {
	bytes, err := os.ReadFile(keyFile(key, perm))
	if err != nil {
		if os.IsNotExist(err) {
			return "", formatError(errorKeyNotFound, perm, key)
		}
		return "", err
	}
	return string(bytes), err
}

func generateKeyPair() (string, string, error) {
	for {
		readKey := rand.Text()
		writeKey := rand.Text()
		_, err := createKeyPair(readKey, writeKey)
		if errors.Is(err, errorKeyExists) {
			continue
		}
		if err != nil {
			return "", "", err
		}
		return readKey, writeKey, nil
	}
}

func createKeyPair(readKey string, writeKey string) (string, error) {
	checkKey := func(key string, perm permission) error {
		_, err := readTableName(key, perm)
		return expectError(err, errorKeyNotFound,
			formatError(errorKeyExists, perm, key))
	}
	err := errors.Join(
		checkKey(readKey, permRead),
		checkKey(writeKey, permWrite),
	)
	if err != nil {
		return "", err
	}
	tn := tableName(readKey, writeKey)
	return tn, errors.Join(
		createPermission(readKey, permRead, tn),
		createPermission(writeKey, permWrite, tn))
}

func createPermission(key string, perm permission, tableName string) error {
	return os.WriteFile(keyFile(key, perm), []byte(tableName), 0666)
}

func keyFile(key string, perm permission) string {
	return path.Join(permLookupDir(perm), key)
}

func tableName(readKey string, writeKey string) string {
	return fmt.Sprintf("%s %s", readKey, writeKey)
}

// not suitable for permBoth
func permLookupDir(perm permission) string {
	switch perm {
	case permRead:
		return dirRead
	case permWrite:
		return dirWrite
	default:
		return ""
	}
}
