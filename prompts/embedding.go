package prompts

import "encoding/json"

type EmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"`
}

func (e *EmbeddingRequest) Decode() []byte {
	data, _ := json.Marshal(e)
	return data
}
