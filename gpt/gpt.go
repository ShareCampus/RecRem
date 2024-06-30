package gpt

import (
	"io"
	"net/http"
	"recrem/prompts"
)

type GPT interface {
	String() string
	Init()
	GetToken() (string, error)
	CallAPI(*prompts.Prompts, string) (*http.Response, error)
	ParseResponse(io.Reader) (int, string, error)
}
