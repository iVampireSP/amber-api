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
	// RoleKnowledge 是知识库的响应
	RoleKnowledge ChatRole = "knowledge"
	// RoleSystemOverride 是特殊的一个角色，可以直接覆盖 Chat 的默认提示词
	RoleSystemOverride ChatRole = "system_override"
	// RoleHumanLater 是一个特殊的角色，不会让 AI 立即回复，也不会被拦截
	RoleHumanLater ChatRole = "user_later"
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
