package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
	"web_app/models"
)

// CreatePostHandler 创建帖子
// @Summary 发布帖子
// @Description 根据用户传入的帖子信息以及请求上下文得到的用户id创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.Post false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {string} json "{"code":1000,"msg":"success","data":null}"
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	// c.ShouldBindJSON() // validator --> binding tag
	var err error
	p := new(models.Post)
	if err = c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c取到当前发请求的用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据传入的帖子id获取帖子详情
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "帖子ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.ApiPostDetail
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取路由参数(帖子id)
	postIDStr := c.Param("id")
	id, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id取出帖子数据
	// 2. 根据id获取社区详情
	data, err := logic.GetPostDetail(id)
	if err != nil {
		zap.L().Error("logic.PostDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 注意不要轻易把服务器报错暴露给用户
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 查询所有的帖子
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodePostNotExist)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListByTimeOrScoreHandler 升级版获取帖子列表的处理函数
// 根据前端传来的参数动态的获取帖子列表
// 按创建时间排序或者按照分数排序
// 1. 获取参数 2. 去redis查询id列表 3. 根据id列表区数据库查询帖子详细信息
func GetPostListByTimeOrScoreHandler(c *gin.Context) {
	// GET请求参数(query string)：/api/v1/postsByTimeOrScore?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page:  1, // 这些常量最好是改成配置文件里的参数
		Size:  10,
		Order: models.OrderScore,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListByTimeOrScoreHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 查询所有的帖子
	data, err := logic.GetPostListByTimeOrScoreNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodePostNotExist)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

//func GetCommunityPostListByTimeOrScoreHandler(c *gin.Context) {
//	// 根据请求参数得到社区id
//	// GET请求参数(query string)：/api/v1/CommunityPostsByTimeOrScore?community_id=1&post=1&size=10&order=time
//	p := new(models.ParamPostList)
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListByTimeOrScoreHandler with invalid param", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	// 查询所有帖子
//	data, err := logic.GetPostListByTimeOrScoreNew(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostListByTimeOrScore failed", zap.Error(err))
//		ResponseError(c, CodePostNotExist)
//		return
//	}
//	// 返回响应
//	ResponseSuccess(c, data)
//	// 返回响应
//}
