package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

type PayCallBack struct {
	UserId      int    `json:"userId"`
	OrderNumber string `json:"orderNumber"`
	OrderStatus int    `json:"orderStatus"`
	Key         string `json:key`
}

func Valid(payCallBack PayCallBack) string {
	valid := "salt=2ddfdca6f8c5b22d36ab21e6f8644650&UserId:" + strconv.Itoa(payCallBack.UserId) + "&OrderNumber:" + payCallBack.OrderNumber + "&OrderStatus:" + strconv.Itoa(payCallBack.OrderStatus)

	hasher := md5.New()

	// 将字符串写入哈希对象
	io.WriteString(hasher, valid)

	// 计算MD5哈希值
	hashedBytes := hasher.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashedString := fmt.Sprintf("%x", hashedBytes)

	return hashedString
}
