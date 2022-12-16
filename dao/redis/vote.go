package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

// 本项目使用简化版的投票分数， 投一票就加432分
/* 投票的几种情况
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票 --> 更新分数和投票纪录 差值的绝对值：1  +432
	2. 之前投反对票，现在改投赞成票 --> 更新分数和投票纪录 差值的绝对值：2  +432*2
direction=0时，有两种情况：
	1. 之前投过赞成票，现在取消投票 --> 更新分数和投票纪录 差值的绝对值：1  -432
	2. 之前投过反对票，现在取消投票 --> 更新分数和投票纪录 差值的绝对值：1  +432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票 --> 更新分数和投票纪录 差值的绝对值：1  -432
	2. 之前投赞成票，现在改投反对票 --> 更新分数和投票纪录 差值的绝对值：2  -432*2
投票的限制：
每个帖子自发表之日起，仅一个星期之内允许用户投票
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPrefix
*/

const (
	oneWeekInseconds = 7 * 24 * 3600
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不能重复投票")
)

func CreatePost(postID int64) error {
	pipeline := client.TxPipeline() // 初始化帖子时间和帖子分数应放在一个事务中
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, directionValue float64) error { // 因为go-redis设计有序集合的分数时用的是float64
	// 1. 判断投票的限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInseconds {
		return ErrVoteTimeExpire
	}
	// 2和3应放在一个事务中

	// 2. 更新分数
	// 先查当前用户给当前帖子的投票记录
	oldValue := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	var op float64
	if oldValue == directionValue {
		return ErrVoteRepeated
	}
	if directionValue > oldValue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(oldValue - directionValue) // 计算两次的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录该用户为该帖子投票的数据
	if directionValue == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  directionValue, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err

}
