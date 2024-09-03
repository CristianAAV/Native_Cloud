package models

import (
	"gorm.io/gorm"
)

type Offer struct {
	gorm.Model
	ID          string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4();column:id"` // UUID v4
	PostId      string  `gorm:"type:uuid;not null;column:post_id"`                          // UUID
	UserId      string  `gorm:"type:uuid;not null;column:user_id"`                          // UUID
	Description string  `gorm:"type:varchar(140);not null;column:description"`              // Máximo 140 caracteres
	Size        string  `gorm:"type:varchar(10);not null;column:size"`                      // Valores posibles: LARGE, MEDIUM, SMALL
	Fragile     bool    `gorm:"default:false;column:fragile"`                               // Booleano
	Offer       float64 `gorm:"not null;column:offer"`                                      // Valor en dólares
}
