package common

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RemoveCookie(c *gin.Context, key string) {
	domain := c.Request.Host
	c.SetCookie("pushToken", "", -1, "/", domain, false, false)
}

func AddCookie(c *gin.Context, key string, value string, maxAge int) {
	domain := c.Request.Host
	c.SetCookie(key, value, maxAge, "/", domain, false, false)
}

func AddPushToken(c *gin.Context, token string) {
	domain := c.Request.Host
	//domain中以.pushplus.plus结尾
	if strings.HasSuffix(domain, ".pushplus.plus") {
		domain = ".pushplus.plus"
	}
	c.SetCookie("pushToken", token, 7*3600*24, "/", domain, false, false)
}
