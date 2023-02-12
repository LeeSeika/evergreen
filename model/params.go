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
