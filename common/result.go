package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回的结果：
type Result struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 成功
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Success = true
	res.Message = "请求成功"
	res.Data = data

	c.JSON(http.StatusOK, res)
}

// 出错
func Error(c *gin.Context, errMsg string) {
	res := Result{}

	res.Success = false
	res.Message = errMsg
	res.Data = gin.H{}

	c.JSON(http.StatusOK, res)
}

func Unauthorized(c *gin.Context, errMsg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": errMsg,
	})
}
