package v1

import (
	"github.com/gin-gonic/gin"
	"net"
	"rag-new/internal/entity"
	"rag-new/internal/schema"
	"time"
)

func (u *ChatController) getPrompt(c *gin.Context, assistant *entity.Assistant, user *schema.UserPublicInfo) string {
	var prompt = ""

	if assistant.DisableDefaultPrompt {
		prompt = assistant.Prompt
	} else {
		var currentTime = time.Now().Format("2006-01-02 15:04:05")
		prompt = `
Time: ` + currentTime + `
Your(assistant) name: ` + assistant.Name + `(current user give you)` + `
Your(assistant) name: ` + assistant.Name + `(current user give you)` + `
Your description: ` + assistant.Description + "(current user given)"
		if user != nil {
			prompt += `
Username: ` + user.Name + `(system hint you)` + `
UserId: ` + user.Id.String() + "(system hint you, user can't change it)"
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
			prompt += `
The user(who is talking with you)'s IP: ` + clientIP + "(Not your IP, system hint you, you not have IP address)"
		}

		if assistant.Prompt != "" {
			prompt += "\n" + assistant.Prompt
		}
	}

	return prompt
}
