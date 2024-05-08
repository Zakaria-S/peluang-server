package domain

import (
	"context"
	"time"
)

type PriceRange struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Min       int       `json:"min" gorm:"unique"`
	Max       int       `json:"max" gorm:"unique"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PriceRangeRepository interface {
	FindAll() ([]PriceRange, error)
	//FindPriceRange(Min int, Max int)
	Store(priceRange *PriceRange) error
	Update(priceRange *PriceRange) error
	Delete(priceRange *PriceRange) error
}

type PriceRangeService interface {
	GetAll() ([]PriceRange, error)
	Create(priceRange *PriceRange, ctx context.Context) (*PriceRange, error)
	Update(priceRange *PriceRange, ctx context.Context) (*PriceRange, error)
	Delete(priceRange *PriceRange, ctx context.Context) error
}
