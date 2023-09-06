package controller

import (
	"fmt"
	"net/http"
	"one-api/common"
	"one-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTokens(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	p, _ := strconv.ParseInt(c.Query("p"), 10, 64)
	if p < 0 {
		p = 0
	}
	page, err := model.GetAllUserTokensPageList(userId, p+1, int64(common.ItemsPerPage))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "",
		"data":        page.Data,
		"total":       page.Total,
		"pages":       page.Pages,
		"currentPage": page.CurrentPage,
		"pageSize":    page.PageSize,
	})
	return
}

// 获取用户下所有的令牌
func GetUserTokens(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	page, err := model.GetUserTokens(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    page,
	})
	return
}

func SearchTokens(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	keyword := c.Query("keyword")

	p, _ := strconv.ParseInt(c.Query("p"), 10, 64)
	if p < 0 {
		p = 0
	}
	page, err := model.SearchUserTokensPageList(userId, keyword, p+1, int64(common.ItemsPerPage))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "",
		"data":        page.Data,
		"total":       page.Total,
		"pages":       page.Pages,
		"currentPage": page.CurrentPage,
		"pageSize":    page.PageSize,
	})
	return
}

func GetToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userId := c.GetInt("id")
	fmt.Printf("c.GetInt中的userId=%v\n", userId)
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}
	fmt.Printf("session中的userId=%v\n", userId)

	token, err := model.GetTokenByIds(id, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    token,
	})
	return
}

func GetTokenStatus(c *gin.Context) {
	tokenId := c.GetInt("token_id")
	userId := c.GetInt("id")
	token, err := model.GetTokenByIds(tokenId, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	expiredAt := token.ExpiredTime
	if expiredAt == -1 {
		expiredAt = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"object":          "credit_summary",
		"total_granted":   token.RemainQuota,
		"total_used":      0, // not supported currently
		"total_available": token.RemainQuota,
		"expires_at":      expiredAt * 1000,
	})
}

/**
 * @api {get} /api/token/quota 获取令牌额度
**/
func GetTokenQuota(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	tokenQuota, err := model.GetTokenQuota(userId, token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    tokenQuota,
	})
}

func SetDefaultToken(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	if userId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}
	err := model.SetDefaultToken(userId, token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

func AddToken(c *gin.Context) {
	token := model.Token{}
	err := c.ShouldBindJSON(&token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if len(token.Name) > 30 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "令牌名称过长",
		})
		return
	}
	userId := c.GetInt("id")
	fmt.Printf("c.GetInt中的userId=%v", userId)
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}
	fmt.Printf("session中的userId=%v", userId)

	cleanToken := model.Token{
		UserId:         userId,
		Name:           token.Name,
		Key:            common.GenerateKey(),
		CreatedTime:    common.GetTimestamp(),
		AccessedTime:   common.GetTimestamp(),
		ExpiredTime:    token.ExpiredTime,
		RemainQuota:    token.RemainQuota,
		UnlimitedQuota: token.UnlimitedQuota,
	}
	err = cleanToken.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func DeleteToken(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	err := model.DeleteTokenById(id, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func UpdateToken(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		userId = common.GetSession[int](c, "id")
	}

	statusOnly := c.Query("status_only")
	token := model.Token{}
	err := c.ShouldBindJSON(&token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if len(token.Name) > 30 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "令牌名称过长",
		})
		return
	}
	cleanToken, err := model.GetTokenByIds(token.Id, userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if token.Status == common.TokenStatusEnabled {
		if cleanToken.Status == common.TokenStatusExpired && cleanToken.ExpiredTime <= common.GetTimestamp() && cleanToken.ExpiredTime != -1 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "令牌已过期，无法启用，请先修改令牌过期时间，或者设置为永不过期",
			})
			return
		}
		if cleanToken.Status == common.TokenStatusExhausted && cleanToken.RemainQuota <= 0 && !cleanToken.UnlimitedQuota {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "令牌可用额度已用尽，无法启用，请先修改令牌剩余额度，或者设置为无限额度",
			})
			return
		}
	}
	if statusOnly != "" {
		cleanToken.Status = token.Status
	} else {
		// If you add more fields, please also update token.Update()
		cleanToken.Name = token.Name
		cleanToken.ExpiredTime = token.ExpiredTime
		cleanToken.RemainQuota = token.RemainQuota
		cleanToken.UnlimitedQuota = token.UnlimitedQuota
	}
	err = cleanToken.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    cleanToken,
	})
	return
}
