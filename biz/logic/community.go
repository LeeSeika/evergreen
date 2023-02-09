package logic

import (
	"evergreen/dao/mysql"
	"evergreen/model"
)

func GetCommunityList() ([]*model.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*model.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
