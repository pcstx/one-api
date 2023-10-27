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

type TokenOrder struct {
	OrderPrice float64 `json:"orderPrice"`
	UserId     int     `json:"userId"`
	OpenId     string  `json:"openId"`
	Token      string  `json:"token"`
}

type PerkAIOrder struct {
	OrderPrice   float64 `json:"orderPrice"`
	PayType      int     `json:"payType"`
	PerkAIUserId int     `json:"perkAIUserId"`
	Token        string  `json:"token"`
}

type PayData struct {
	OrderNumber string      `json:"orderNumber"`
	ImgUrl      string      `json:"imgUrl"`
	Form        string      `json:"form"`
	PayDataDto  interface{} `json:"payDataDto"`
	PayUrl      string      `json:"payUrl"`
}

type QueryOrder struct {
	OrderNumber string `json:"orderNumber"`
	Token       string `json:"token"`
}

type Order struct {
	Id          int64  `json:"id"`
	OrderNumber string `json:"orderNumber"`
	OrderPrice  int64  `json:"orderPrice"`
	OrderTime   string `json:"orderTime"`
	PayTime     string `json:"payTime"`
	PayStatus   int    `json:"payStatus"`
	OrderStatus int    `json:"orderStatus"`
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
	common.AddPushToken(c, token)
	common.SetSession[string](c, "pushToken", token)

	//根据token获取pushplus中用户详情
	userInfo, _ := GetMyInfo(token)

	//系统内的账号登录
	weChatLogin(c, userInfo)
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
	//未登录或者其他错误
	if res.Code != 200 {
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
	//未登录或者其他错误
	if res.Code != 200 {
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

func AuthLogin(c *gin.Context, userInfo *UserInfo) {
	wechatId := userInfo.OpenId

	user := model.User{
		WeChatId: wechatId,
	}
	if model.IsWeChatIdAlreadyTaken(wechatId) {
		err := user.FillUserByWeChatId()
		if err != nil {
			common.Error(c, err.Error())
			c.Abort()
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
				c.Abort()
				return
			}
		} else {
			common.Error(c, "管理员关闭了新用户注册")
			c.Abort()
			return
		}
	}

	if user.Status != common.UserStatusEnabled {
		common.Error(c, "用户已被封禁")
		c.Abort()
		return
	}
	saveLoginInfo(&user, c)
}

/*
登录操作
先调用pushplus登录接口
*/
func (con PushplusController) WechatLogout(c *gin.Context) {

	//调用pushplus接口
	loginOut(c)
	//删除cookie
	common.RemovePushToken(c)
	common.RemoveCookie(c, "session")
	//c.SetCookie("pushToken", "", -1, "/", common.PushPlusDomain, false, false)
	common.ClearSession(c)
	common.Success(c, nil)
}

func loginOut(c *gin.Context) {
	token := common.GetPushToken(c)
	//token := common.GetSession[string](c, "pushToken")

	fmt.Printf("退出时候获取token=%v\n", token)
	url := fmt.Sprintf("%s/customer/login/loginOut", common.PushPlusApiUrl)
	_, err := common.HttpGet[string](url, token)
	if err != nil {
		fmt.Println(err.Error())
	}
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

func perkTokenOrder(tokenOrder *TokenOrder) (string, error) {
	url := fmt.Sprintf("%s/customer/payPage/tokenOrder?orderPrice=%s&userId=%s&openId=%s", common.PushPlusApiUrl, tokenOrder.OrderPrice, tokenOrder.UserId, tokenOrder.OpenId)
	result, err := common.HttpGet[string](url, tokenOrder.Token)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return result.Data, err
}

func (con PushplusController) TokenOrder(c *gin.Context) {
	orderPriceStr := c.Query("orderPrice")

	// 将字符串转换为 float32
	orderPrice, err := strconv.ParseFloat(orderPriceStr, 32)
	if err != nil {
		fmt.Println("无法将字符串转换为浮点数:", err)
		return
	}

	//服务端校验
	if orderPrice <= 1 {
		common.Error(c, "充值金额最少1元")
		return
	}

	var tokenOrder TokenOrder

	//获取请求token

	token := common.GetPushToken(c)
	//token := common.GetSession[string](c, "pushToken")
	userId := common.GetSession[int](c, "id")

	user, err := model.GetUserById(userId, false)
	if err != nil {
		common.Error(c, err.Error())
		return
	}
	tokenOrder.OpenId = user.WeChatId
	tokenOrder.Token = token
	tokenOrder.UserId = userId
	tokenOrder.OrderPrice = orderPrice

	result, err := perkTokenOrder(&tokenOrder)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	common.Success(c, result)
	return
}

// 调用pushplus下单
func perkAIOrder(perkAIOrder *PerkAIOrder) (*PayData, error) {
	url := fmt.Sprintf("%s/customer/pay/perkAIOrder", common.PushPlusApiUrl)
	//perkAIOrder.PayType = 0
	result, err := common.HttpPost[PayData](url, perkAIOrder, perkAIOrder.Token)
	if err != nil {
		fmt.Println(err.Error())
		return &PayData{}, err
	}

	return &result.Data, err
}

func (con PushplusController) PerkAIOrder(c *gin.Context) {
	var perkAIOrderDto PerkAIOrder
	if err := c.ShouldBind(&perkAIOrderDto); err != nil {
		common.Error(c, "请求对象有误")
		return
	}

	//服务端校验
	if perkAIOrderDto.OrderPrice < 1 {
		common.Error(c, "充值金额最少1元")
		return
	}

	//获取请求token
	token := common.GetPushToken(c)
	fmt.Printf("token:%s", token)
	//token := common.GetSession[string](c, "pushToken")
	userId := common.GetSession[int](c, "id")
	perkAIOrderDto.Token = token
	perkAIOrderDto.PerkAIUserId = userId

	//1元等于5万token
	quota := perkAIOrderDto.OrderPrice * 50000

	result, err := perkAIOrder(&perkAIOrderDto)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	if len(result.OrderNumber) > 0 {
		//支付订单入库
		redemption := model.Redemption{
			UserId:      userId,
			Name:        "pushplus现金充值",
			Key:         result.OrderNumber,
			CreatedTime: common.GetTimestamp(),
			Quota:       int(quota),
		}
		err = redemption.Insert()
		if err != nil {
			common.Error(c, err.Error())
			return
		}
	} else {
		common.Error(c, "生成订单失败，请稍后重试")
		return
	}

	common.Success(c, result)
	return
}

// 主动查询
func queryOrder(queryOrder *QueryOrder) (*Order, error) {
	url := fmt.Sprintf("%s/customer/pay/queryOrder?orderNumber=%s", common.PushPlusApiUrl, queryOrder.OrderNumber)
	result, err := common.HttpGet[Order](url, queryOrder.Token)
	if err != nil {
		fmt.Println(err.Error())
		return &Order{}, err
	}

	return &result.Data, err
}

func (con PushplusController) QueryOrder(c *gin.Context) {
	orderNumber := c.Query("orderNumber")

	//获取请求token
	token := common.GetPushToken(c)
	//token := common.GetSession[string](c, "pushToken")
	userId := common.GetSession[int](c, "id")

	var queryOrderDto QueryOrder
	queryOrderDto.Token = token
	queryOrderDto.OrderNumber = orderNumber

	result, err := queryOrder(&queryOrderDto)
	if err != nil {
		common.Error(c, err.Error())
		return
	}

	//判断状态入库处理
	if result.OrderStatus != 0 {
		if result.OrderStatus == 1 {
			//支付成功
			_, err = model.UpdatePaySuccess(result.OrderNumber, userId)
			// if err != nil {
			// 	common.Error(c, err.Error())
			// 	return
			// }
		} else if result.OrderStatus == -1 {
			//支付取消
			_, err = model.UpdatePayCancel(result.OrderNumber, userId)
			// if err != nil {
			// 	common.Error(c, err.Error())
			// 	return
			// }
		}

	}

	//result.OrderStatus = 1

	common.Success(c, result)
	return
}

// 回调方法
func (con PushplusController) CallBackPay(c *gin.Context) {
	var payCallBack common.PayCallBack
	if err := c.ShouldBind(&payCallBack); err != nil {
		common.Error(c, "请求对象有误")
		return
	}

	//校验是否伪造
	key := common.Valid(payCallBack)

	if key != payCallBack.Key {
		fmt.Println("未通过校验")
		common.Error(c, "未通过请求校验")
		return
	}

	// 判断状态入库处理
	if payCallBack.OrderStatus != 0 {
		if payCallBack.OrderStatus == 1 {
			//支付成功
			_, err := model.UpdatePaySuccess(payCallBack.OrderNumber, payCallBack.UserId)
			if err != nil {
				common.Error(c, err.Error())
				return
			}
		} else if payCallBack.OrderStatus == -1 {
			//支付取消
			_, err := model.UpdatePayCancel(payCallBack.OrderNumber, payCallBack.UserId)
			if err != nil {
				common.Error(c, err.Error())
				return
			}
		}
	}
	common.Success(c, 1)
	return
}
