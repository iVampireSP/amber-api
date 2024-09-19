package builtin_tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"rag-new/internal/schema"
	"strconv"
)

type searchWebRequest struct {
	Query string `json:"query" mapstructure:"query"`
}

//type SerpResults struct {
//	Results []struct {
//		Body  string `json:"body"`
//		Href  string `json:"href"`
//		Title string `json:"title"`
//	} `json:"results"`
//}

const SerpMaxResult = 10

var ErrOnSearch = fmt.Errorf("暂时无法搜索")

func (s *Service) SearchWeb(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}

	var params searchWebRequest
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	var queryParams = map[string]string{
		"q":           params.Query,
		"max_results": strconv.Itoa(SerpMaxResult),
	}

	var url = s.config.ThirdParty.InternalSerpAPI

	if url == "" {
		response.Content = ErrOnSearch.Error()
		return response, ErrOnSearch
	}

	// 拼接 get 参数
	url += "/search?"
	for k, v := range queryParams {
		url = fmt.Sprintf("%s&%s=%s", url, k, v)
	}

	rsp, err := http.Get(url)
	if err != nil {
		s.logger.Sugar.Fatal(err)
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
		// handle error
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
