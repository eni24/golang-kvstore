package main

import (
	"fmt"
	"os"
	"strings"
	"stdlog"
)

var cmdChan chan KvsCommand

func init() {
	cmdChan = startService(createStore())
}

func main() {
	main2(os.Args[1:])
}

func main2(args []string) {
	onearg := strings.Join(args, "")

	backChan := make(chan string)

	// is this an assignment k=v?
	if splitarg := strings.Split(onearg, "="); len(splitarg) == 2 {
		cmdChan <- KvsCommand{splitarg[0], splitarg[1], nil, backChan}
		resp := <- backChan
		stdlog.Debugf("Response in backchannel is: %q", resp)
		return
	}

	// print kv-pairs. any args are interpreted as key filters
	cmdChan <- KvsCommand{"", "", args, backChan}
	resp := <- backChan
	stdlog.Debugf("Response in backchannel is: %q", resp)
}

func createStore() *kvstore {
	var f *os.File
	var err error

	if f, err = getDbfile(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	//defer f.Close()

	return NewFileboundKvstore(f)
}

func getDbfile() (*os.File, error) {
	dirname := fmt.Sprintf("%v/.kv", os.Getenv("HOME"))

	// check dir, create if necessary
	if dir, err := os.Open(dirname); err == nil {
		dir.Close() // check done, close file
	} else {
		// create dir
		if err = os.Mkdir(dirname, 0700); err != nil {
			return nil, fmt.Errorf("unable to create dir '%v'\n", dirname)
		}
	}

	// check file, create if necessary
	filename := fmt.Sprintf("%v/kvdb.txt", dirname)
	var kvfile *os.File
	var err error

	if kvfile, err = os.OpenFile(filename, os.O_RDWR, 0700); err != nil {
		if kvfile, err = os.Create(filename); err != nil {
			return nil, fmt.Errorf("unable to create file '%v'\n", filename)
		}
	}

	// all good
	return kvfile, nil
}
