package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type TemplateSearch struct {
	Name  string `json:"name" form:"name"`
	Model string `json:"model" form:"model"`
	request.PageInfo
}

type CreateTemplate struct {
	Name        string             `json:"name" binding:"required"`
	ProductName string             `json:"productName"`
	Model       string             `json:"model"`
	Items       []TemplateItemSort `json:"items"`
}

type UpdateTemplate struct {
	ID          uint               `json:"ID" binding:"required"`
	Name        string             `json:"name" binding:"required"`
	ProductName string             `json:"productName"`
	Model       string             `json:"model"`
	Status      *int               `json:"status"`
	Items       []TemplateItemSort `json:"items"`
}

type CopyTemplate struct {
	ID   uint   `json:"ID" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type TemplateItemSort struct {
	ItemID uint `json:"itemID"`
	Sort   int  `json:"sort"`
}
