package prompts

import (
	"encoding/json"

	"github.com/bytedance/gopkg/util/logger"
)

const (
	DefaultMaxTokenSize = 1200
)

type Role string

const (
	RoleEmpty     Role = ""
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

type ResponseFormatType string

const (
	ResponseFormatJSON ResponseFormatType = "json_object"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Prompts struct {
	Stream         bool            `json:"stream,omitempty"`
	Model          string          `json:"model,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type ResponseFormatType `json:"type,omitempty"`
}

func New() *Prompts {
	return &Prompts{
		MaxTokens: DefaultMaxTokenSize,
		Stream:    false,
	}
}

func (p *Prompts) SetResponseFormat(content ResponseFormatType) {
	p.ResponseFormat = &ResponseFormat{
		Type: content,
	}
}

func (p *Prompts) AddSingleMessage(role Role, content string) {
	if content == "" {
		logger.Debugf("empty content")
		return
	}

	p.Messages = append(p.Messages, Message{
		Role:    role,
		Content: content,
	})
}

func (p *Prompts) AddSystemPrompt(content string) {
	if len(p.Messages) > 0 && p.Messages[0].Role == RoleSystem {
		p.Messages = p.Messages[1:]
	}

	systemMessage := Message{
		Role:    RoleSystem,
		Content: content,
	}

	p.Messages = append([]Message{systemMessage}, p.Messages...)
}

func (p *Prompts) String() string {
	return string(p.Decode())
}

func (p *Prompts) Decode() []byte {
	data, _ := json.Marshal(p)
	return data
}
