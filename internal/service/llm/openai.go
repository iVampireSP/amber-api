package llm

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/tmc/langchaingo/llms"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strconv"
)

func (s *Service) GenerateContent(ctx context.Context, llmChat *schema.LLMChat, llmTools []llms.Tool, historyContent []llms.MessageContent) (response *llms.ContentResponse, err error) {
	// 上一个字
	var lastWord = ""
	// 重复次数
	var lastWordRepeatCount = 0

	resp, err := s.OpenAI.GenerateContent(ctx,
		historyContent,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			// 检测长度
			if len(chunk) == 0 {
				return nil
			}

			var stringChunk = string(chunk)

			// 检测是否可以转换为数字或者 float
			if !s.isNumeric(stringChunk) {
				// 检测是否 json，判断是否是工具调用
				var isJson = sonic.Valid(chunk)
				if !isJson {
					// 取 chunk 中最后一个字
					var chunkLastWord = string(chunk[len(chunk)-1])
					// 检测是否是上一个字
					if lastWord == chunkLastWord {
						lastWordRepeatCount++
					} else {
						lastWordRepeatCount = 0
						lastWord = chunkLastWord
					}
					// 如果上一个字重复次数大于 10，就终止
					if lastWordRepeatCount >= 10 {
						s.Logger.Sugar.Warnf("Detected repeated word: %s, chunk: %s", lastWord, string(chunk))
						return consts.ErrWordRepeatedDetected
					}

					s.write(ctx, llmChat, &schema.AssistantResponse{
						State: schema.StateChunk,
						ChunkMessage: &schema.ChunkMessage{
							Content: stringChunk,
						},
						Content: stringChunk,
					})
				}

			} else {
				s.write(ctx, llmChat, &schema.AssistantResponse{
					State: schema.StateChunk,
					ChunkMessage: &schema.ChunkMessage{
						Content: stringChunk,
					},
					Content: stringChunk,
				})
			}

			return nil
		}),
		llms.WithTools(llmTools),
		llms.WithN(llmChat.N),
		llms.WithMaxTokens(llmChat.MaxTokens),
		llms.WithTemperature(llmChat.Temperature),
		llms.WithTopP(llmChat.TopP),
		llms.WithModel(llmChat.Model),
		llms.WithTopK(llmChat.TopK))
	return resp, err
}

func (s *Service) isNumeric(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
