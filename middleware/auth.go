package middleware

import (
	"fmt"
	"strings"
	"time"
	"web_app/controllers"
	"web_app/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL
		// 这里假设Token放在Header的Authorization中， 并使用Bearer开头
		// Authorization：Bearer xxxxxx.xxx.xxx
		// 这里的具体实现方式要根据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2003,
			//	"msg":  "请求头中的auth为空",
			//})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2004,
			//	"msg":  "请求头中auth格式有误",
			//})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			//c.JSON(http.StatusOK, gin.H{
			//	"code": 2005,
			//	"msg":  "无效的token",
			//})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		// c.Set("userID", mc.UserID) // 不要在代码里写这种莫名其妙的字符串写死，用一个常量代替
		c.Set(controllers.ContextUserID, mc.UserID)
		c.Next() // 后续的处理函数可以通过c.Get("username")来获取当前请求的用户信息
	}
}

// 请求时间的中间件， 虽然zap日志库已经封装了请求时间而且比你的好哈哈
func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		fmt.Printf("the request URL %s cost %v\n", url, costTime)
	}
}
