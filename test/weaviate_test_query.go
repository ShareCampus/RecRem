package main

import (
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func main() {
	weaviateUrl := os.Getenv("WEAVIATE_URL")
	weaviateKey := os.Getenv("WEAVIATE_API_KEY")
	openai_key := os.Getenv("OPENAI_KEY")
	fmt.Println(weaviateKey, weaviateUrl)
	cfg := weaviate.Config{
		Host:   weaviateUrl, // Replace with your Weaviate endpoint
		Scheme: "http",
		// AuthConfig: auth.ApiKey{Value: weaviateKey}, // Replace with your Weaviate instance API key
		Headers: map[string]string{
			"X-OpenAI-Api-Key": openai_key, // Replace with your inference API key
		},
		// Headers: nil,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"biology"})

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", result)
}
