package v1

import "github.com/google/wire"

var ProviderApiControllerSet = wire.NewSet(
	NewUserController,
	NewToolController,
	NewAssistantController,
	NewChatController,
)
