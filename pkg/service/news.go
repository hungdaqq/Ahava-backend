package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type NewsService interface {
	AddNews(n models.News) (models.News, error)
	UpdateNews(news_id uint, n models.News) (models.News, error)
	DeleteNews(news_id uint) error
	ListAllNews(limit, offset int) (models.ListNews, error)
	GetFeaturedNews() ([]models.News, error)
	GetNewsByID(news_id uint) (models.News, error)
}

type newsService struct {
	repository repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
	return &newsService{
		repository: repo,
	}
}

func (ns *newsService) AddNews(n models.News) (models.News, error) {
	return ns.repository.AddNews(n)
}

func (ns *newsService) UpdateNews(news_id uint, n models.News) (models.News, error) {
	return ns.repository.UpdateNews(news_id, n)
}

func (ns *newsService) DeleteNews(news_id uint) error {
	return ns.repository.DeleteNews(news_id)
}

func (ns *newsService) ListAllNews(limit, offset int) (models.ListNews, error) {
	return ns.repository.ListAllNews(limit, offset)
}

func (ns *newsService) GetFeaturedNews() ([]models.News, error) {
	return ns.repository.GetFeaturedNews()
}

func (ns *newsService) GetNewsByID(news_id uint) (models.News, error) {
	return ns.repository.GetNewsByID(news_id)
}
