package controller

import (
	"one-api/common"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type TestController struct{}

func (con TestController) GetSession(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	common.Success(c, id)
}

func (con TestController) SetSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("id", 1)
	err := session.Save()
	if err != nil {
		common.Error(c, err.Error())
		return
	}
	common.Success(c, nil)
}

func (con TestController) ClearSession(c *gin.Context) {
	session := sessions.Default(c)

	session.Options(sessions.Options{MaxAge: -1})

	session.Clear()
	err := session.Save()
	if err != nil {
		common.Error(c, err.Error())
		return
	}
	common.Success(c, nil)
}
