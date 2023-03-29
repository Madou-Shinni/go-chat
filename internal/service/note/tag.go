package note

import (
	"context"
	"errors"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type ArticleTagService struct {
	*repo.Source
}

func NewArticleTagService(source *repo.Source) *ArticleTagService {
	return &ArticleTagService{Source: source}
}

func (s *ArticleTagService) Create(ctx context.Context, uid int, tag string) (int, error) {
	data := &model.ArticleTag{
		UserId:  uid,
		TagName: tag,
		Sort:    1,
	}

	if err := s.Db().WithContext(ctx).Create(data).Error; err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (s *ArticleTagService) Update(ctx context.Context, uid int, tagId int, tag string) error {
	return s.Db().Table("article_tag").Where("id = ? and user_id = ?", tagId, uid).UpdateColumn("tag_name", tag).Error
}

func (s *ArticleTagService) Delete(ctx context.Context, uid int, tagId int) error {

	db := s.Db().WithContext(ctx)

	var num int64
	if err := db.Table("article").Where("user_id = ? and FIND_IN_SET(?,tags_id)", uid, tagId).Count(&num).Error; err != nil {
		return err
	}

	if num > 0 {
		return errors.New("标签已被使用不能删除")
	}

	return db.Delete(&model.ArticleTag{}, "id = ? and user_id = ?", tagId, uid).Error
}

func (s *ArticleTagService) List(ctx context.Context, uid int) ([]*model.TagItem, error) {

	db := s.Db().WithContext(ctx)

	items := make([]*model.TagItem, 0)

	err := db.Table("article_tag").Select("id", "tag_name").Where("user_id = ?", uid).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var num int64
		if err := db.Table("article").Where("user_id = ? and status = 1 and FIND_IN_SET(?,tags_id)", uid, item.Id).Count(&num).Error; err == nil {
			item.Count = int(num)
		}
	}

	return items, nil
}
