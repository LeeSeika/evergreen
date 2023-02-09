package mysql

import (
	"database/sql"
	"evergreen/model"

	"go.uber.org/zap"
)

func GetCommunityList() ([]*model.Community, error) {
	sqlStr := "select community_id, community_name from community"
	var communityList []*model.Community
	err := db.Select(&communityList, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
		} else {
			zap.L().Error("select community list error:", zap.Error(err))
		}
		return nil, err
	}
	return communityList, nil
}

func GetCommunityDetailByID(id int64) (*model.CommunityDetail, error) {
	var communityDetail model.CommunityDetail
	sqlStr := "select community_id, community_name, introduction, create_time from community where community_id = ?"
	err := db.Get(&communityDetail, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
		return nil, err
	}
	return &communityDetail, nil
}
