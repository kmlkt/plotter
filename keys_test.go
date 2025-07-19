package plotter

import (
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.RemoveAll(dirData)
	os.RemoveAll(dirRead)
	os.RemoveAll(dirWrite)
	os.MkdirAll(dirData, 0750)
	os.MkdirAll(dirRead, 0750)
	os.MkdirAll(dirWrite, 0750)
	m.Run()
	os.RemoveAll(dirData)
	os.RemoveAll(dirRead)
	os.RemoveAll(dirWrite)
}

func TestDifferentKeys(t *testing.T) {
	A := "abacaba"
	B := "abcdef"
	tn1, err := createKeyPair(A, B)
	if err != nil {
		t.Fatal(err)
	}
	tn2, err := readOrCreateTableName(A, permRead)
	if err != nil {
		t.Fatal(err)
	}
	tn3, err := readOrCreateTableName(B, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	if tn1 != tn2 || tn1 != tn3 {
		t.Errorf("Table names = %s, %s, %s, want three same", tn1, tn2, tn3)
	}
	_, err = readOrCreateTableName(A, permWrite)
	if err == nil {
		t.Error("Did not return errorKeyNoPermission when has no rights")
	} else if !errors.Is(err, errorKeyNoPermission) {
		t.Fatal(err)
	}
	_, err = readOrCreateTableName(B, permRead)
	if err == nil {
		t.Error("Did not return errorKeyNoPermission when has no rights")
	} else if !errors.Is(err, errorKeyNoPermission) {
		t.Fatal(err)
	}
}

func TestSameKeys(t *testing.T) {
	A := "reggin"
	tn1, err := createKeyPair(A, A)
	if err != nil {
		t.Fatal(err)
	}
	tn2, err := readOrCreateTableName(A, permRead)
	if err != nil {
		t.Fatal(err)
	}
	tn3, err := readOrCreateTableName(A, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	if tn1 != tn2 || tn1 != tn3 {
		t.Errorf("Table names = %s, %s, %s, want three same", tn1, tn2, tn3)
	}
}

func TestImlicitKeys(t *testing.T) {
	A := "kcuf"
	tn1, err := readOrCreateTableName(A, permRead)
	if err != nil {
		t.Fatal(err)
	}
	tn2, err := readOrCreateTableName(A, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	if tn1 != tn2 {
		t.Errorf("Table names = %s, %s, want two same", tn1, tn2)
	}
}

func TestGeneratedKeys(t *testing.T) {
	A, B, err := generateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	tn2, err := readOrCreateTableName(A, permRead)
	if err != nil {
		t.Fatal(err)
	}
	tn3, err := readOrCreateTableName(B, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	if tn2 != tn3 {
		t.Errorf("Table names = %s, %s, want two same", tn2, tn3)
	}
	_, err = readOrCreateTableName(A, permWrite)
	if err == nil {
		t.Error("Did not return errorKeyNoPermission when has no rights")
	} else if !errors.Is(err, errorKeyNoPermission) {
		t.Fatal(err)
	}
	_, err = readOrCreateTableName(B, permRead)
	if err == nil {
		t.Error("Did not return errorKeyNoPermission when has no rights")
	} else if !errors.Is(err, errorKeyNoPermission) {
		t.Fatal(err)
	}
}
