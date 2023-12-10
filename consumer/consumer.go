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

	// Define your queue
	queue := "/queue/queue1"

	// Establish initial connection
	conn := connect(broker, port, username, password)
	defer conn.Disconnect()

	// Start listener
	go listen(conn, queue)

	// Send and receive messages
	for {
		time.Sleep(time.Second * 2)
	}
}

func connect(broker, port, username, password string) *stomp.Conn {
	for {
		conn, err := stomp.Dial("tcp", broker+":"+port,
			stomp.ConnOpt.Login(username, password),
			stomp.ConnOpt.AcceptVersion(stomp.V12),
			stomp.ConnOpt.HeartBeat(5*time.Second, 5*time.Second),
			stomp.ConnOpt.HeartBeatGracePeriodMultiplier(5),
		)
		if err == nil {
			fmt.Println("Connected to ActiveMQ")
			return conn
		}

		fmt.Printf("Error connecting to ActiveMQ: %s. Retrying...\n", err)
		time.Sleep(time.Second * 5)
	}
}

func listen(conn *stomp.Conn, queue string) {
	sub, err := conn.Subscribe(queue, stomp.AckAuto)
	if err != nil {
		log.Fatalf("Error subscribing to %s: %s", queue, err)
	}
	defer sub.Unsubscribe()

	for {
		msg := <-sub.C
		if msg.Err != nil {
			log.Printf("Error receiving message from %s: %s", queue, msg.Err)
			time.Sleep(time.Second * 5) // Wait before reconnecting
			conn = connect("localhost", "61613", "your_username", "your_password")
			sub, err = conn.Subscribe(queue, stomp.AckAuto)
			if err != nil {
				log.Fatalf("Error re-subscribing to %s: %s", queue, err)
			}
			continue
		}

		fmt.Printf("Received message from %s: %s\n", queue, msg.Body)
	}
}
