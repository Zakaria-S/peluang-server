package pricerange

import (
	"context"
	"peluang-server/domain"
)

type service struct {
	priceRangeRepo domain.PriceRangeRepository
}

func NewService(priceRangeRepo domain.PriceRangeRepository) domain.PriceRangeService {
	return &service{
		priceRangeRepo,
	}
}

func (s *service) GetAll() ([]domain.PriceRange, error) {
	priceranges, err := s.priceRangeRepo.FindAll()
	if err != nil {
		return []domain.PriceRange{}, err
	}
	return priceranges, nil
}

func (s *service) Create(priceRange *domain.PriceRange, ctx context.Context) (*domain.PriceRange, error) {

	if err := s.priceRangeRepo.Store(priceRange); err != nil {
		return &domain.PriceRange{}, err
	}
	return priceRange, nil
}

func (s *service) Update(priceRange *domain.PriceRange, ctx context.Context) (*domain.PriceRange, error) {
	panic("unimplemented")
}
func (s *service) Delete(priceRange *domain.PriceRange, ctx context.Context) error {
	panic("unimplemented")
}
