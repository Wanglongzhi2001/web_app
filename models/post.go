package models

import "time"

// 注意内存对齐

type Post struct {
	ID          int64     `json:"id,string" db:"post_id"` // json加上string是为了防止id值大于2^53 - 1时导致的js数字失真，int64最大值为2^63 - 1
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	*Post                               // 嵌入社区信息
	*CommunityDetail `json:"community"` // 使得返回的与community有关的字段会整齐的放在一起
}
