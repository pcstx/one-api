package controller

import (
	"one-api/common"
	"one-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllLogs(c *gin.Context) {
	p, _ := strconv.ParseInt(c.Query("p"), 10, 64)
	if p < 0 {
		p = 0
	}
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	username := c.Query("username")
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	page, err := model.GetAllLogsPageList(logType, startTimestamp, endTimestamp, modelName, username, tokenName, p+1, int64(common.ItemsPerPage))
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success":     true,
		"message":     "",
		"data":        page.Data,
		"total":       page.Total,
		"pages":       page.Pages,
		"currentPage": page.CurrentPage,
		"pageSize":    page.PageSize,
	})
}

func GetUserLogs(c *gin.Context) {
	//获取页数
	p, _ := strconv.ParseInt(c.Query("p"), 10, 64)
	if p < 0 {
		p = 0
	}
	userId := c.GetInt("id")
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	page, err := model.GetUserLogsPageList(userId, logType, startTimestamp, endTimestamp, modelName, tokenName, p+1, int64(common.ItemsPerPage))
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success":     true,
		"message":     "",
		"data":        page.Data,
		"total":       page.Total,
		"pages":       page.Pages,
		"currentPage": page.CurrentPage,
		"pageSize":    page.PageSize,
	})
}

func SearchAllLogs(c *gin.Context) {
	keyword := c.Query("keyword")
	logs, err := model.SearchAllLogs(keyword)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
}

func SearchUserLogs(c *gin.Context) {
	keyword := c.Query("keyword")
	userId := c.GetInt("id")
	logs, err := model.SearchUserLogs(userId, keyword)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data":    logs,
	})
}

func GetLogsStat(c *gin.Context) {
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	username := c.Query("username")
	modelName := c.Query("model_name")
	quotaNum := model.SumUsedQuota(logType, startTimestamp, endTimestamp, modelName, username, tokenName)
	//tokenNum := model.SumUsedToken(logType, startTimestamp, endTimestamp, modelName, username, "")
	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"quota": quotaNum,
			//"token": tokenNum,
		},
	})
}

func GetLogsSelfStat(c *gin.Context) {
	username := c.GetString("username")
	logType, _ := strconv.Atoi(c.Query("type"))
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	tokenName := c.Query("token_name")
	modelName := c.Query("model_name")
	quotaNum := model.SumUsedQuota(logType, startTimestamp, endTimestamp, modelName, username, tokenName)
	//tokenNum := model.SumUsedToken(logType, startTimestamp, endTimestamp, modelName, username, tokenName)
	c.JSON(200, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"quota": quotaNum,
			//"token": tokenNum,
		},
	})
}
