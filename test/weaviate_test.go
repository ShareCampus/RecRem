package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
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

	// Check the connection
	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库链接成功", live)

	// // add the schema
	// classObj := &models.Class{
	// 	Class:      "Question1",
	// 	Vectorizer: "text2vec-openai",
	// 	ModuleConfig: map[string]interface{}{
	// 		"generative-openai": map[string]interface{}{},
	// 	},
	// }

	// if err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); err != nil {
	// 	panic(err)
	// }

	// Retrieve the data
	items, err := getJSONdata()
	if err != nil {
		panic(err)
	}

	// convert items into a slice of models.Object
	objects := make([]*models.Object, len(items))
	for i := range items {
		objects[i] = &models.Object{
			Class: "Question1",
			Properties: map[string]any{
				"category": items[i]["Category"],
				"question": items[i]["Question"],
				"answer":   items[i]["Answer"],
			},
		}
	}

	// batch write items
	batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}
	for _, res := range batchRes {
		if res.Result.Errors != nil {
			panic(res.Result.Errors.Error)
		}
	}
}

func getJSONdata() ([]map[string]string, error) {
	// Retrieve the data
	data, err := http.DefaultClient.Get("https://raw.githubusercontent.com/weaviate-tutorials/quickstart/main/data/jeopardy_tiny.json")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()

	// Decode the data
	var items []map[string]string
	if err := json.NewDecoder(data.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}
