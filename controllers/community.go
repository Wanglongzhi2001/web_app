package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"web_app/logic"
)

// -----跟社区相关的-----
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 注意不要轻易把服务器报错暴露给用户
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name)以列表的形式返回
	// 1. 获取社区id
	communityIDStr := c.Param("id") // /community/:id c.Param获取路径参数
	id, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 根据id获取社区详情
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam) // 注意不要轻易把服务器报错暴露给用户
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}
