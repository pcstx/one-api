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
	// session := sessions.Default(c)
	// username := session.Get("username")
	// role := session.Get("role")
	// id := session.Get("id")
	// status := session.Get("status")
	id := common.GetSession[int](c, "id")

	//先从cookie中读取pushToken，不存在的话从hearder中读取,最后从session中读取
	pushToken := common.GetPushToken(c)
	// if len(pushToken) <= 0 {
	// 	pushToken = c.GetHeader("pushToken")

	// 	if len(pushToken) <= 0 {
	// 		pushToken = common.GetSession[string](c, "pushToken")
	// 	}
	// }

	if len(pushToken) > 0 {
		//pushToken存在，但是session不存在用户对象，说明是在第三方中登录，要实现单点
		//if id <= 0 {
		//根据pushToken获取用户信息，再回写到session中
		//等于登录逻辑
		userInfo, _ := controller.GetMyInfo(pushToken)
		if userInfo != nil && userInfo.Id > 0 {
			common.SetSession[string](c, "pushToken", pushToken)
			controller.AuthLogin(c, userInfo)
		} else {
			//c.SetCookie("pushToken", "", -1, "/", common.PushPlusDomain, false, false)
			common.RemovePushToken(c)
			common.RemoveCookie(c, "session")
			common.ClearSession(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未登录或登录已过期，请重新登录！",
			})
			c.Abort()
			return
		}
		//}
	} else {
		//如果pushToken不存在，说明没有登录或者登录超时
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未登录或登录已过期，请重新登录！",
		})
		c.Abort()
		return
	}

	status := common.GetSession[int](c, "status")
	role := common.GetSession[int](c, "role")
	username := common.GetSession[string](c, "username")

	if status == common.UserStatusDisabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已被封禁",
		})
		c.Abort()
		return
	}
	if role < minRole {
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
		key = strings.TrimPrefix(key, "sp-")
		parts := strings.Split(key, "-")
		key = parts[0]
		token, err := model.ValidateUserToken(key)
		if err != nil {
			abortWithMessage(c, http.StatusUnauthorized, err.Error())
			return
		}
		userEnabled, err := model.CacheIsUserEnabled(token.UserId)
		if err != nil {
			abortWithMessage(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !userEnabled {
			abortWithMessage(c, http.StatusForbidden, "用户已被封禁")
			return
		}
		c.Set("id", token.UserId)
		c.Set("token_id", token.Id)
		c.Set("token_name", token.Name)
		if len(parts) > 1 {
			if model.IsAdmin(token.UserId) {
				c.Set("channelId", parts[1])
			} else {
				abortWithMessage(c, http.StatusForbidden, "普通用户不支持指定渠道")
				return
			}
		}
		c.Next()
	}
}
