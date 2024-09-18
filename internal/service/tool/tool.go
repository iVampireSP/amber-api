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
	"rag-new/pkg/random"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/bytedance/sonic"
)

func (s *Service) ListToolFromUserId(ctx context.Context, userId schema.UserId) ([]*entity.Tool, error) {
	tools, err := s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.UserId.Eq(userId.String())).Find()

	if err != nil {
		return nil, err
	}

	return tools, nil
}

func (s *Service) CreateTool(ctx context.Context, tool *schema.ToolCreateRequest, userId schema.UserId) (*entity.Tool, error) {
	if !s.config.Debug.Enabled {
		internalAddress, err := s.IsAllowed(tool.DiscoveryUrl)
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
	toolEntity.DiscoveryUrl = tool.DiscoveryUrl
	toolEntity.ApiKey = tool.ApiKey

	toolData, err := s.getToolData(tool.DiscoveryUrl)
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

	err = s.dao.WithContext(ctx).Tool.Create(&toolEntity)
	if err != nil {
		return nil, err
	}

	toolData.ToolId = toolEntity.Id

	_, err = s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(toolEntity.Id))).Update(s.dao.Tool.Data, *toolData.Output())

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

	_, err = s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(tool.Id))).Update(s.dao.Tool.Data, *toolData.Output())

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

	// 最多只能有 64 个
	if len(toolData.Functions) > 64 {
		return nil, consts.ErrToolTooManyFunctions
	}

	return &toolData, nil
}

func (s *Service) DeleteTool(ctx context.Context, id schema.EntityId) error {
	// 做检查，不能删除已经绑定的 tool
	count, err := s.dao.WithContext(ctx).AssistantTool.Where(s.dao.AssistantTool.ToolId.Eq(uint(id))).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ErrToolFailedDeleteBecauseHasBindAssistant
	}

	_, err = s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(id))).Delete()

	return err
}

//func (s *Service) UpdateTool(ctx context.Context, id schema.ToolId, tool *schema.ToolUpdateRequest) error {
//	_, err := s.x.Context(ctx).ID(id).AllCols().Update(tool)
//	return err
//}

func (s *Service) GetTool(ctx context.Context, id schema.EntityId) (*entity.Tool, error) {
	tool, err := s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(id))).First()

	return tool, err
}

func (s *Service) CheckTool(ctx context.Context, url string, userId schema.UserId) (bool, error) {
	count, err := s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.UserId.Eq(userId.String())).Where(s.dao.Tool.DiscoveryUrl.Eq(url)).Count()

	return count > 0, err
}

func (s *Service) Exists(ctx context.Context, id schema.EntityId) (bool, error) {
	count, err := s.dao.WithContext(ctx).Tool.Where(s.dao.Tool.Id.Eq(uint(id))).Count()

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

func (s *Service) GenerateToolCallToken(ctx context.Context, chat *entity.Chat) (*entity.ToolCallToken, error) {
	var expiredAt = time.Now().Add(time.Hour * 24)

	var tct = entity.ToolCallToken{
		Token:     random.String(24),
		ChatId:    chat.Id,
		ExpiredAt: expiredAt,
	}

	err := s.dao.WithContext(ctx).ToolCallToken.Create(&tct)

	return &tct, err
}
