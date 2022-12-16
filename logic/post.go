package logic

import (
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 根据model.post生成post_id, author_id, title ...
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID)
	return
}

func GetPostDetail(id int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们想要的接口想要的数据
	post, err := mysql.GetPostDetail(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetail failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	author, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUsernameById failed", zap.Error(err))
		return
	}
	// 根据社区id查询社区名称
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		Post:            post,
		AuthorName:      author.Username,
		CommunityDetail: communityDetail,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (postDetailList []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}
	postDetailList = make([]*models.ApiPostDetail, 0, len(posts))

	var (
		author          *models.User
		communityDetail *models.CommunityDetail
	)
	for _, post := range posts {

		author, err = mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区名称
		communityDetail, err = mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			continue
		}
		postDetailList = append(postDetailList, &models.ApiPostDetail{
			AuthorName:      author.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		})
	}
	return
}

func GetPostListByTimeOrScore(p *models.ParamPostList) (postDetailList []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostListByTimeOrScore return 0 data")
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 根据id去mysql中查找帖子详细信息，返回的数据要按照给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	postDetailList = make([]*models.ApiPostDetail, 0, len(posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	var (
		author          *models.User
		communityDetail *models.CommunityDetail
	)
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		author, err = mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区名称
		communityDetail, err = mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Error(err))
			continue
		}
		postDetailList = append(postDetailList, &models.ApiPostDetail{
			AuthorName:      author.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: communityDetail,
		})
	}
	return

}
