package redis

import (
	"evergreen/model"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetPostListInOrder(p *model.ParamPostListInOrder) ([]string, error) {
	var key string
	if p.Order == model.OrderByTime {
		key = getRedisKey(KeyPostTimeZSet)
	} else {
		key = getRedisKey(KeyPostScoreZSet)
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostVoteData(ids []string) ([]int64, error) {
	var data []int64
	pipeline := rdb.TxPipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1").Val()
	}
	exec, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range exec {
		intCmd := cmder.(*redis.IntCmd)
		data = append(data, intCmd.Val())
	}
	return data, nil
}

func GetCommunityPostIDsInOrder(p *model.ParamCommunityPostList) ([]string, error) {
	var orderKey string
	if p.Order == model.OrderByScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	} else {
		orderKey = getRedisKey(KeyPostTimeZSet)
	}
	communityKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + ":" + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(key).Val() < 1 {
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, communityKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return rdb.ZRevRange(key, start, end).Result()
}

func CreatePost(communityID, postID int64) error {
	pipeline := rdb.TxPipeline()

	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()

	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	}).Result()

	communityIDStr := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(communityIDStr, postID)
	_, err := pipeline.Exec()
	return err
}
