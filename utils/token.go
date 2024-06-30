package utils

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/pkoukk/tiktoken-go"
)

// copy from https://github.com/pkoukk/tiktoken-go/blob/main/README_zh-hans.md
func NumTokensFromMessages(content string, model string) int {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		logger.Errorf("EncodingForModel: %v", err)
		return len(content)
	}

	return len(tkm.Encode(content, nil, nil))
}
