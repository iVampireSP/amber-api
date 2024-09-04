package schema

import (
	"strconv"
	"time"
)

type UserTokenInfo struct {
	Aud              string    `json:"aud"`
	Iss              string    `json:"iss"`
	Iat              float64   `json:"iat"`
	Exp              float64   `json:"exp"`
	Sub              UserId    `json:"sub" mapstructure:"-"`
	Scopes           []string  `json:"scopes"`
	Id               int       `json:"id"`
	Uuid             string    `json:"uuid"`
	Avatar           string    `json:"avatar"`
	Name             string    `json:"name"`
	EmailVerified    bool      `json:"email_verified"`
	RealNameVerified bool      `json:"real_name_verified"`
	PhoneVerified    bool      `json:"phone_verified"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	CreatedAt        time.Time `json:"created_at"`
}

type User struct {
	Token UserTokenInfo
	Valid bool
}

type UserId int64

func (u UserId) String() string {
	return strconv.FormatInt(int64(u), 10)
}

type JWTTokenTypes string

const (
	JWTAccessToken JWTTokenTypes = "access_token"
	JWTIDToken     JWTTokenTypes = "id_token"
)

type UserPublicInfo struct {
	Name      string    `json:"name"`
	Id        *UserId   `json:"id"`
	GuestId   *string   `json:"guest_id"`
	ChatOwner ChatOwner `json:"chat_owner"`
}

func (jwtTokenTypes JWTTokenTypes) String() string {
	return string(jwtTokenTypes)
}
