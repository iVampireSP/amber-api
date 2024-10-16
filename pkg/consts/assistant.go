package consts

import "errors"

var (
	ErrAssistantAlreadyBindTheTool    = errors.New("这个助理已经绑定过此工具了")
	ErrAssistantNotFound              = errors.New("未找到该助理")
	ErrAssistantHasBindToolCantDelete = errors.New("这个助理有绑定的工具，请先移除所有的工具，然后再尝试删除该助理")
	// ErrToolNotBind ErrAssistantHasBindLibraryCantDelete = errors.New("这个助理有绑定的资料库，请先移除助理绑定的资料库，然后再尝试删除该助理")
	ErrToolNotBind                 = errors.New("该工具没有绑定该助理")
	ErrAlreadyFavorite             = errors.New("already favorite")
	ErrNotFavorite                 = errors.New("没有 favorite")
	ErrAssistantCannotFavoriteSelf = errors.New("不能收藏自己的助理")
	ErrAssistantNotPublic          = errors.New("该助理不是公开的")
)
