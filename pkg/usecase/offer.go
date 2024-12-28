package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type offerUseCase struct {
	repository repository.OfferRepository
}

type OfferUseCase interface {
	AddNewOffer(model models.Offer) (models.Offer, error)
	ExpireOffer(offer_id uint) (models.Offer, error)
	UpdateOffer(offer_id uint, model models.Offer) (models.Offer, error)
	GetOffers() ([]models.Offer, error)
}

func NewOfferUseCase(repo repository.OfferRepository) *offerUseCase {
	return &offerUseCase{
		repository: repo,
	}
}

func (off *offerUseCase) AddNewOffer(model models.Offer) (models.Offer, error) {

	result, err := off.repository.AddNewOffer(model)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (off *offerUseCase) UpdateOffer(offer_id uint, model models.Offer) (models.Offer, error) {

	result, err := off.repository.UpdateOffer(offer_id, model)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (off *offerUseCase) ExpireOffer(offer_id uint) (models.Offer, error) {

	result, err := off.repository.ExpireOffer(offer_id)
	if err != nil {
		return models.Offer{}, err
	}

	return result, nil
}

func (o *offerUseCase) GetOffers() ([]models.Offer, error) {

	offers, err := o.repository.GetOffers()
	if err != nil {
		return []models.Offer{}, err
	}

	return offers, nil
}
