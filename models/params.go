package models

// 定义请求的参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
	QueueName  = "createPostQueue"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`              // 帖子id(让前端传string，我们再转回int64)
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1)，反对票(-1)，取消投票(0)（这里应为必填字段，但是如果加上required就会读取不到0,因为validator库自动过滤false，0这些以为你没填
}

// ParamPOstList 获取帖子列表query string参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空，如果为空则默认查询所有帖子，否则查询某个社区的帖子
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

type ParamCommunityPostList struct {
	*ParamPostList
}
