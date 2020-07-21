package app

import (
	"github.com/gin-gonic/gin"
	"github.com/xzy773599929/blog-service/pkg/errcode"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx:ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	//JSON将给定的结构作为JSON序列化到响应主体中。
	//它还将内容类型设置为“application/json”。
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			GetPage(r.Ctx),
			GetPageSize(r.Ctx),
			totalRows,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error)  {
	response := gin.H{
		"code": err.Code(),
		"msg": err.Msg(),
	}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}