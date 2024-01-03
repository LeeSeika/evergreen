package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSec = 7 * 24 * 60 * 60 * 365
	valuePerVote = 432

	valuePerComment = 250
)

var (
	ErrorVoteTimeExpired = errors.New("vote time expired")
)

func VoteForPost(userID, postID string, voteValue float64) error {
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSec {
		return ErrorVoteTimeExpired
	}
	oldVoteValue := rdb.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	var voteDirection float64
	if voteValue == oldVoteValue {
		return nil
	}
	if voteValue > oldVoteValue {
		voteDirection = 1
	} else {
		voteDirection = -1
	}
	diff := math.Abs(oldVoteValue - voteValue)

	pipeline := rdb.TxPipeline()
	_, err := pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), voteDirection*diff*valuePerVote, postID).Result()
	if err != nil {
		return err
	}
	if voteValue == 0 {
		_, err = pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Result()
	} else {
		_, err = pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  voteValue,
			Member: userID,
		}).Result()
	}
	_, err = pipeline.Exec()

	return err
}

func AddCommentToPost(postID string) error {
	_, err := rdb.ZIncrBy(getRedisKey(KeyPostScoreZSet), valuePerComment, postID).Result()
	if err != nil {
		return err
	}
	return nil
}
