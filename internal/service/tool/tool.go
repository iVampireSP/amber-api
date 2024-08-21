package tool

import (
	"context"
	"io"
	"net"
	"net/http"
	url2 "net/url"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"strings"

	"github.com/go-playground/validator/v10"

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
	if !s.config.Debug.Enabled {
		internalAddress, err := s.IsAllowed(tool.Url)
		if err != nil {
			return nil, err
		}
		if internalAddress {
			return nil, consts.ErrToolAddressIsInternal
		}
	}

	var toolEntity entity.Tool

	toolEntity.UserId = userId

	toolEntity.Name = tool.Name
	toolEntity.Description = tool.Description
	toolEntity.DiscoveryUrl = tool.Url
	toolEntity.ApiKey = tool.ApiKey

	toolData, err := s.getToolData(tool.Url)
	if err != nil {
		return nil, err
	}

	err = s.ValidateSyntax(toolData)
	if err != nil {
		return nil, err
	}

	if !s.config.Debug.Enabled {
		internalAddress, err := s.IsAllowed(toolData.CallbackUrl)
		if err != nil {
			return nil, err
		}
		if internalAddress {
			return nil, consts.ErrToolAddressIsInternal
		}
	}

	_, err = s.x.Context(ctx).Insert(&toolEntity)
	toolData.ToolId = toolEntity.Id
	toolEntity.Data = *toolData.Output()

	_, err = s.x.Context(ctx).ID(toolEntity.Id).Cols("data").Update(&toolEntity)

	return &toolEntity, err
}

func (s *Service) UpdateToolData(ctx context.Context, tool *entity.Tool) error {
	if !s.config.Debug.Enabled {
		internalAddress, err := s.IsAllowed(tool.DiscoveryUrl)
		if err != nil {
			return err
		}
		if internalAddress {
			return consts.ErrToolAddressIsInternal
		}

		internalAddress, err = s.IsAllowed(tool.Data.CallbackUrl)
		if err != nil {
			return err
		}
		if internalAddress {
			return consts.ErrToolAddressIsInternal
		}
	}

	toolData, err := s.getToolData(tool.DiscoveryUrl)

	if err != nil {
		return err
	}

	toolData.ToolId = tool.Id

	err = s.ValidateSyntax(toolData)
	if err != nil {
		return err
	}

	tool.Data = *toolData.Output()

	_, err = s.x.Context(ctx).ID(tool.Id).Cols("data").Update(tool)

	return err
}

func (s *Service) getToolData(url string) (*schema.ToolDiscoveryInput, error) {
	var toolData schema.ToolDiscoveryInput

	// post url
	resp, err := http.Get(url)
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

	return &toolData, err
}

func (s *Service) DeleteTool(ctx context.Context, id int64) error {
	// 做检查，不能删除已经绑定的 tool
	count, err := s.x.Context(ctx).Where("tool_id = ?", id).Count(&entity.AssistantTool{})
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ErrToolFailedDeleteBecauseHasBindAssistant
	}

	_, err = s.x.Context(ctx).ID(id).Delete(&entity.Tool{})
	return err
}

//func (s *Service) UpdateTool(ctx context.Context, id schema.ToolId, tool *schema.ToolUpdateRequest) error {
//	_, err := s.x.Context(ctx).ID(id).AllCols().Update(tool)
//	return err
//}

func (s *Service) GetTool(ctx context.Context, id int64) (*entity.Tool, error) {
	var tool = entity.Tool{}
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

func (s *Service) ValidateSyntax(toolDiscoveryOutput *schema.ToolDiscoveryInput) error {
	var validate = validator.New()
	err := validate.Struct(toolDiscoveryOutput)

	return err
}

// IsAllowed 检测是否允许使用此 URL
func (s *Service) IsAllowed(url string) (bool, error) {
	urlParse, err := url2.Parse(url)
	if err != nil {
		return false, err
	}

	host := urlParse.Hostname()

	// 如果在集群内
	if strings.HasSuffix(host, "cluster.local") || strings.HasSuffix(host, ".svc") {
		return true, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return false, err
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return true, nil
		}
	}

	return false, nil
}
