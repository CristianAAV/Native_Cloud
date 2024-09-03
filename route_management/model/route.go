package model

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	ID                 string `gorm:"type:uuid;primary_key;default:uuid_generate_v4();column:id"`
	FlightId           string
	SourceAirportCode  string
	SourceCountry      string
	DestinyAirportCode string
	DestinyCountry     string
	BagCost            int
	PlannedStartDate   time.Time
	PlannedEndDate     time.Time
}

func (route *Route) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	id := uuid.NewString()
	route.ID = id
	return
}

func (r Route) ParseToDTO() gin.H {
	return gin.H{
		"id":                 r.ID,
		"flightId":           r.FlightId,
		"sourceAirportCode":  r.SourceAirportCode,
		"sourceCountry":      r.SourceCountry,
		"destinyAirportCode": r.DestinyAirportCode,
		"destinyCountry":     r.DestinyCountry,
		"bagCost":            r.BagCost,
		"plannedStartDate":   r.PlannedStartDate,
		"plannedEndDate":     r.PlannedEndDate,
		"createdAt":          r.CreatedAt,
	}
}

func ParseToDTO(r []Route) []gin.H {
	parsedRoutes := make([]gin.H, len(r))
	for i, route := range r {
		parsedRoutes[i] = route.ParseToDTO()
	}
	return parsedRoutes
}
