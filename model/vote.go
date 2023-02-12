package model

type VoteData struct {
	PostID int64 `json:"post_id,string"`
	VoteUp int   `json:"vote_up"`
}
