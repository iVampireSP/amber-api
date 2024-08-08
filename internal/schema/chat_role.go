package schema

type ChatRole string

var (
	RoleAssistant  ChatRole = "assistant"
	RoleHuman      ChatRole = "user"
	RoleSystem     ChatRole = "system"
	RoleHideSystem ChatRole = "system_hide"
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
