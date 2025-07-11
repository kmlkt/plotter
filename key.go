package pltt

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
)

type permission byte

const (
	permRead permission = iota
	permWrite
)

func readOrCreateTableName(key string, perm permission) (*string, error) {
	tn, err := readTableName(key, perm)
	if tn != nil || err != nil {
		return tn, err
	}

	withAnotherPerm, err := readTableName(key, 1-perm)
	if withAnotherPerm != nil || err != nil {
		return nil, err
	}

	newTn, err := createKeyPair(key, key)
	if err != nil {
		return nil, err
	}
	return &newTn, nil
}

func readTableName(key string, perm permission) (*string, error) {
	bytes, err := os.ReadFile(keyFile(key, perm))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	s := string(bytes)
	return &s, err
}

func generateKeyPair() (string, string, error) {
	for {
		readKey := rand.Text()
		writeKey := rand.Text()
		tn1, err := readTableName(readKey, permRead)
		tn2, err := readTableName(writeKey, permWrite)
		if  != nil ||
		 != nil {

		}
	}
}

func createKeyPair(readKey string, writeKey string) (string, error) {
	tn := tableName(readKey, writeKey)
	err := createPermission(readKey, permRead, tn)
	if err != nil {
		return "", err
	}
	err = createPermission(writeKey, permWrite, tn)
	if err != nil {
		return "", err
	}
	return tn, nil
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
