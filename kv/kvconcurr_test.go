package main

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"stdlog"
)

const NUM_OF_THREADS = 20

func concRoutine(t *testing.T, startSignal chan int, checkSignal chan int, wg1 *sync.WaitGroup, wg2 *sync.WaitGroup, myid int) {

	<-startSignal

	stdlog.Debugf("Thread %v has received start signal. starting...\n", myid)


	// create value as uuid
	u, err := uuid.NewV4()

	if err != nil {
		t.Logf("Error in uuid creation: %q", err)
		t.FailNow()
	}

	val := u.String()
	key := fmt.Sprintf("testkey_%v", myid)

	cliArgs := []string{fmt.Sprintf("%v=%v", key,val)}
	main2(cliArgs)

	wg1.Done()

	// wait for check signal
	stdlog.Debugf("Thread %v made its change. Waiting for check signal now\n", myid)
	<-checkSignal

	// get value from store and check it
	kvs := createStore()

	a := assert.New(t)
	if !a.EqualValues(kvs.Get(key), val, fmt.Sprintf("Thread %v: values not equal", myid)) {
		wg2.Done()
		t.FailNow()
	}

	wg2.Done()
}

func Test_ConcurrentAccess(t *testing.T) {
	startChan := make(chan int)
	checkChan := make(chan int)

	waitGroup1 := &sync.WaitGroup{}
	waitGroup2 := &sync.WaitGroup{}

	for i := 0; i < NUM_OF_THREADS; i++ {
		go concRoutine(t, startChan, checkChan, waitGroup1, waitGroup2, i)
		waitGroup1.Add(1)
		waitGroup2.Add(1)
	}

	stdlog.Debugf("%v Threads started. Sending start Signal!\n", NUM_OF_THREADS)
	close(startChan)

	waitGroup1.Wait()

	stdlog.Debug("All Threads returned - sending check signal!")
	close(checkChan)
	waitGroup2.Wait()

	stdlog.Debugf("All checks completed")
}
