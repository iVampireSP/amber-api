package v1

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"rag-new/pkg/random"
)

func (u *ChatController) getCacheKey(key string) string {
	return fmt.Sprintf("chat:%s", key)
}

type ChatStreamCache struct {
	ChatId    schema.EntityId
	Variables map[string]string
}

func (u *ChatController) generateChatStream(c context.Context,
	chatId schema.EntityId,
	userPublic *schema.UserPublicInfo,
	variables map[string]string) (streamId string, err error) {
	var randomId = random.String(32)
	// 保存 chat stream id
	err = u.redis.Client.Set(c, u.getCacheKey("entity:"+chatId.String()), randomId, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	var csc = ChatStreamCache{
		ChatId:    chatId,
		Variables: variables,
	}

	chatJson, err := sonic.MarshalString(csc)
	if err != nil {
		return "", err
	}

	err = u.redis.Client.Set(c, u.getCacheKey("stream:"+randomId), chatJson, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	userJson, err := sonic.Marshal(userPublic)
	if err != nil {
		return "", err
	}

	err = u.redis.Client.Set(c, u.getCacheKey("stream:"+randomId+":user"), userJson, consts.ChatStreamExpire).Err()
	if err != nil {
		return "", err
	}

	return randomId, nil
}

//func (u *ChatController) GenerateCallToken(ctx context.Context, chatEntity *entity.Chat) (string, error) {
//	tct, err := u.toolService.GenerateToolCallToken(ctx, chatEntity)
//	if err != nil {
//		return "", err
//	}
//
//	if tct == nil {
//		return "", fmt.Errorf("failed to generate tool call token")
//	}
//
//	return tct.Token, nil
//}
