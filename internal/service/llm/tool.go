package llm

import (
	"bytes"
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strings"
)

// spiltFunctionName 将函数名分割为 entity_toolName （entity 为 *entity.Tool，toolName 为 string）的形式
func (s *Service) spiltFunctionName(functionName string) (prefix string, realFunctionName string) {
	// 根据 _ 分割
	var functionNames = strings.Split(functionName, "_")

	// 从第 1 个开始取到最后一个
	var toolName = strings.Join(functionNames[1:], "_")

	return functionNames[0], toolName
}

func (s *Service) GetToolById(ctx context.Context, id schema.EntityId) (*entity.Tool, error) {
	return s.ToolService.GetTool(ctx, id)
}

// callRemoteFunction 可以调用远程函数
func (s *Service) callRemoteFunction(tool *entity.Tool, llmChat *schema.LLMChat, functionName string, args schema.FunctionCallArguments) (*schema.ToolRemoteResponse, error) {
	if !s.config.Debug.Enabled {
		internalAddress, err := s.ToolService.IsAllowed(tool.Data.CallbackUrl)
		if err != nil {
			return nil, err
		}
		if internalAddress {
			return nil, consts.ErrToolAddressIsInternal
		}
	}

	var toolRequest = &schema.ToolRemoteRequest{
		FunctionName: functionName,
		Parameters:   args,
		Chat:         llmChat.Chat,
		//ToolCallToken: llmChat.ToolCallToken,
	}

	s.Logger.Sugar.Infof("Calling remote function: %v", toolRequest)

	if llmChat.UserPublicInfo != nil {
		toolRequest.User = llmChat.UserPublicInfo
	}

	toolRequestJson, err := sonic.Marshal(toolRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tool.Data.CallbackUrl, bytes.NewBuffer(toolRequestJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", consts.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	if tool.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+tool.ApiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyJson := &schema.ToolRemoteResponse{}

	err = sonic.Unmarshal(body, bodyJson)
	if err != nil {
		return nil, err
	}

	if bodyJson.Success {
		return bodyJson, nil
	}

	return bodyJson, errors.New(bodyJson.Content)
}
