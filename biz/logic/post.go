package logic

import (
	"evergreen/dao/mysql"
	"evergreen/dao/redis"
	"evergreen/model"
	"evergreen/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *model.Post) error {
	p.ID = snowflake.GenID()
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.CommunityID, p.ID)
	return err
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

func GetPostListInOrder(p *model.ParamPostListInOrder) ([]*model.ApiPostDetail, error) {
	ids, err := redis.GetPostListInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("there is no post for the given page and size", zap.Int64("page", p.Page), zap.Int64("size", p.Size))
		return nil, nil
	}
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	voteNumber, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	apiPostDetailList := make([]*model.ApiPostDetail, 0, len(postList))
	for idx, post := range postList {
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
			VoteNumbers:     voteNumber[idx],
			Post:            post,
			CommunityDetail: community,
		}
		apiPostDetailList = append(apiPostDetailList, &postDetail)
	}
	return apiPostDetailList, nil
}

func GetCommunityPostList(p *model.ParamCommunityPostList) ([]*model.ApiPostDetail, error) {
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("there is no post for the given community and page and size", zap.Int64("community_id", p.CommunityID), zap.Int64("page", p.Page), zap.Int64("size", p.Size))
		return nil, nil
	}
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	VoteNumber, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	apiPostDetailList := make([]*model.ApiPostDetail, 0, len(postList))
	for idx, post := range postList {
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
			VoteNumbers:     VoteNumber[idx],
			Post:            post,
			CommunityDetail: community,
		}
		apiPostDetailList = append(apiPostDetailList, &postDetail)
	}
	return apiPostDetailList, nil
}
