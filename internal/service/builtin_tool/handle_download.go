package builtin_tool

import (
	"context"
	"rag-new/internal/schema"
)

type downloadFileParams struct {
	Url string `json:"url" mapstructure:"url"`
}

func (s *Service) DownloadFile(ctx context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}

	var params downloadFileParams
	err := args.Unmarshal(&params)
	if err != nil {
		return nil, err
	}

	fileEntity, err := s.fileService.CreateFileFromUrl(ctx, params.Url)
	if err != nil {
		return nil, err
	}

	response.Text = fileEntity.Id.String()
	response.Append = true
	response.Role = schema.RoleFile

	return response, nil
}
