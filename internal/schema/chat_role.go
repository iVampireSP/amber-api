package schema

type ChatRole string

var (
	RoleAssistant ChatRole = "assistant"
	RoleHuman     ChatRole = "user"
	// RoleHideHuman 为用户隐藏消息，不会在聊天历史中显示，但是会交给 LLM 处理
	RoleHideHuman ChatRole = "user_hide"
	RoleSystem    ChatRole = "system"
	// RoleHideSystem 为系统隐藏消息，不会在聊天历史中显示，但是会交给 LLM 处理
	RoleHideSystem ChatRole = "system_hide"
	// RoleFile 处理用于上传的文件，content 为 files 中的 id
	RoleFile ChatRole = "file"
	// RoleToolCall 是工具调用信息，需要携带 Tool Call Info
	RoleToolCall ChatRole = "assistant_call"
	// RoleTool 是工具的响应，需要携带 Tool ID
	RoleTool ChatRole = "tool"
)

type ChatOwner string

var (
	OwnerUser  ChatOwner = "user"
	OwnerGuest ChatOwner = "guest"
)

func (cr ChatRole) String() string {
	return string(cr)
}
func (co ChatOwner) String() string {
	return string(co)
}
