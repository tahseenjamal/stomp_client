package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-stomp/stomp"
)

var conn *stomp.Conn
var err error

type packet struct {
	queue   string
	message string
}

var brokerQueue chan packet

func getSession() *stomp.Conn {

	for {

		fmt.Println("Attempting connection")
		// ActiveMQ broker details
		broker := "localhost"
		port := "61613"
		username := ""
		password := ""
		// Create connection with Reconnect option and HeartBeatGracePeriodMultiplier
		conn_local, err := stomp.Dial("tcp", broker+":"+port,
			stomp.ConnOpt.Login(username, password),
			stomp.ConnOpt.AcceptVersion(stomp.V12),
			stomp.ConnOpt.HeartBeat(5*time.Second, 5*time.Second), // Set the initial heartbeats
			stomp.ConnOpt.HeartBeatGracePeriodMultiplier(5),       // Set HeartBeatGracePeriodMultiplier
		)

		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
		} else {

			fmt.Println("connected")
			return conn_local
		}

	}

}

func sendMessage() {

	for {

		data := <-brokerQueue
		queue := data.queue
		message := data.message

		// Will break out of this loop on successful connection and
		// message delivery. So messages are not lost
		for {

			err = conn.Send(queue, "text/plain", []byte(message))
			if err != nil {
				log.Printf("Error sending message to %s: %s", queue, err)
				conn = getSession()
				fmt.Println("Inside send message connected")
			} else {
				fmt.Printf("Message sent to %s: %s\n", queue, message)
				break
			}
		}
	}
}

func init() {

	brokerQueue = make(chan packet)
	conn = getSession()
	go sendMessage()
}

func main() {

	// Define your queues
	queue := "/queue/queue1"

	for i := 1; i <= 20; i++ {

		// Message content
		message := fmt.Sprintf("Message number %d", i)
		brokerQueue <- packet{queue, message}
		time.Sleep(time.Second * 2)
	}
}
