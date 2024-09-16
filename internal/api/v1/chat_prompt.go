package v1

import (
	"github.com/gin-gonic/gin"
	"net"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"time"
)

func (u *ChatController) getPrompt(c *gin.Context, assistant *entity.Assistant, user *schema.UserPublicInfo, owner schema.ChatOwner) (string, error) {
	var prompt = ""

	var currentTime = time.Now().Format("2006-01-02 15:04:05")
	var userPrompt = "Now Time: " + currentTime

	if assistant != nil {
		userPrompt += "\nYour(assistant) name: " + assistant.Name
	}
	if user != nil {
		userPrompt += "\nUsername: " + user.Name + "\nUserId: " + user.Id.String()
	}

	if assistant == nil {
		// 默认启用记忆
		memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, user.Id)
		if err != nil {
			return "", err
		}

		prompt = consts.DefaultPrompt

		prompt += "\nUser memory you know: " + memoryPrompt + "\n"
		prompt += userPrompt

	} else if assistant.DisableDefaultPrompt {
		// 如果禁用了默认的 Prompt
		prompt = assistant.Prompt

		// 那还是可以获取记忆
		memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, user.Id)
		if err != nil {
			return "", err
		}

		prompt += "\nUser memory you know: " + memoryPrompt + "\n"
	} else {
		prompt += userPrompt

		var clientIP = ""

		// 如果 header 里面有 HeaderUserIp
		if c.GetHeader(HeaderUserIp) != "" {
			var headerIP = c.GetHeader(HeaderUserIp)
			var ip = net.ParseIP(headerIP)
			if ip != nil && !ip.IsLoopback() && !ip.IsPrivate() {
				clientIP = headerIP
			}
		}

		if clientIP == "" {
			var cIP = c.ClientIP()
			var ip = net.ParseIP(cIP)
			// 如果是内部 IP
			if ip != nil && !ip.IsLoopback() && !ip.IsPrivate() {
				clientIP = ip.String()
			}
		}

		if clientIP != "" {
			prompt += `
The user(who is talking with you)'s IP: ` + clientIP + "(Not your IP, system hint you, you not have IP address)"
		}

		// 记忆
		var useMemory = true
		if assistant.DisableMemory {
			useMemory = false
		}

		if owner == schema.OwnerGuest {
			if assistant.EnableMemoryForAssistantAPI && !assistant.DisableMemory {
				useMemory = true
			}
		}

		if useMemory {
			memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, assistant.UserId)
			if err != nil {
				return "", err
			}

			prompt += "\nUser memory you know: " + memoryPrompt + "\n"
		}

		if assistant.Prompt != "" {
			prompt += "\n" + assistant.Prompt
		}
	}

	return prompt, nil
}
