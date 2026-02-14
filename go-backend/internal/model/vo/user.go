package vo

import "time"

type UserVO struct {
	ID          int64     `json:"id,string"`
	UserAccount string    `json:"userAccount"`
	UserName    string    `json:"userName"`
	UserAvatar  string    `json:"userAvatar"`
	UserProfile string    `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	UserStatus  string    `json:"userStatus"`
	TokenQuota  int64     `json:"tokenQuota,string"`
	UsedTokens  int64     `json:"usedTokens,string"`
	Balance     float64   `json:"balance"`
	CreateTime  time.Time `json:"createTime"`
}

type LoginUserVO struct {
	ID          int64     `json:"id,string"`
	UserAccount string    `json:"userAccount"`
	UserName    string    `json:"userName"`
	UserAvatar  string    `json:"userAvatar"`
	UserProfile string    `json:"userProfile"`
	UserRole    string    `json:"userRole"`
	UserStatus  string    `json:"userStatus"`
	TokenQuota  int64     `json:"tokenQuota,string"`
	UsedTokens  int64     `json:"usedTokens,string"`
	Balance     float64   `json:"balance"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}
