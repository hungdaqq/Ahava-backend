package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type NewsRepository interface {
	AddNews(n models.News) (models.News, error)
	UpdateNews(news_id uint, n models.News) (models.News, error)
	DeleteNews(news_id uint) error
	ListAllNews(limit, offset int) (models.ListNews, error)
	GetFeaturedNews() ([]models.News, error)
	GetNewsByID(news_id uint) (models.News, error)
}

type newsRepository struct {
	DB *gorm.DB
}

func NewNewsRepository(DB *gorm.DB) NewsRepository {
	return &newsRepository{
		DB: DB,
	}
}

func (r *newsRepository) AddNews(n models.News) (models.News, error) {

	var news models.News

	if err := r.DB.Create(&domain.News{
		Title:          n.Title,
		Description:    n.Description,
		Content:        n.Content,
		DefaultImage:   n.DefaultImage,
		IsFeatured:     &n.IsFeatured,
		IsHomepage:     &n.IsHomepage,
		IsDisplay:      &n.IsDisplay,
		Category:       n.Category,
		TitleSEO:       n.TitleSEO,
		DescriptionSEO: n.DescriptionSEO,
		LinkSEO:        n.LinkSEO,
	}).Scan(&news).Error; err != nil {
		return news, err
	}

	return news, nil
}

func (r *newsRepository) UpdateNews(news_id uint, n models.News) (models.News, error) {

	var news models.News

	if err := r.DB.Model(&domain.News{}).Where("id = ?", news_id).Updates(domain.News{
		Title:          n.Title,
		Description:    n.Description,
		Content:        n.Content,
		DefaultImage:   n.DefaultImage,
		IsFeatured:     &n.IsFeatured,
		IsHomepage:     &n.IsHomepage,
		IsDisplay:      &n.IsDisplay,
		Category:       n.Category,
		TitleSEO:       n.TitleSEO,
		DescriptionSEO: n.DescriptionSEO,
	}).Scan(&news).Error; err != nil {
		return news, err
	}

	return news, nil
}

func (r *newsRepository) DeleteNews(news_id uint) error {

	if err := r.DB.Where("id = ?", news_id).Delete(&domain.News{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *newsRepository) ListAllNews(limit, offset int) (models.ListNews, error) {
	// Define the list of users
	var news []models.News
	var total int64
	// Define the query
	query := r.DB.Model(&domain.News{})
	if err := query.Count(&total).Error; err != nil {
		return models.ListNews{}, err
	}
	if err := query.Offset(offset).Limit(limit).Find(&news).Error; err != nil {
		return models.ListNews{}, err
	}
	// Return the list of users
	return models.ListNews{
		News:   news,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (r *newsRepository) GetFeaturedNews() ([]models.News, error) {
	var news []models.News
	if err := r.DB.Where("is_featured = ?", true).Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) GetNewsByID(news_id uint) (models.News, error) {
	var news models.News
	if err := r.DB.Where("id = ?", news_id).Find(&news).Error; err != nil {
		return models.News{}, err
	}
	return news, nil
}
