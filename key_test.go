package pltt

import (
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
	if tn1 != *tn2 || tn1 != *tn3 {
		t.Errorf("Table names = %s, %s, %s, want three same", tn1, *tn2, *tn3)
	}
	tn4, err := readOrCreateTableName(A, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	tn5, err := readOrCreateTableName(B, permRead)
	if err != nil {
		t.Fatal(err)
	}
	if tn4 != nil || tn5 != nil {
		t.Error("Returned tablename when have no rights (want nil)")
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
	if tn1 != *tn2 || tn1 != *tn3 {
		t.Errorf("Table names = %s, %s, %s, want three same", tn1, *tn2, *tn3)
	}
}

func TestImlicitkeys(t *testing.T) {
	A := "kcuf"
	tn1, err := readOrCreateTableName(A, permRead)
	if err != nil {
		t.Fatal(err)
	}
	tn2, err := readOrCreateTableName(A, permWrite)
	if err != nil {
		t.Fatal(err)
	}
	if *tn1 != *tn2 {
		t.Errorf("Table names = %s, %s, want two same", *tn1, *tn2)
	}
}
