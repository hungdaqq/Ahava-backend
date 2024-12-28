package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type offerRepository struct {
	DB *gorm.DB
}

func NewOfferRepository(db *gorm.DB) *offerRepository {
	return &offerRepository{
		DB: db,
	}
}

type OfferRepository interface {
	AddNewOffer(model models.Offer) (offer models.Offer, err error)
	ExpireOffer(offer_id uint) (offer models.Offer, err error)
	UpdateOffer(offer_id uint, model models.Offer) (offer models.Offer, err error)
	FindOfferRate(product_id uint) (percentage uint, err error)
	GetOffers() ([]models.Offer, error)
}

func (repo *offerRepository) AddNewOffer(model models.Offer) (models.Offer, error) {

	var offer models.Offer

	if err := repo.DB.Exec("INSERT INTO offers (product_id,offer_rate,expire_at) values($1,$2,$3) RETURNING *",
		model.ProductID, model.OfferRate).Scan(&offer).Error; err != nil {
		return models.Offer{}, err
	}

	return offer, nil
}

func (repo *offerRepository) ExpireOffer(offer_id uint) (models.Offer, error) {

	var offer models.Offer

	result := repo.DB.Model(&domain.Offer{}).Where("id = ?", offer_id).Update("expire_at", time.Now()).Scan(&offer)

	if result.Error != nil {
		return models.Offer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Offer{}, errors.New("no matching offer found")
	}

	return offer, nil
}

func (repo *offerRepository) UpdateOffer(offer_id uint, model models.Offer) (models.Offer, error) {

	var offer models.Offer

	result := repo.DB.Model(&domain.Offer{}).Where("id = ?", offer_id).Updates(
		domain.Offer{
			OfferRate: model.OfferRate,
			ExpireAt:  model.ExpireAt,
		}).Scan(&offer)

	if result.Error != nil {
		return models.Offer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Offer{}, errors.New("no matching offer found")
	}

	return offer, nil
}

func (repo *offerRepository) FindOfferRate(product_id uint) (uint, error) {

	var percentage uint

	err := repo.DB.Raw(
		"SELECT offer_rate FROM offers WHERE product_id=$1 AND expire_at > NOW()",
		product_id,
	).Scan(&percentage).Error

	if err != nil {
		return 0, err
	}

	return percentage, nil
}

func (c *offerRepository) GetOffers() ([]models.Offer, error) {

	var offers []models.Offer

	err := c.DB.Raw("SELECT * FROM offers").Scan(&offers).Error
	if err != nil {
		return []models.Offer{}, err
	}

	return offers, nil
}
