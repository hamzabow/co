package main

import (
	"co/internal/apikeyinput"
	"co/internal/genmessage"
	"co/internal/messagetextarea"
	"errors"
	"fmt"
	"log"
	"os"
)

var ErrMissingAPIKey = errors.New("missing OPENAI_API_KEY environment variable")

func main() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		// log.Fatalf("Error message: %v", ErrMissingAPIKey)
		key = apikeyinput.PromptApiKey()
		if key == "" {
			max_retries := 3
			retries := max_retries
			for key == "" && retries > 0 {
				fmt.Println("-------------------------------------------------")
				fmt.Println("Key is empty! Please enter a valid OpenAI API KEY")
				fmt.Println("-------------------------------------------------")
				key = apikeyinput.PromptApiKey()
				retries -= 1
			}
			if retries <= 0 && key == "" {
				fmt.Println("---------------------------------------------------------------------")
				fmt.Printf("You have entered an invalid API key %d times. Please try again later.\n", max_retries)
				fmt.Println("---------------------------------------------------------------------")
				os.Exit(1)
			}
		}
	}
	// from now on, key is not empty
	response, err := genmessage.GenerateCommitMessage(key)
	if err != nil {
		log.Fatalf("Error message: %v", err)
	}
	messagetextarea.MessageTextArea(response)
}
