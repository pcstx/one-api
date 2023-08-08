package common

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession[T any](c *gin.Context, key string) T {
	var t T

	session := sessions.Default(c)
	value := session.Get(key)
	if value != nil {
		t = value.(T)
	}
	return t
}

func SetSession[T any](c *gin.Context, key string, value T) bool {
	session := sessions.Default(c)
	session.Set(key, value)
	err := session.Save()
	if err != nil {
		return false
	}
	return true
}

func ClearSession(c *gin.Context) bool {
	session := sessions.Default(c)

	session.Options(sessions.Options{MaxAge: -1})

	session.Clear()
	err := session.Save()
	if err != nil {
		return false
	}
	return true
}
