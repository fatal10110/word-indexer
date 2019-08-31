package main

import (
	"os"
	"log"
)

const workersCount = 5
const webPort = 8080
var StatsStore Store
var ChannelBroker Broker

func main() {
	log.SetOutput(os.Stdout)
	
	StatsStore = NewMapStore()
	ChannelBroker = NewChannelBroker()

	ChannelBroker.Connect()
	defer ChannelBroker.Disconnect()

	for i := 0; i < workersCount; i++ {
		go ChannelBroker.Handle(CreateWorker())
	}

	StartServer(webPort, ChannelBroker)
	
}
