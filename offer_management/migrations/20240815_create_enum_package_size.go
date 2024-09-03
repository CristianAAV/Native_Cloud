package migrations

import (
	"gorm.io/gorm"
)

func CreateEnumTypes(db *gorm.DB) error {
	return db.Exec("CREATE TYPE package_size AS ENUM ('LARGE', 'MEDIUM', 'SMALL');").Error
}

func CreateUUID4(db *gorm.DB) error {
	return db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
}

func DropEnumTypes(db *gorm.DB) error {
	return db.Exec("DROP TYPE IF EXISTS package_size;").Error
}

func DropUUID4(db *gorm.DB) error {
	return db.Exec("DROP EXTENSION IF EXISTS \"uuid-ossp\";").Error
}
