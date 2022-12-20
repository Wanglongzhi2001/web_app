package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*

{
	"code":10001, // 程序中的错误码
	"msg": xx, // 提示信息
	"data": {}, // 数据
}
*/

type ResonseDate struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	//rd := &ResonseDate{
	//	Code: code,
	//	Msg:  code.Msg(),
	//	Data: nil,
	//}
	//c.JSON(http.StatusOK, rd)
	// 少写一个变量
	c.JSON(http.StatusOK, &ResonseDate{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// 自定义错误
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	//rd := &ResonseDate{
	//	Code: code,
	//	Msg:  code.Msg(),
	//	Data: nil,
	//}
	//c.JSON(http.StatusOK, rd)
	// 少写一个变量
	c.JSON(http.StatusOK, &ResonseDate{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	//rd := &ResonseDate{
	//	Code: CodeSuccess,
	//	Msg:  CodeSuccess.Msg(),
	//	Data: data,
	//}
	//c.JSON(http.StatusOK, rd)

	c.JSON(http.StatusOK, &ResonseDate{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
