package repository

import (
	"ahava/pkg/domain"
	errors "ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"
	"time"

	"gorm.io/gorm"
)

type OfferRepository interface {
	AddNewOffer(model models.Offer) (offer models.Offer, err error)
	ExpireOffer(product_id uint) (offer models.Offer, err error)
	UpdateOffer(product_id uint, model models.Offer) (offer models.Offer, err error)
	FindOfferRate(product_id uint) (percentage uint, err error)
	GetOffers() ([]models.Offer, error)
}

type offerRepository struct {
	DB *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{
		DB: db,
	}
}

func (r *offerRepository) AddNewOffer(model models.Offer) (models.Offer, error) {
	var offer models.Offer

	// Set default expiration if not provided
	if model.ExpireAt.IsZero() {
		model.ExpireAt = time.Now().AddDate(0, 1, 0)
	}

	err := r.DB.Raw(
		"INSERT INTO offers (product_id, offer_rate, expire_at) VALUES ($1, $2, $3) RETURNING *",
		model.ProductID, model.OfferRate, model.ExpireAt,
	).Scan(&offer).Error

	if err != nil {
		return models.Offer{}, err
	}

	return offer, nil
}

func (r *offerRepository) ExpireOffer(product_id uint) (models.Offer, error) {

	var offer models.Offer

	result := r.DB.
		Model(&domain.Offer{}).
		Where("product_id = ?", product_id).
		Update("expire_at", time.Now()).
		Scan(&offer)

	if result.Error != nil {
		return models.Offer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Offer{}, errors.ErrEntityNotFound
	}

	return offer, nil
}

func (r *offerRepository) UpdateOffer(product_id uint, model models.Offer) (models.Offer, error) {

	var offer models.Offer

	result := r.DB.
		Model(&domain.Offer{}).
		Where("product_id = ?", product_id).
		Updates(domain.Offer{
			OfferRate: model.OfferRate,
			ExpireAt:  model.ExpireAt,
		}).
		Scan(&offer)

	if result.Error != nil {
		return models.Offer{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Offer{}, errors.ErrEntityNotFound
	}

	return offer, nil
}

func (r *offerRepository) FindOfferRate(product_id uint) (uint, error) {

	var percentage uint

	err := r.DB.Raw(
		"SELECT offer_rate FROM offers WHERE product_id=$1 AND expire_at > NOW()",
		product_id,
	).Scan(&percentage).Error

	if err != nil {
		return 0, err
	}

	return percentage, nil
}

func (r *offerRepository) GetOffers() ([]models.Offer, error) {

	var offers []models.Offer

	err := r.DB.Raw("SELECT * FROM offers").Scan(&offers).Error
	if err != nil {
		return []models.Offer{}, err
	}

	return offers, nil
}
