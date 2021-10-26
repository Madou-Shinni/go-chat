package service

import (
	"go-chat/app/model"
	"gorm.io/gorm"
)

type MemberItem struct {
	UserId   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Gender   int    `json:"gender"`
	Motto    string `json:"motto"`
	Leader   int    `json:"leader"`
	UserCard string `json:"user_card"`
}

type GroupMemberService struct {
	db *gorm.DB
}

func NewGroupMemberService(db *gorm.DB) *GroupMemberService {
	return &GroupMemberService{
		db: db,
	}
}

// isMember 判断用户是否是群成员
func (s *GroupMemberService) IsMember(groupId, userId int) bool {
	result := &model.GroupMember{}

	count := s.db.Select("id").
		Where("group_id = ? and user_id = ? and is_quit = ?", groupId, userId, 0).Unscoped().
		First(result).RowsAffected

	return count != 0
}

// GetMemberIds 获取所有群成员用户ID
func (s *GroupMemberService) GetMemberIds(groupId int) []int {
	var ids []int

	_ = s.db.Model(&model.GroupMember{}).Select("user_id").Where("group_id = ? and is_quit = ?", groupId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetMemberIds 获取所有群成员ID
func (s *GroupMemberService) GetUserGroupIds(userId int) []int {
	var ids []int

	_ = s.db.Model(&model.GroupMember{}).Select("id").Where("user_id = ? and is_quit = ?", userId, 0).Unscoped().Scan(&ids)

	return ids
}

// GetGroupMembers 获取群组成员列表
func (s *GroupMemberService) GetGroupMembers(groupId int) []*MemberItem {
	var items []*MemberItem

	fields := []string{
		"lar_group_member.leader",
		"lar_group_member.user_card",
		"lar_group_member.user_id",
		"lar_users.avatar",
		"lar_users.nickname",
		"lar_users.gender",
		"lar_users.motto",
	}

	s.db.Table("lar_group_member").
		Select(fields).
		Joins("left join lar_users on lar_users.id = lar_group_member.user_id").
		Where("lar_group_member.group_id = ? and lar_group_member.is_quit = ?", groupId, 0).
		Order("lar_group_member.leader desc").
		Unscoped().
		Scan(&items)

	return items
}

// GetMemberRemarks 获取指定群成员的备注信息
func (s *GroupMemberService) GetMemberRemarks(groupId int, userId int) string {
	var remarks string

	s.db.Model(&model.GroupMember{}).
		Select("user_card").
		Where("group_id = ? and user_id = ?", groupId, userId).
		Unscoped().
		Scan(&remarks)

	return remarks
}
