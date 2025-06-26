package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

var logChan = make(chan DataChannelMessage, 100)

func initLogging() {

	go func() {
		file, err := os.Create("datachannel_messages.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header once
		writer.Write([]string{"sep=,"}) // This line is to ensure Excel opens the file correctly
		writer.Write([]string{
			"FrameID", "MessageSentTimeLocalMachine1", "MessageSentTimeVM1", "MessageSentTimeVM2",
			"MessageSentTimeLocalMachine2", "MessageSendRate", "LatencyEndToEnd",
		})

		for msg := range logChan {
			// write to csv
			record := []string{
				strconv.FormatInt(msg.FrameID, 10),
				strconv.FormatInt(msg.MessageSentTimeLocalMachine1, 10),
				strconv.FormatInt(msg.MessageSentTimeVM1, 10),
				strconv.FormatInt(msg.MessageSentTimeVM2, 10),
				strconv.FormatInt(msg.MessageSentTimeLocalMachine2, 10),
				strconv.FormatInt(msg.MessageSendRate, 10),
				strconv.FormatInt(msg.LatencyEndToEnd, 10),
			}
			writer.Write(record)
			writer.Flush() // Optional: remove or batch flush for better performance
		}
	}()
}
