package tool

import (
	"context"
	"io"
	"net/http"
	"rag-new/internal/entity"
	"rag-new/internal/schema"

	"github.com/bytedance/sonic"
)

func (s *Service) ListToolFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Tool, error) {
	var tools []*entity.Tool
	err := s.x.Context(ctx).Where("user_id = ?", userId).Find(&tools)
	if err != nil {
		return nil, err
	}

	return tools, nil
}

func (s *Service) CreateTool(ctx context.Context, tool *schema.ToolCreateRequest, userId schema.UserId) (*entity.Tool, error) {
	var toolEntity entity.Tool
	var toolData schema.ToolDiscoveryInput

	// map to entity
	toolEntity.UserId = userId

	toolEntity.Name = tool.Name
	toolEntity.Description = tool.Description
	toolEntity.DiscoveryUrl = tool.Url
	toolEntity.ApiKey = tool.ApiKey

	// post url
	resp, err := http.Get(tool.Url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// convert to byte
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = sonic.Unmarshal(body, &toolData)
	if err != nil {
		return nil, err
	}

	toolEntity.Data = *toolData.Output()

	_, err = s.x.Context(ctx).Insert(&toolEntity)

	toolEntity.Data.ToolId = toolEntity.ID
	// update
	//_, err = s.x.Context(ctx).ID(toolEntity.ID).AllCols().Update(&toolEntity)

	// only update data
	_, err = s.x.Context(ctx).ID(toolEntity.ID).Cols("data").Update(&toolEntity)

	return &toolEntity, err
}

func (s *Service) DeleteTool(ctx context.Context, id int64) error {
	_, err := s.x.Context(ctx).ID(id).Delete(&entity.Tool{})
	return err
}

//func (s *Service) UpdateTool(ctx context.Context, id schema.ToolId, tool *schema.ToolUpdateRequest) error {
//	_, err := s.x.Context(ctx).ID(id).AllCols().Update(tool)
//	return err
//}

func (s *Service) GetTool(ctx context.Context, id int64) (*entity.Tool, error) {
	var tool entity.Tool
	_, err := s.x.Context(ctx).Where("id = ?", id).Get(&tool)
	return &tool, err
}

func (s *Service) CheckTool(ctx context.Context, url string, userId schema.UserId) (bool, error) {
	count, err := s.x.Context(ctx).Where("user_id = ?", userId).Where("discovery_url = ?", url).Count(&entity.Tool{})
	return count > 0, err
}

func (s *Service) Exists(ctx context.Context, id int64) (bool, error) {
	count, err := s.x.Context(ctx).Where("id = ?", id).Count(&entity.Tool{})
	return count > 0, err
}
