package model

const (
	DefaultPostPageValue = 1
	DefaultPostSizeValue = 10

	DefaultCommentPageValue = 1
	DefaultCommentSizeValue = 15
)

const (
	OrderByTime  = "time"
	OrderByScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	PostID        string `json:"post_id" binding:"required"`
	VoteDirection int8   `json:"direction" binding:"oneof=-1 0 1"`
}

type ParamPostListInOrder struct {
	Page  int64  `form:"page"`
	Size  int64  `form:"size"`
	Order string `form:"order"`
}

type ParamCommunityPostList struct {
	*ParamPostListInOrder
	CommunityID int64 `json:"community_id"`
}

type ParamCommentsInOrder struct {
	Page  int64  `form:"page"`
	Size  int64  `form:"size"`
	Order string `form:"order"`
}

type ParamComments struct {
	*ParamCommentsInOrder
	PostID int64 `json:"post_id"`
}
