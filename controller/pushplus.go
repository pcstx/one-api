package controller

import (
	"errors"
	"fmt"
	"one-api/common"
	"one-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Data struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	QrCode    string `json:"qrCode"`
}

func getQrCodeUrl() (*Data, error) {

	url := fmt.Sprintf("%s/common/wechat/getQrcode", common.PushPlusApiUrl)
	data, err := common.HttpGet[Data](url)
	if err != nil {
		return &Data{}, err
	}
	return &data.Data, nil
}

// 返回微信登录二维码
func QrCode(c *gin.Context) {
	qrCodeUrl, err := getQrCodeUrl()
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	common.Success(c, qrCodeUrl)
	return
}

func getConfirmLogin(qrCode string) (string, error) {
	if qrCode == "" {
		return "", errors.New("无效的参数")
	}

	url := fmt.Sprintf("%s/common/wechat/confirmLogin?key=%s&code=1001", common.PushPlusApiUrl, qrCode)
	confirmLoginResponse, err := common.HttpGet[string](url)
	if err != nil {
		common.SysLog(err.Error())
		return "", err
	}
	if confirmLoginResponse.Code == 500 {
		return "", errors.New(confirmLoginResponse.Msg)
	}
	return confirmLoginResponse.Data, nil
}

func ConfirmLogin(c *gin.Context) {
	qrCode := c.Query("key")
	token, err := getConfirmLogin(qrCode)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	//设置cookie，实现单点登录
	c.SetCookie("pushToken", token, 7*3600*24, "/", common.PushPlusDomain, false, false)

	//系统内的账号登录
	weChatLogin(c, token)
	return
}

func weChatLogin(c *gin.Context, wechatId string) {
	user := model.User{
		WeChatId: wechatId,
	}
	if model.IsWeChatIdAlreadyTaken(wechatId) {
		err := user.FillUserByWeChatId()
		if err != nil {
			common.Error(c, err.Error())
			return
		}
	} else {
		if common.RegisterEnabled {
			user.Username = "wechat_" + strconv.Itoa(model.GetMaxUserId()+1)
			user.DisplayName = "WeChat User"
			user.Role = common.RoleCommonUser
			user.Status = common.UserStatusEnabled

			if err := user.Insert(0); err != nil {
				common.Error(c, err.Error())
				return
			}
		} else {
			common.Error(c, "管理员关闭了新用户注册")
			return
		}
	}

	if user.Status != common.UserStatusEnabled {
		common.Error(c, "用户已被封禁")
		return
	}
	setupLogin(&user, c)
}
