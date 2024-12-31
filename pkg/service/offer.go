package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type OfferService interface {
	AddNewOffer(model models.Offer) (models.Offer, error)
	ExpireOffer(product_id uint) (models.Offer, error)
	UpdateOffer(product_id uint, model models.Offer) (models.Offer, error)
	GetOffers() ([]models.Offer, error)
}

type offerService struct {
	repository repository.OfferRepository
}

func NewOfferService(repo repository.OfferRepository) OfferService {
	return &offerService{
		repository: repo,
	}
}

func (off *offerService) AddNewOffer(model models.Offer) (models.Offer, error) {

	result, err := off.repository.AddNewOffer(model)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (off *offerService) UpdateOffer(product_id uint, model models.Offer) (models.Offer, error) {

	result, err := off.repository.UpdateOffer(product_id, model)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (off *offerService) ExpireOffer(product_id uint) (models.Offer, error) {

	result, err := off.repository.ExpireOffer(product_id)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (o *offerService) GetOffers() ([]models.Offer, error) {

	offers, err := o.repository.GetOffers()
	if err != nil {
		return []models.Offer{}, err
	}

	return offers, nil
}
