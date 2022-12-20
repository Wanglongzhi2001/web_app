package routes

import (
	"net/http"
	"web_app/controllers"
	"web_app/logger"
	"web_app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("api/v1")
	// 注册业务路由
	v1.POST("/SignUp", controllers.SignUpHandler)

	// 登录业务路由
	v1.POST("/login", controllers.LoginHandler)
	v1.Use(middleware.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		// 发布帖子
		v1.POST("/post", controllers.CreatePostHandler)
		// 点击帖子获取帖子详情
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts/", controllers.GetPostListHandler)
		// 根据时间或分数获取帖子列表
		v1.GET("/postsByTimeOrScore", controllers.GetPostListByTimeOrScoreHandler)
		// 在特定社区中根据时间或分数获取帖子列表
		v1.GET("/communityPostsByTimeOrScore", controllers.GetPostListByTimeOrScoreHandler)
		// 投票
		v1.POST("/vote", controllers.PostVoteHandler)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "欢迎光临")
	})
	r.GET("/ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是已登录的用户(判断请求头中是否有有效的JWT token)
		//isLogin := true
		//if isLogin {
		//	c.String(http.StatusOK, "pong")
		//} else {
		//	// 否则就返回请登录
		//	c.String(http.StatusOK, "请登录")
		//}
		// 认证的操作已经放到中间件里了
		c.String(http.StatusOK, "pong")
	})
	return r
}
