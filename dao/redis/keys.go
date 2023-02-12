package redis

const (
	KeyPrefix              = "evergreen:"
	KeyPostTimeZSet        = "post:time"
	KeyPostScoreZSet       = "post:score"
	KeyPostVotedZSetPrefix = "post:voted:"
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
