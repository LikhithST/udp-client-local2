package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

type DataChannelMessage struct {
	FrameID                      int64  `json:"frameID"`
	MessageSentTimeLocalMachine1 int64  `json:"messageSentTime_LocalMachine1,omitempty"`
	MessageSentTimeVM1           int64  `json:"messageSentTime_VM1,omitempty"`
	MessageSentTimeVM2           int64  `json:"messageSentTime_VM2,omitempty"`
	MessageSentTimeLocalMachine2 int64  `json:"messageSentTime_LocalMachine2,omitempty"`
	LatencyEndToEnd              int64  `json:"latency_end_to_end,omitempty"`
	MessageSendRate              int64  `json:"message_send_rate,omitempty"`
	Payload                      []byte `json:"payload"`
}

var msg DataChannelMessage

func main() {

	addrStr := flag.String("addr", ":9000", "UDP listen address")
	bufSize := flag.Int("bufsize", 4096, "Buffer size for incoming UDP packets")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *addrStr)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	initLogging()

	buf := make([]byte, *bufSize)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			fmt.Println("Invalid message from", remoteAddr, ":", err)
			continue
		}

		fmt.Printf("Received from %v: %+v\n", remoteAddr, msg)
		msg.MessageSentTimeLocalMachine2 = time.Now().UnixMilli()
		msg.LatencyEndToEnd = msg.MessageSentTimeLocalMachine2 - msg.MessageSentTimeLocalMachine1

		logChan <- msg
		fmt.Println("Message logged to channel")

	}
}
