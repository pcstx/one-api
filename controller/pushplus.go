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

type Recharge struct {
	OutNumber string `json:"outNumber"`
	Point     int    `json:"point"`
	Token     string `json:"token"`
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
	session.Save()

	//根据token获取pushplus中用户详情
	userInfo, _ := GetMyInfo(token)

	//系统内的账号登录
	WeChatLogin(c, userInfo)
	return
}

// 实时获取用户资料
func getUserInfo(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/customer/user/userInfo", common.PushPlusApiUrl)
	res, err := common.HttpGet[UserInfo](url, token)
	if err != nil {
		common.SysLog(err.Error())
		return &UserInfo{}, err
	}

	return &res.Data, nil
}

// 获取用户资料（缓存中）
func GetMyInfo(token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/customer/user/myInfo", common.PushPlusApiUrl)
	res, err := common.HttpGet[UserInfo](url, token)
	if err != nil {
		common.SysLog(err.Error())
		return &UserInfo{}, err
	}

	return &res.Data, nil
}

func WeChatLogin(c *gin.Context, userInfo *UserInfo) {
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

	fmt.Printf("退出时候获取token=%v\n", token)
	url := fmt.Sprintf("%s/customer/login/loginOut", common.PushPlusApiUrl)
	common.HttpGet[string](url, token)
}

func perkAIRecharge(recharge *Recharge) (int, error) {
	url := fmt.Sprintf("%s/customer/userPoint/perkAIRecharge", common.PushPlusApiUrl)
	result, err := common.HttpPost[int](url, recharge, recharge.Token)
	if err == nil {
		if result.Data <= 0 {
			return result.Data, errors.New(result.Msg)
		}
	}

	return result.Data, err
}

/**
 * 充值接口
 */
func (con PushplusController) Recharge(c *gin.Context) {
	var recharge Recharge

	// 根据请求的Content-Type自动选择绑定器并将数据绑定到User结构体
	if err := c.ShouldBind(&recharge); err != nil {
		common.Error(c, "请求对象有误")
		return
	}

	//服务端校验
	if recharge.Point <= 0 {
		common.Error(c, "积分不能小于0")
		return
	}

	//获取请求token
	session := sessions.Default(c)
	token := session.Get("pushToken")
	userId := session.Get("id")
	recharge.Token = token.(string)

	//生成流水号
	recharge.OutNumber = common.GetUUID()
	//积分与额度兑换比例为：1积分=500额度
	quota := recharge.Point * 500
	//记录到数据库中
	redemption := model.Redemption{
		UserId:      userId.(int),
		Name:        "pushplus积分兑换",
		Key:         recharge.OutNumber,
		CreatedTime: common.GetTimestamp(),
		Quota:       quota,
	}
	err := redemption.Insert()
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	result, err := perkAIRecharge(&recharge)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	if result > 0 {
		//结果入库,更新表状态，user表中增加积分
		_, err := model.Recharge(redemption.Id, redemption.UserId)
		if err != nil {
			common.Error(c, err.Error())
			return
		}
		common.Success(c, quota)
		return
	} else {
		common.Error(c, "充值失败")
		return
	}
}
