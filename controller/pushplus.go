package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"one-api/common"
	"time"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Data struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	QrCode    string `json:"qrCode"`
}

func getQrCodeUrl() (*Data, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/common/wechat/getQrcode", common.PushPlusApiUrl), nil)
	if err != nil {
		return &Data{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		common.SysLog(err.Error())
		return &Data{}, errors.New("请求失败，请稍后重试！")
	}
	defer res.Body.Close()
	var qrcodeResponse Response[Data]
	err = json.NewDecoder(res.Body).Decode(&qrcodeResponse)
	if err != nil {
		return &Data{}, err
	}

	return &qrcodeResponse.Data, nil
}

// 返回微信登录二维码
func QrCode(c *gin.Context) {
	qrCodeUrl, err := getQrCodeUrl()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    qrCodeUrl,
	})
	return
}

func getConfirmLogin(qrCode string) (string, error) {
	if qrCode == "" {
		return "", errors.New("无效的参数")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/common/wechat/confirmLogin?key=%s&code=1001", common.PushPlusApiUrl, qrCode), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		common.SysLog(err.Error())
		return "", errors.New("请求失败，请稍后重试！")
	}
	defer res.Body.Close()
	var confirmLoginResponse Response[string]
	err = json.NewDecoder(res.Body).Decode(&confirmLoginResponse)
	if err != nil {
		common.SysLog(err.Error())
		return "", err
	}
	fmt.Printf("返回内容：%v\n", &confirmLoginResponse)

	return confirmLoginResponse.Data, nil
}

func ConfirmLogin(c *gin.Context) {
	qrCode := c.Query("key")
	token, err := getConfirmLogin(qrCode)
	fmt.Printf("token: %s\n", token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "尚未登录",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    token,
	})
	return
}
