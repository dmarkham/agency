package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dmarkham/agency/agent"
	"github.com/dmarkham/agency/client/ollama"
	log "github.com/sirupsen/logrus"
)

func main() {
	//client := openai.New(os.Getenv("OPENAI_API_KEY"))
	client := ollama.New("http://"+os.Getenv("OLLAMA_HOST"), http.DefaultClient)

	// Initialize a poet agent and a critic agent
	poet := agent.New("poet",
		agent.WithClient(client),
		//agent.WithModel("dolphin-llama3:8b-256k-v2.9-q4_K_M"),
		agent.WithModel("llama3:8b"),
		agent.WithMaxTokens(2000),
		agent.WithTemperature(0.9))

	critic := agent.New("critic",
		agent.WithClient(client),
		agent.WithModel("llama3:8b"),
		agent.WithTemperature(0.9),
		agent.WithMaxTokens(2000))

	// Set the topic for the haiku
	topic := "golang"

	// The poet writes a haiku about the given topic
	_, err := poet.Listen(fmt.Sprintf("Write a haiku about a %s", topic))
	if err != nil {
		log.Fatalf("Poet Listen failed: %v", err)
	}

	haiku, err := poet.Respond(context.Background())
	if err != nil {
		log.Fatalf("Poet Respond failed: %v", err)
	}

	fmt.Printf("Haiku:\n%s\n\n", haiku)

	// The critic critiques the haiku
	_, err = critic.Listen(fmt.Sprintf("Please provide harsh critique's of this haiku: \n%s", haiku))
	if err != nil {
		log.Fatalf("Critic Listen failed: %v", err)
	}
  poet.Append()
	critique, err := critic.Respond(context.Background())
	if err != nil {
		log.Fatalf("Critic Respond failed: %v", err)
	}

	fmt.Printf("Critique: %s\n\n", critique)

	// Create a loop for the poet to improve and the critic to critique
	for i := 0; i < 10; i++ {
		// Poet takes into account the critique and tries to improve
		_, err = poet.Listen(fmt.Sprintf("Feedback received: '%s'. Please improve the haiku. only respond with the haiku.", critique))
		if err != nil {
			log.Fatalf("Poet Listen failed: %v", err)
		}

		haiku, err = poet.Respond(context.Background())
		if err != nil {
			log.Fatalf("Poet Respond failed: %v", err)
		}
		fmt.Printf("\nImproved Haiku:\n %s\n\n", haiku)

		// The critic critiques the improved haiku
		_, err = critic.Listen(fmt.Sprintf("Please provide harsh critique's of this haiku: \n%s", haiku))
		if err != nil {
			log.Fatalf("Critic Listen failed: %v", err)
		}

		critique, err = critic.Respond(context.Background())
		if err != nil {
			log.Fatalf("Critic Respond failed: %v", err)
		}

		fmt.Printf("Critique: %s\n\n", critique)
	}
}
