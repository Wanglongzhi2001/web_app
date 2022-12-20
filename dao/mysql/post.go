package mysql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"web_app/models"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post( 
	post_id, title, content, author_id, community_id)
    values(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostDetail 根据id查询单个帖子数据
func GetPostDetail(id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    		post_id, title, content, author_id, community_id 
			from post where post_id = ?`
	err = db.Get(post, sqlStr, id)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlStr := `select title, content, post_id, author_id, community_id from post ORDER BY create_time DESC limit ?,?` // 限制5条，不能一下全查出来
	postList = make([]*models.Post, 0, 5)                                                                             // 0 -> 长度, 5 -> 容量
	if err = db.Select(&postList, sqlStr, (page-1)*size, size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no post in db")
			err = nil
		}
	}
	return
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
			from post
			where post_id in (?)
			order by FIND_IN_SET(post_id, ?)` // FIND_IN_SET 按照给定的id的顺序返回
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

// 判断该社区是否有该帖子
func IsExistPostInCommunity(community_id int64, post_id int64) bool {
	sqlStr := `select 
    		post_id, title, content, author_id, community_id 
			from post where post_id = ?`
	post := new(models.Post)
	err := db.Get(post, sqlStr, post_id)
	if err == sql.ErrNoRows || post.CommunityID != community_id {
		return false
	}
	return true
}
