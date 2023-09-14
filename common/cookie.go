package common

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func RemoveCookie(c *gin.Context, key string) {
	domain := c.Request.Host
	fmt.Println("remove域名：" + domain)
	c.SetCookie(key, "", -1, "/", domain, false, false)
}

func AddCookie(c *gin.Context, key string, value string, maxAge int) {
	domain := c.Request.Host
	fmt.Println("add域名：" + domain)
	c.SetCookie(key, value, maxAge, "/", domain, false, false)
}

func AddPushToken(c *gin.Context, token string) {
	domain := c.Request.Host
	//domain中以.pushplus.plus结尾
	if strings.HasSuffix(domain, ".pushplus.plus") {
		domain = ".pushplus.plus"
	}
	fmt.Println("addpush域名：" + domain)
	c.SetCookie("pushToken", token, 7*3600*24, "/", domain, false, false)
}

func RemovePushToken(c *gin.Context) {
	domain := c.Request.Host
	//domain中以.pushplus.plus结尾
	if strings.HasSuffix(domain, ".pushplus.plus") {
		domain = ".pushplus.plus"
	}
	c.SetCookie("pushToken", "", -1, "/", domain, false, false)
}

func GetPushToken(c *gin.Context) string {
	pushToken, _ := c.Cookie("pushToken")
	if len(pushToken) <= 0 {
		pushToken = c.GetHeader("pushToken")

		// if len(pushToken) <= 0 {
		// 	pushToken = GetSession[string](c, "pushToken")
		// }
	}
	return pushToken
}
