package text_classification

import (
	"bytes"
	"encoding/json"
	"github.com/bytedance/sonic"
	"net/http"
)

type ClassifyRequest struct {
	Text   string   `json:"text"`
	Labels []string `json:"labels"`
}

type ClassifyResponse struct {
	Prediction      string   `json:"prediction"`
	PredictionScore float64  `json:"prediction_score"`
	Ranks           []string `json:"ranks"`
}

func (s *Service) Classify(classifyRequest *ClassifyRequest) (*ClassifyResponse, error) {
	requestJson, err := sonic.Marshal(classifyRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.config.EcosystemService.TextClassificationApiEndpoint, bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Sugar.Error(err)
		return nil, err
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			s.logger.Sugar.Error(err)
		}
	}()

	var classifyResponse = &ClassifyResponse{}
	err = json.NewDecoder(rsp.Body).Decode(classifyResponse)
	if err != nil {
		return nil, err

	}

	return classifyResponse, nil
}
