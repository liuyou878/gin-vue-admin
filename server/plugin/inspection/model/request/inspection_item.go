package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type InspectionItemSearch struct {
	Name       string `json:"name" form:"name"`
	ResultType string `json:"resultType" form:"resultType"`
	request.PageInfo
}

type CreateInspectionItem struct {
	Name       string   `json:"name" binding:"required"`
	ResultType string   `json:"resultType" binding:"required,oneof=pass_fail number both"`
	Unit       string   `json:"unit"`
	MinValue   *float64 `json:"minValue"`
	MaxValue   *float64 `json:"maxValue"`
	Remark     string   `json:"remark"`
}

type UpdateInspectionItem struct {
	ID         uint     `json:"ID" binding:"required"`
	Name       string   `json:"name" binding:"required"`
	ResultType string   `json:"resultType" binding:"required,oneof=pass_fail number both"`
	Unit       string   `json:"unit"`
	MinValue   *float64 `json:"minValue"`
	MaxValue   *float64 `json:"maxValue"`
	Remark     string   `json:"remark"`
}
