package user

import (
	"light/models"
)

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateResponse struct {
	Username string `json:"username"`
}

type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

type ListResponse struct {
	TotalCount uint64             `json:"totalCount"`
	UserList   []*models.UserInfo `json:"userList"`
}

type SwaggerListResponse struct {
	TotalCount uint64            `json:"totalCount"`
	UserList   []models.UserInfo `json:"userList"`
}
