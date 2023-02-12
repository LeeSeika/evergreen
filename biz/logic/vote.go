package logic

import (
	"evergreen/dao/redis"
	"evergreen/model"
	"strconv"

	"go.uber.org/zap"
)

func VoteForPost(userID int64, p *model.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userID", userID), zap.String("postID", p.PostID), zap.Int8("direction", p.VoteDirection))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.VoteDirection))
}
