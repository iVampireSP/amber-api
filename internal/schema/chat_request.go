package schema

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}

	// 超过 2038 年后的时间将会失败

	tz, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}

	ct.Time, err = time.ParseInLocation(ctLayout, s, tz)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

type ChatRequest struct {
	ChatId EntityId `uri:"id" binding:"required"`
}

type ChatCreateRequest struct {
	Name        string      `json:"name" binding:"required" validate:"max=30"`
	AssistantId *EntityId   `json:"assistant_id"`
	ExpiredAt   *CustomTime `json:"expired_at" time_format:"2006-01-02" time_utc:"1"`
	UserId      UserId      `json:"user_id" swaggerignore:"true" binding:"-"`
}

type ChatUpdateRequest struct {
	Name        string      `json:"name" binding:"required" validate:"max=30"`
	AssistantId *EntityId   `json:"assistant_id"`
	ExpiredAt   *CustomTime `json:"expired_at" time_format:"2006-01-02" time_utc:"1"`
}

type ChatGuestCreateRequest struct {
	Name        string   `json:"name" binding:"required" validate:"max=30"`
	AssistantId EntityId `json:"assistant_id" binding:"required"`
	GuestID     string   `json:"guest_id" binding:"required" validate:"max=32"`
}

type ChatMessageAddRequest struct {
	AssistantId *EntityId         `json:"assistant_id"`
	Message     string            `json:"message" binding:"required" validate:"max=255"`
	Role        ChatRole          `json:"role" binding:"required" enums:"user,user_hide,system,system_hide,assistant,image"`
	Variables   map[string]string `json:"variables"`
}

type ChatMessageAddImageRequest struct {
	Image *multipart.FileHeader `form:"image" swaggerignore:"true"`
}

type ChatMessageAddFileRequest struct {
	File *multipart.FileHeader `form:"file" swaggerignore:"true"`
}

type ChatMessageResponse struct {
	StreamId string `json:"stream_id"`
	Stream   bool   `json:"stream"`
}

type ChatDownloadRemoteFileRequest struct {
	Url string `form:"url"   binding:"required"`
}

type ChatStreamRequest struct {
	StreamId string `uri:"stream_id" binding:"required"`
}
