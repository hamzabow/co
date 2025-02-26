package main

import (
	"log"
	"os"

	"github.com/hamzabow/co/internal/apikeyinput"
	"github.com/hamzabow/co/internal/genmessage"
	"github.com/hamzabow/co/internal/messagetextarea"
)

func main() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		var err error
		key, err = apikeyinput.PromptApiKeyWithRetries()
		if err != nil {
			log.Fatalf("Error message: %v", err)
		}
	}
	// from now on, key is not empty
	response, err := genmessage.GenerateCommitMessage(key)
	if err != nil {
		log.Fatalf("Error message: %v", err)
	}
	messagetextarea.MessageTextArea(response)
}
