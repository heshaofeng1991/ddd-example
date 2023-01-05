// Package interfaces provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package interfaces

import (
	"time"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Announcement defines model for Announcement.
type Announcement struct {
	// 公告id
	AnnouncementId int `json:"announcement_id"`

	// 公告内容
	Content string `json:"content"`

	// 创建时间
	CreatedAt time.Time `json:"created_at"`

	// 过期时间
	Expiration string `json:"expiration"`

	// 公告状态1:启用 0:停用
	Status int `json:"status"`

	// 公告标题
	Title string `json:"title"`

	// 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// AnnouncementsResp defines model for AnnouncementsResp.
type AnnouncementsResp struct {
	// code (错误码).
	Code int            `json:"code"`
	Data []Announcement `json:"data"`

	// message (错误信息).
	Message string `json:"message"`
}