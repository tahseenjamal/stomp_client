package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-stomp/stomp"
)

func main() {
	// ActiveMQ broker details
	broker := "localhost"
	port := "61613"
	username := ""
	password := ""

	// Create connection with Reconnect option and HeartBeatGracePeriodMultiplier
	conn, err := stomp.Dial("tcp", broker+":"+port,
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.HeartBeat(5*time.Second, 5*time.Second), // Set the initial heartbeats
		stomp.ConnOpt.HeartBeatGracePeriodMultiplier(5),       // Set HeartBeatGracePeriodMultiplier
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Disconnect()

	// Define your queues
	queue1 := "/queue/queue1"

	// Message content
	message := "Hello, ActiveMQ!"

	// Send message to both queues
	err = conn.Send(queue1, "text/plain", []byte(message))
	if err != nil {
		log.Printf("Error sending message to %s: %s", queue1, err)
	} else {
		fmt.Printf("Message sent to %s: %s\n", queue1, message)
	}

	time.Sleep(time.Second * 2)
}
