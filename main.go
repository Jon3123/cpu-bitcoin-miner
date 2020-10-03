package main

import (
	"cpu-bitcoin-miner/gobtcminer/client"
	"cpu-bitcoin-miner/gobtcminer/config"
	"cpu-bitcoin-miner/gobtcminer/logger"
	"cpu-bitcoin-miner/gobtcminer/mining"
	"time"
)

var monitor logger.Logger

func main() {

	// Read and parse the configuration file
	conf := config.ReadConf("config.json")
	monitor = logger.NewLogger(conf.Log)
	diff, err := client.GetDifficulty(conf.User, conf.Password, conf.Host)
	if err != nil {
		monitor.Print("info", "Error getting difficulty: "+err.Error())
	}
	dispatcher := mining.NewDispatcher(monitor)
	dispatcher.Run()
	//Add Chunks on a regular basis
	for {
		//Get a new Chunk and split it accordingly to the machin settings
		for _, chunk := range mining.NewChunkList(2, uint32(time.Now().Unix()), diff) {
			if len(dispatcher.ChunkQueueIn) < cap(dispatcher.ChunkQueueIn) {
				dispatcher.ChunkQueueIn <- chunk
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}