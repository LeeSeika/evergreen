package logic

import (
	"encoding/json"
	"evergreen/dao/redis"
	"evergreen/middleware/mq"
	"evergreen/model"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(userID int64, p *model.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", p.PostID), zap.Int8("direction", p.VoteDirection))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.VoteDirection))
}

func HandleAddComment() {
	msgs, err := mq.GetVoteConsumerMsg()
	if err != nil {
		zap.L().Error("get vote consumer msg failed", zap.Error(err))
		return
	}
	for m := range msgs {
		comment := model.Comment{}
		err := json.Unmarshal(m.Body, &comment)
		if err != nil {
			zap.L().Error("json unmarshal comment failed, going to deliver to dlx queue", zap.Error(err))
			m.Nack(false, false)
			continue
		}
		err = redis.AddCommentToPost(strconv.Itoa(int(comment.PostId)))
		if err != nil {
			zap.L().Error("add post to redis failed, going to deliver to dlx queue", zap.Error(err))
			m.Nack(false, false)
		}
	}
}
