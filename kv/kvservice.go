package main

import "stdlog"

type KvsCommand struct {
	setKey      string
	setVal      string
	printFilter []string
	backChan chan string
}

func startService(kvs *kvstore) chan KvsCommand {
	cc := make(chan KvsCommand)
	go processLoop(kvs, cc)
	return cc
}

func processLoop(kvs *kvstore, commandChannel chan KvsCommand) {
	for {
		cmd := <-commandChannel
		stdlog.Debugf("Command received: %q", cmd)
		if cmd.setKey != "" {
			kvs.Set(cmd.setKey, cmd.setVal)
			// always persist store after write

			if err := kvs.PersistStore(); err != nil {
				stdlog.Fatalf("Cannot persist store:", err)
			}


			cmd.backChan <- "ok" // todo: return result

		} else if len(cmd.printFilter) > 0 {
			kvs.Printkeys(cmd.printFilter)
			cmd.backChan <- "ok" // todo: return result
		} else {
			kvs.Printall()
			cmd.backChan <- "ok" // todo: return result
		}


	}
}


