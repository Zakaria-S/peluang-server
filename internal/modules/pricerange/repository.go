package pricerange

import (
	"peluang-server/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(con *gorm.DB) domain.PriceRangeRepository {
	return &repository{
		db: con,
	}
}

func (r *repository) FindAll() ([]domain.PriceRange, error) {
	var pricerange []domain.PriceRange
	result := r.db.Find(&pricerange)
	if result.Error != nil {
		return []domain.PriceRange{}, result.Error
	}
	return pricerange, nil
}

/*
	func (r *repository) FindPriceRange(Min int, Max int) (*domain.PriceRange, error) {
		tx := r.db.Where("min = ?", &Min).First()
	}
*/
func (r *repository) Store(priceRange *domain.PriceRange) error {
	if tx := r.db.Create(&priceRange); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *repository) Update(priceRange *domain.PriceRange) error {
	if tx := r.db.Save(&priceRange); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *repository) Delete(priceRange *domain.PriceRange) error {
	if tx := r.db.Delete(&priceRange); tx.Error != nil {
		return tx.Error
	}
	return nil
}
