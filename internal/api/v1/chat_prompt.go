package v1

import (
	"github.com/gin-gonic/gin"
	"net"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"rag-new/pkg/safetpl"
	"time"
)

func (u *ChatController) getPrompt(c *gin.Context,
	assistant *entity.Assistant,
	user *schema.UserPublicInfo,
	owner schema.ChatOwner,
	variable map[string]string) (string, error) {
	var prompt = "When encountering problems, you must first observe the problem and then think about what to do next, and output your thoughts.\n"

	var currentTime = time.Now().Format("2006-01-02 15:04:05")
	var userPrompt = "Server Time: " + currentTime

	if assistant != nil {
		userPrompt += "\nYour(assistant) name: " + assistant.Name
	}
	if user != nil {
		userPrompt += "\nUsername: " + user.Name + "\nUserId: " + user.Id.String()
	}

	// 如果没有指定 Assistant
	if assistant == nil {
		// 如果有用户的情况下才启用记忆
		if user != nil {
			// 默认启用记忆
			memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, user.Id)
			if err != nil {
				return "", err
			}

			prompt += consts.DefaultPrompt

			prompt += "\nUser memory you know: " + memoryPrompt + "\n"
			prompt += userPrompt
		}
		// 如果用户是 nil 的话，使用默认 Prompt
		prompt += consts.DefaultPrompt
	} else if assistant.DisableDefaultPrompt {
		// 如果禁用了默认的 Prompt
		prompt += assistant.Prompt

		if user != nil {
			// 那还是可以获取记忆
			memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, user.Id)
			if err != nil {
				return "", err
			}

			prompt += "\nUser memory you know: " + memoryPrompt + "\n"
		}

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

		// 助理 API 是禁用记忆的
		if owner == schema.OwnerGuest {
			useMemory = false

			// 例外情况：如果用户要求 启用记忆
			if assistant.EnableMemoryForAssistantAPI {
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
			if variable != nil && len(variable) > 0 {
				prompt += "\n" + safetpl.RenderTemplate(assistant.Prompt, variable)
			} else {
				prompt += "\n" + assistant.Prompt
			}
		}
	}

	prompt += "\n" + safetpl.RenderTemplate("User time: {now}", variable)

	return prompt, nil
}
