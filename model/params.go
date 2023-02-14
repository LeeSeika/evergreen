package model

type ParamSignUp struct {
	Username   string `json:"user_name" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	PostID        string `json:"post_id" binding:"required"`
	VoteDirection int8   `json:"vote_direction" binding:"oneof=-1 0 1"`
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
