package builtin_tool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"rag-new/internal/schema"
	"strings"
	"time"
)

const WebSearchPrompt = `You have the tool browser. Use browser in the following circumstances:

- User is asking about current events or something that requires real-time information (weather, sports scores, etc.)
- User is asking about some term you are totally unfamiliar with (it might be new)
- User explicitly asks you to browse or provide links to references
`
const BingAPI = "https://api.bing.microsoft.com/v7.0/search"

const JinaReader = "https://r.jina.ai"

var (
	ErrOnSearch = fmt.Errorf("暂时无法搜索")
	ErrOnReader = fmt.Errorf("无法读取网页内容")
)

type browserRequest struct {
	QueryOrUrl string `json:"query_or_url" mapstructure:"query_or_url"`
}

type BingAPIResult struct {
	Type     string `json:"_type"`
	WebPages struct {
		WebSearchUrl          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			Id                         string    `json:"id"`
			Name                       string    `json:"name"`
			Url                        string    `json:"url"`
			DatePublished              string    `json:"datePublished"`
			DatePublishedDisplayText   string    `json:"datePublishedDisplayText,omitempty"`
			IsFamilyFriendly           bool      `json:"isFamilyFriendly"`
			DisplayUrl                 string    `json:"displayUrl"`
			Snippet                    string    `json:"snippet"`
			DateLastCrawled            time.Time `json:"dateLastCrawled"`
			CachedPageUrl              string    `json:"cachedPageUrl"`
			Language                   string    `json:"language"`
			IsNavigational             bool      `json:"isNavigational"`
			NoCache                    bool      `json:"noCache"`
			SiteName                   string    `json:"siteName"`
			DatePublishedFreshnessText string    `json:"datePublishedFreshnessText,omitempty"`
			ThumbnailUrl               string    `json:"thumbnailUrl,omitempty"`
			PrimaryImageOfPage         struct {
				ThumbnailUrl string `json:"thumbnailUrl"`
				Width        int    `json:"width"`
				Height       int    `json:"height"`
				SourceWidth  int    `json:"sourceWidth"`
				SourceHeight int    `json:"sourceHeight"`
				ImageId      string `json:"imageId"`
			} `json:"primaryImageOfPage,omitempty"`
		} `json:"value"`
	} `json:"webPages"`
}

func (s *Service) Browser(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}

	var params browserRequest
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	// 验证 params.QueryOrUrl 是否是真的 URL
	if isUrl(params.QueryOrUrl) {
		err = s.readWebContent(params.QueryOrUrl, response)
		if err != nil {
			response.Content = ErrOnReader.Error()
		}
	}

	var queryParams = map[string]string{
		"q": params.QueryOrUrl,
	}
	err = s.serp(queryParams, response)
	if err != nil {
		response.Content = ErrOnSearch.Error()
	}

	return response, nil
}

func isUrl(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func (s *Service) serp(queryParams map[string]string, response *schema.CallBuiltInResponse) error {
	// 拼接 get 参数
	url2 := BingAPI + "?"

	// 构建为 URL Query 参数
	var queryData = url.Values{}
	for key, value := range queryParams {
		queryData.Set(key, value)
	}

	// 拼接 url
	url2 += httpBuildQuery(queryData)

	req, err := http.NewRequest("GET", url2, nil)
	if err != nil {
		response.Content = ErrOnReader.Error()
		return ErrOnReader
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.ThirdParty.BingAPIKey)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Sugar.Error(err)
		response.Content = ErrOnSearch.Error()
		return ErrOnSearch
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			s.logger.Sugar.Error(err)
		}
	}()

	//body, err := io.ReadAll(rsp.Body)
	//if err != nil {
	//	// handle error
	//	response.Content = ErrOnSearch.Error()
	//	return response, ErrOnSearch
	//}

	var serpResults = BingAPIResult{}
	err = json.NewDecoder(rsp.Body).Decode(&serpResults)
	if err != nil {
		response.Content = ErrOnSearch.Error()
		return ErrOnSearch
	}

	for _, result := range serpResults.WebPages.Value {
		response.Content += fmt.Sprintf(`
Title: %s
URL: %s
datePublished: %s
datePublishedFreshnessText: %s
datePublishedDisplayText: %s
snippet: %s
language: %s
siteName: %s

`, result.Name,
			result.Url,
			result.DatePublished,
			result.DatePublishedFreshnessText,
			result.DatePublishedDisplayText,
			result.Snippet,
			result.Language,
			result.SiteName,
		)
	}

	return nil
}

func (s *Service) readWebContent(url string, response *schema.CallBuiltInResponse) error {
	var jinaKey = s.config.ThirdParty.JinaAIKey
	var fullUrl = JinaReader + "/" + url

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		response.Content = ErrOnReader.Error()
		return ErrOnReader
	}

	if jinaKey != "" {
		req.Header.Set("Authorization", "Bearer "+jinaKey)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Sugar.Error(err)
		response.Content = ErrOnSearch.Error()
		return ErrOnSearch
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
		return ErrOnSearch
	}

	//
	//var serpResults = SerpResults{}
	//err = json.NewDecoder(rsp.Body).Decode(&serpResults)
	//if err != nil {
	//	response.Content = ErrOnSearch.Error()
	//	return response, ErrOnSearch
	//}

	response.Content = string(body)

	return nil
}

func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
