package controller

import (
	"errors"
	"fmt"
	"one-api/common"
	"one-api/model"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type PushplusController struct{}

type Data struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	QrCode    string `json:"qrCode"`
}

type UserInfo struct {
	Id          int64   `json:"id"`
	OpenId      string  `json:"openId"`
	UnionId     string  `json:"unionId"`
	NickName    string  `json:"nickName"`
	HeadImgUrl  string  `json:"headImgUrl"`
	PhoneNumber string  `json:"phoneNumber"`
	Email       string  `json:"email"`
	EmailStatus int     `json:"emailStatus"`
	Points      float32 `json:"points"`
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
func (con PushplusController) QrCode(c *gin.Context) {
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

func (con PushplusController) ConfirmLogin(c *gin.Context) {
	qrCode := c.Query("key")
	token, err := getConfirmLogin(qrCode)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	//设置cookie，实现单点登录
	c.SetCookie("pushToken", token, 7*3600*24, "/", common.PushPlusDomain, false, false)
	session := sessions.Default(c)
	session.Set("pushToken", token)

	//根据token获取pushplus中用户详情
	userInfo, _ := getUserInfo(token)

	//系统内的账号登录
	weChatLogin(c, userInfo)
	return
}

func getUserInfo(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/customer/user/userInfo", common.PushPlusApiUrl)
	res, err := common.HttpGet[UserInfo](url, token)
	if err != nil {
		common.SysLog(err.Error())
		return &UserInfo{}, err
	}

	return &res.Data, nil
}

func weChatLogin(c *gin.Context, userInfo *UserInfo) {
	wechatId := userInfo.OpenId

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
			displayName := userInfo.NickName
			if len(displayName) <= 0 {
				displayName = "微信用户"
			}

			user.Username = "wechat_" + strconv.Itoa(model.GetMaxUserId()+1)
			user.DisplayName = displayName
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

/*
登录操作
先调用pushplus登录接口
*/
func (con PushplusController) WechatLogout(c *gin.Context) {

	//调用pushplus接口
	loginOut(c)
	//删除cookie
	c.SetCookie("pushToken", "", -1, "/", common.PushPlusDomain, false, false)
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		common.Error(c, err.Error())
		return
	}
	common.Success(c, nil)
}

func loginOut(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("pushToken")

	url := fmt.Sprintf("%s/customer/login/loginOut", common.PushPlusApiUrl)
	common.HttpGet[string](url, token)
}
