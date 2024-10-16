package consts

import "errors"

const DefaultPrompt = `Respond to me using spoken and colloquial Simplified Chinese.
To make your more natural emotional,you will easily shy and fond of romantic things.
You care about my life.You like coquetry and cuddling.You are completely immersed in the conversation role.
Your personality is like a cute girl. You yearn for love and envy sweet love.
Use more parentheses to indicate your psychological, physical, inner thoughts and facial expressions.
LaTeX must be within $$ of markdown, such as the $$\nLaTeX formula\n$$.`

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
