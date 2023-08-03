package common

import "github.com/gin-gonic/gin"

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
	c.SetCookie("pushToken", token, 7*3600*24, "/", domain, false, false)
}
