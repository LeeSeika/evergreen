package logic

import (
	"evergreen/dao/mysql"
	"evergreen/model"
	"evergreen/pkg/snowflake"
)

func CreatePost(p *model.Post) error {
	p.ID = snowflake.GenID()
	return mysql.CreatePost(p)
}

func GetPostDetailByID(postID int64) (*model.ApiPostDetail, error) {
	post, err := mysql.GetPostDetailByID(postID)
	if err != nil {
		return nil, err
	}
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		return nil, err
	}
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		return nil, err
	}
	postDetail := model.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return &postDetail, nil
}

func GetPostDetailList(page, size int64) ([]*model.ApiPostDetail, error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	apiPostDetailList := make([]*model.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			return nil, err
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			return nil, err
		}
		postDetail := model.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		apiPostDetailList = append(apiPostDetailList, &postDetail)
	}
	return apiPostDetailList, nil
}
