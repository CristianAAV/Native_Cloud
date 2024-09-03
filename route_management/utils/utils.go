package utils

import (
	"fmt"
	"route_management/model"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	User_url      string
	Authenticator IAuthenticator
}

func IsValidUUID(u string) (uuid.UUID, bool) {
	parsedId, err := uuid.Parse(u)
	return parsedId, err == nil
}

func ValidateTrayecto(trayecto model.Route, db *gorm.DB, c *gin.Context) error {
	if trayecto.DestinyAirportCode == "" ||
		trayecto.DestinyCountry == "" ||
		trayecto.FlightId == "" ||
		trayecto.PlannedEndDate.IsZero() ||
		trayecto.PlannedStartDate.IsZero() ||
		trayecto.SourceAirportCode == "" ||
		trayecto.SourceCountry == "" {
		c.Status(400)
		return fmt.Errorf("invalid trayecto")
	}

	if trayecto.PlannedEndDate.Before(time.Now()) || trayecto.PlannedStartDate.Before(time.Now()) {
		c.JSON(412, gin.H{"msg": "Las fechas del trayecto no son válidas"})
		return fmt.Errorf("invalid dates")
	}

	if trayecto.PlannedStartDate.After(trayecto.PlannedEndDate) {
		c.JSON(412, gin.H{"msg": "Las fechas del trayecto no son válidas"})
		return fmt.Errorf("invalid dates")
	}

	var flightIdExists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM routes WHERE flight_id = ?) AS found",
		trayecto.FlightId).Scan(&flightIdExists)
	if flightIdExists {
		c.Status(412)
		return fmt.Errorf("FlightId %s exists", trayecto.FlightId)
	}
	return nil
}
