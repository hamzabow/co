package main

import (
	"co/internal/genmessage"
	"co/internal/messagetextarea"
	"log"
)

func main() {
	response, err := genmessage.GenerateCommitMessage()
	if err != nil {
		log.Fatalf("Error message: %v", err)
	}
	messagetextarea.MessageTextArea(response)
}
