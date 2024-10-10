package v1

import (
	"github.com/gin-gonic/gin"
	"net"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"rag-new/pkg/consts"
	"rag-new/pkg/safetpl"
	"strings"
	"time"
)

type promptOptions struct {
	Assistant             *entity.Assistant
	User                  *schema.UserPublicInfo
	Owner                 schema.ChatOwner
	OverrideDefaultPrompt string
	Variables             map[string]string
}

func (u *ChatController) getPrompt(c *gin.Context, options *promptOptions) (string, error) {
	var prompts []string

	var disableDefaultPrompt = false
	var disableMemory = false
	var disableUserPrompt = true
	var disableAssistantPrompt = true

	// 以下设置按照优先级排序

	// 1. 如果用户不是 nil，则启用记忆和 user prompt
	if options.User != nil {
		disableMemory = false
		disableUserPrompt = false
	}

	// 如果没有指定 Assistant，则启用默认 prompt
	if options.Assistant == nil {
		// 那么使用默认的 prompt
		disableDefaultPrompt = false
	} else {
		disableAssistantPrompt = false

		// 如果助理设置了禁用系统默认 Prompt
		if options.Assistant.DisableDefaultPrompt {
			// 禁用系统默认 prompt
			disableDefaultPrompt = true
		}
		// 如果要禁用记忆
		if options.Assistant.DisableMemory {
			disableMemory = true
		}
	}

	// 如果有覆盖的 prompt，则禁用默认 prompt
	if options.OverrideDefaultPrompt != "" {
		disableDefaultPrompt = true

		prompts = append(prompts, options.OverrideDefaultPrompt)
	}

	// 访客模式下，禁用记忆
	if options.Owner == schema.OwnerGuest {
		disableMemory = true

		// 例外情况：如果用户要求 启用记忆
		if options.Assistant.EnableMemoryForAssistantAPI {
			disableMemory = false
		}
	}

	// 应用更改
	if !disableDefaultPrompt {
		prompts = append(prompts, consts.DefaultPrompt)
	}

	if !disableMemory {
		var userId schema.UserId
		if options.User != nil {
			userId = options.User.Id
		} else if options.Owner == schema.OwnerGuest && options.Assistant != nil {
			userId = options.Assistant.UserId
		}

		if userId != "" {
			memoryPrompt, err := u.memoryService.GenerateMemoryPrompt(c, userId)
			if err != nil {
				return "", err
			}
			if memoryPrompt != "" {
				memoryPrompt = "用户的喜好: \n" + memoryPrompt
				prompts = append(prompts, memoryPrompt)
			}
		}

	}

	if !disableUserPrompt {
		var currentTime = time.Now().Format("2006-01-02 15:04:05")
		var userPrompt = "关于：\n - 服务器时间: " + currentTime

		if options.Assistant != nil {
			userPrompt += "\n - 你的名字: " + options.Assistant.Name
		}
		if options.User != nil {
			userPrompt += "\n - 用户的名字: " + options.User.Name + "\n- 用户 ID: " + options.User.Id.String() + "\n"
		}

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
			userPrompt += `- 用户 IP: ` + clientIP
		}

		prompts = append(prompts, userPrompt)
	}

	// 如果没有禁用助理的 prompt
	if !disableAssistantPrompt {

		prompts = append(prompts, options.Assistant.Prompt)
	}

	var prompt = strings.Join(prompts, "\n")

	// 渲染模板
	if options.Variables != nil && len(options.Variables) > 0 {
		prompt = safetpl.RenderTemplate(prompt, options.Variables)

		// 如果 options.Variables 有 now
		if options.Variables["now"] != "" {
			prompt += "\n用户时间: " + options.Variables["now"]
		}
	}

	return prompt, nil
}
