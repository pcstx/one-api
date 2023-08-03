package middleware

import (
	"net/http"
	"one-api/common"
	"one-api/controller"
	"one-api/model"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func authHelper(c *gin.Context, minRole int) {
	session := sessions.Default(c)
	username := session.Get("username")
	role := session.Get("role")
	id := session.Get("id")
	status := session.Get("status")

	//先从cookie中读取pushToken，不存在的话从hearder中读取,最后从session中读取
	pushToken, _ := c.Cookie("pushToken")
	if len(pushToken) <= 0 {
		pushToken = c.GetHeader("pushToken")
	}
	if len(pushToken) <= 0 {
		pushToken = session.Get("pushToken").(string)
	}

	if len(pushToken) > 0 {
		//pushToken存在，但是session不存在用户对象，说明是在第三方中登录，要实现单点
		if id == nil {
			//根据pushToken获取用户信息，再回写到session中
			//等于登录逻辑
			session.Set("pushToken", pushToken)
			session.Save()
			userInfo, _ := controller.GetMyInfo(pushToken)
			controller.WeChatLogin(c, userInfo)
		}
	} else {
		//如果pushToken不存在，说明没有登录或者登录超时
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录或登录已过期，请重新登录！",
		})
		c.Abort()
		return
	}

	if status.(int) == common.UserStatusDisabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已被封禁",
		})
		c.Abort()
		return
	}
	if role.(int) < minRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权进行此操作，权限不足",
		})
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Set("role", role)
	c.Set("id", id)
	c.Set("pushToken", pushToken)
	c.Next()
}

func authHelper2(c *gin.Context, minRole int) {
	session := sessions.Default(c)
	username := session.Get("username")
	role := session.Get("role")
	id := session.Get("id")
	status := session.Get("status")
	if username == nil {
		// Check access token
		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无权进行此操作，未登录且未提供 access token",
			})
			c.Abort()
			return
		}
		user := model.ValidateAccessToken(accessToken)
		if user != nil && user.Username != "" {
			// Token is valid
			username = user.Username
			role = user.Role
			id = user.Id
			status = user.Status
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "无权进行此操作，access token 无效",
			})
			c.Abort()
			return
		}
	}
	if status.(int) == common.UserStatusDisabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已被封禁",
		})
		c.Abort()
		return
	}
	if role.(int) < minRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无权进行此操作，权限不足",
		})
		c.Abort()
		return
	}
	c.Set("username", username)
	c.Set("role", role)
	c.Set("id", id)
	c.Next()
}

func UserAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleCommonUser)
	}
}

func AdminAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleAdminUser)
	}
}

func RootAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleRootUser)
	}
}

func TokenAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("Authorization")
		key = strings.TrimPrefix(key, "Bearer ")
		key = strings.TrimPrefix(key, "sk-")
		parts := strings.Split(key, "-")
		key = parts[0]
		token, err := model.ValidateUserToken(key)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "one_api_error",
				},
			})
			c.Abort()
			return
		}
		if !model.CacheIsUserEnabled(token.UserId) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"message": "用户已被封禁",
					"type":    "one_api_error",
				},
			})
			c.Abort()
			return
		}
		c.Set("id", token.UserId)
		c.Set("token_id", token.Id)
		c.Set("token_name", token.Name)
		requestURL := c.Request.URL.String()
		consumeQuota := true
		if strings.HasPrefix(requestURL, "/v1/models") {
			consumeQuota = false
		}
		c.Set("consume_quota", consumeQuota)
		if len(parts) > 1 {
			if model.IsAdmin(token.UserId) {
				c.Set("channelId", parts[1])
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"message": "普通用户不支持指定渠道",
						"type":    "one_api_error",
					},
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
