package main

import "os"

func main() {
	arguments := os.Args
	if arguments[1] == "MessagingNode" {
		runMessageNode(arguments[HOST], arguments[PORT], arguments[DESTINATION_HOST], arguments[DESTINATION_PORT], arguments[MESSAGE_COUNT])
	} else {
		runRegistry()
	}
}
