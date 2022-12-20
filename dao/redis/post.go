package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"web_app/models"
)

func getIDsfromKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZREVRANGE 按分数从大到小查询
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 根据p.size, p.page, p.order从redis查询出post的id并返回
	// 1. 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2. 确定查询的索引起始点
	// 3. ZREVRANGE 按分数从大到小查询
	return getIDsfromKey(key, p.Page, p.Size)

}

// 按社区对帖子进行分数或时间的排序

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用 zinterstore 把分区的帖子set与帖子分数的 zset 生成一个新的zset
	// 针对新的zset按之前的逻辑(REVRANGE)取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	// 因为zinterstore比较重，所以利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) // zinterstore计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话直接根据key查询ids
	return getIDsfromKey(key, p.Page, p.Size)
}

// 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量， 查反对票的话只需把min的1改成-1
	//	v := client.ZCount(key, "1", "1").Val()
	//	fmt.Printf("---------------v是---------------%d", v, "\n")
	//	data = append(data, v)
	//}
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
