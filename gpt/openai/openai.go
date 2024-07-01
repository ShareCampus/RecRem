package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"recrem/config/db"
	"recrem/models"
	"recrem/utils"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
)

// open ai API calls
const (
	OPENAI                      = "openai"
	DataPrefix                  = "data: "
	DataEndStr                  = "DONE"
	DefaultChannelTimeoutSecond = 5
)

type OpenAI struct {
	accessKeyChan chan models.OpenAI
}

func NewOpenAI() *OpenAI {
	return &OpenAI{}
}

var O *OpenAI

// var _ gpt.GPT = &OpenAI{} // 类型断言

func InitGpt() {
	O = NewOpenAI()
	var accessKeys []models.OpenAI
	if err := db.Db.Model(models.OpenAI{}).Find(&accessKeys).Error; err != nil {
		panic(err)
	}
	O.accessKeyChan = make(chan models.OpenAI, len(accessKeys))
	for _, key := range accessKeys {
		O.accessKeyChan <- key
	}
}

func (o *OpenAI) CallEmbeddingAPI(prompt *models.EmbeddingRequest) (*http.Response, error) {
	token, err := o.GetToken()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, EmbeddingEndPoint, bytes.NewBuffer(prompt.Decode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("bad response status code %d", resp.StatusCode)
	}

	return resp, nil
}

func (o *OpenAI) GetToken() (string, error) {
	accessKey, err := getAccessKey()
	if err != nil {
		return "", err
	}

	return accessKey.Token, nil
}

func getAccessKey() (*models.OpenAI, error) {
	select {
	case key := <-O.accessKeyChan:
		O.accessKeyChan <- key
		return &key, nil
	case <-time.After(time.Second * DefaultChannelTimeoutSecond):
		return nil, errors.New("no access key available")
	}
}

func (o *OpenAI) ParseResponse(stream io.Reader) (int, string, error) {
	data, err := io.ReadAll(stream)
	if err != nil {
		return 0, "", err
	}
	var content []string
	model, responses := chatResponses(string(data))
	for _, item := range responses {
		for _, choice := range item.Choices {
			content = append(content, choice.Delta.Content, choice.Message.Content)
		}
	}

	result := strings.Join(content, "")
	return utils.NumTokensFromMessages(result, model), result, nil
}

func chatResponses(input string) (gptModel string, chatresponses []models.ChatResponse) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		jsonStr := strings.TrimPrefix(line, DataPrefix)
		if strings.Contains(jsonStr, DataEndStr) {
			break
		}

		var response models.ChatResponse
		if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
			logger.Errorf("marshal json error: %v", err)
			continue
		}
		if len(response.ID) == 0 {
			continue
		}

		if len(gptModel) == 0 && len(response.Model) > 0 {
			gptModel = response.Model
		}
		chatresponses = append(chatresponses, response)
	}
	return
}

func (o *OpenAI) String() string {
	return OPENAI
}
