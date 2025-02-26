package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
	commitMessage, commitResult := messagetextarea.MessageTextArea(response)

	if commitMessage == "" {
		fmt.Println("No commit message provided")
		return
	}

	if commitResult == messagetextarea.ResultCommit {
		commit(commitMessage)
	} else {
		fmt.Println("Commit cancelled")
		return
	}
}

func commit(msg string) {
	cmd := exec.Command("git", "commit", "-m", msg)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error message: %v", err)
	}
	fmt.Println("Commit successful")
}
