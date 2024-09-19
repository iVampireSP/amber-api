package builtin_tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"rag-new/internal/schema"
)

type readUrlRequest struct {
	Url string `json:"url" mapstructure:"url"`
}

var ErrOnReader = fmt.Errorf("无法读取网页内容")

const JinaReader = "https://r.jina.ai"

func (s *Service) ReadUrl(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}

	var params readUrlRequest
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	var jinaKey = s.config.ThirdParty.JinaAIKey

	if jinaKey == "" {
		response.Content = ErrOnReader.Error()
		return response, ErrOnReader
	}

	var fullUrl = JinaReader + "/" + params.Url

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		response.Content = ErrOnReader.Error()
		return response, ErrOnReader
	}

	req.Header.Set("Authorization", "Bearer "+jinaKey)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Sugar.Error(err)
		response.Content = ErrOnSearch.Error()
		return response, ErrOnSearch
	}

	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			s.logger.Sugar.Error(err)
		}
	}()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		response.Content = ErrOnSearch.Error()
		return response, ErrOnSearch
	}

	//
	//var serpResults = SerpResults{}
	//err = json.NewDecoder(rsp.Body).Decode(&serpResults)
	//if err != nil {
	//	response.Content = ErrOnSearch.Error()
	//	return response, ErrOnSearch
	//}

	response.Content = string(body)

	return response, nil
}
