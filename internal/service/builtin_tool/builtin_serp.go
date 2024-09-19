package builtin_tool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"rag-new/internal/schema"
	"time"
)

type searchWebRequest struct {
	Query string `json:"query" mapstructure:"query"`
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

const BingAPI = "https://api.bing.microsoft.com/v7.0/search"

var ErrOnSearch = fmt.Errorf("暂时无法搜索")

func (s *Service) SearchWeb(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}

	var params searchWebRequest
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	var queryParams = map[string]string{
		"q": params.Query,
	}

	var url2 = BingAPI

	// 拼接 get 参数
	url2 += "?"

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
		return response, ErrOnReader
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", s.config.ThirdParty.BingAPIKey)

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
		return response, ErrOnSearch
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

	return response, nil
}

func httpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
