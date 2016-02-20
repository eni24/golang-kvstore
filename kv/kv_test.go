package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSimpleSetValue(t *testing.T) {
	a := assert.New(t)

	kvs := NewKvstore()
	kvs.Set("Wurst", "Kaese")

	a.Equal(kvs.Get("Wurst"), "Kaese")

	//Change and check again
	kvs.Set("Wurst", "Brot")
	a.Equal(kvs.Get("Wurst"), "Brot")
}

func TestSetCrosswiseValues(t *testing.T) {
	a := assert.New(t)

	kvs := NewKvstore()
	kvs.Set("1", "2")
	kvs.Set("2", "1")
	kvs.Set("3", "1")
	kvs.Set("4", "2")

	a.Equal(kvs.Get("1"), "2")
	a.Equal(kvs.Get("2"), "1")
	a.Equal(kvs.Get("3"), "1")
	a.Equal(kvs.Get("4"), "2")
}

func TestFileRead(t *testing.T) {
	a := assert.New(t)
	var f *os.File
	var err error

	if f, err = os.OpenFile("kvtestdata.txt", os.O_RDONLY, 0400); err != nil {
		a.Fail("Error in Test: Can't read input file")
	}

	kvs := NewFileboundKvstore(f)
	a.Equal(kvs.Get("xyxyxyxyxyx"), "uewfhreiufhireh")
}
