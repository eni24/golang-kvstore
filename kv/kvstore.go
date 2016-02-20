package main

import (
	"bufio"
	"fmt"
	"os"
	"stdlog"
	"strings"
	//"sync"
)

type kvstore struct {
	kvmap map[string]string
	storageFile *os.File
}

func (kvs kvstore) Set(k string, v string) {
	stdlog.Debugf("setting key '%v' = '%v'", k, v)
	kvs.kvmap[k] = v
}

func (kvs kvstore) Get(k string) string {
	return kvs.kvmap[k]
}

func NewKvstore() *kvstore {
	x := new(kvstore)
	x.kvmap = make(map[string]string)
	x.storageFile = nil
	return x
}

func (kvs kvstore) Printkeys(filter []string) {
	for _, k := range filter {
		fmt.Printf("%v=%v\n", k, kvs.Get(k))
	}
}

func (kvs kvstore) Printall() {
	for k, _ := range kvs.kvmap {
		fmt.Printf("'%v'='%v'\n", k, kvs.Get(k))
	}
}

func NewFileboundKvstore(f *os.File) *kvstore {
	kvs := NewKvstore()
	kvs.storageFile = f
	defer kvs.storageFile.Close()

	reader := bufio.NewReader(kvs.storageFile)
	scanner := bufio.NewScanner(reader)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		stdlog.Debugf("Reading line from file: '%v'", line)

		x := strings.Split(line, "=")

		if len(x) != 2 {
			stdlog.Errorf("error in dbfile line #%v. Invalid number of elements: '%v'", i+1, line)
			os.Exit(1)
		}

		kvs.Set(x[0], x[1])
	}

	return kvs
}

func (kvs kvstore) PersistStore() error {

	stdlog.Debug("kvstore.PersistStore called")
	if kvs.storageFile == nil {
		return fmt.Errorf("Not a filebound kvstore. can't persist. Me=%q", kvs)
	}

	var err error
	kvs.storageFile, err = os.OpenFile(kvs.storageFile.Name(),os.O_RDWR, 0700)
	if err != nil {
		return fmt.Errorf("Error opening file %q for write access: %q", kvs.storageFile.Name(), err)
	}

	//stdlog.Debugf("kvs persist: %q", kvs)
	for k := range kvs.kvmap {
		line := fmt.Sprintf("%v=%v\n", k, kvs.Get(k))
		_, err = kvs.storageFile.WriteString(line)

		if err != nil {
			return fmt.Errorf("Error writing line '%v'. Check dbfile manually. error is: %q\n", line, err)
		}
	}

	return nil
}

