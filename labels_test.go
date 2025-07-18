package pltt

import (
	"fmt"
	"testing"
)

func TestValueStringInt(t *testing.T) {
	testValueString(0*e9, "0", t)
	testValueString(1*e9, "1", t)
	testValueString(100500*e9, "100500", t)
	testValueString(-1*e9, "-1", t)
	testValueString(-100500*e9, "-100500", t)
}

func TestValueStringFloat(t *testing.T) {
	testValueString(1*e9/100, "0.01", t)
	testValueString(1005*e9/10, "100.5", t)
	testValueString(-1*e9/100, "-0.01", t)
	testValueString(-1005*e9/10, "-100.5", t)
}

func testValueString(arg int64, expected string, t *testing.T) {
	assertEquals(fmt.Sprintf("valueString of %v", arg), valueString(0, arg), expected, t)
}

func assertEquals[T comparable](key string, given T, expected T, t *testing.T) {
	if given != expected {
		t.Errorf("%s was %v, want %v", key, given, expected)
	}
}
