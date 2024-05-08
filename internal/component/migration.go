package component

import (
	"peluang-server/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&domain.User{},
		&domain.UserOtp{},
		&domain.PriceRange{},
	)
}
