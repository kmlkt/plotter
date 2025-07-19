package plotter

import (
	"fmt"
	"testing"
)

func TestDStringInt(t *testing.T) {
	testDString(0*d1, "0", t)
	testDString(1*d1, "1", t)
	testDString(100500*d1, "100500", t)
	testDString(-1*d1, "-1", t)
	testDString(-100500*d1, "-100500", t)
}

func TestDStringFloat(t *testing.T) {
	testDString(1*d1/100, "0.01", t)
	testDString(1005*d1/10, "100.5", t)
	testDString(-1*d1/100, "-0.01", t)
	testDString(-1005*d1/10, "-100.5", t)
}

func TestDParseInt(t *testing.T) {
	testDParse(1*d1, "1", t)
	testDParse(0*d1, "0", t)
	testDParse(100500*d1, "100500", t)
	testDParse(-1*d1, "-1", t)
	testDParse(-100500*d1, "-100500", t)
}

func TestDParseFloat(t *testing.T) {
	testDParse(1*d1/100, "0.01", t)
	testDParse(1005*d1/10, "100.5", t)
	testDParse(-1*d1/100, "-0.01", t)
	testDParse(-1005*d1/10, "-100.5", t)
}

func testDString(arg decimal, expected string, t *testing.T) {
	assertEquals(fmt.Sprintf("decimal of %v .String()", int64(arg)), arg.String(), expected, t)
}

func testDParse(expected decimal, arg string, t *testing.T) {
	ans, err := parseDecimal(arg)
	if err != nil {
		t.Fatal(err)
	}
	assertEquals(fmt.Sprintf("parseDecimal( %v )", arg), ans, expected, t)
}

func assertEquals[T comparable](key string, given T, expected T, t *testing.T) {
	if given != expected {
		t.Errorf("%s was %v, want %v", key, given, expected)
	}
}
