package routes

import (
	"route_management/model"
	"route_management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoutes(c *gin.Context, db *gorm.DB, config utils.Config) {

	isValid := config.Authenticator.ValidateAuth(c, config)
	if !isValid {
		return
	}

	paramFlightId := c.Query("flight")
	paramFlightIdExists := paramFlightId != ""

	trayectos := []model.Route{}
	if paramFlightIdExists {
		db.
			Where(model.Route{FlightId: paramFlightId}).
			Find(&trayectos)
	} else {
		db.
			Find(&trayectos)
	}
	c.JSON(200, model.ParseToDTO(trayectos))
}

func GetRoute(c *gin.Context, db *gorm.DB, config utils.Config) {
	isValid := config.Authenticator.ValidateAuth(c, config)
	if !isValid {
		return
	}

	id := c.Params.ByName("id")
	parsedId, isValidId := utils.IsValidUUID(id)
	if !isValidId {
		c.Status(400)
		return
	}

	route := model.Route{}
	db.First(&route, parsedId)
	if route.ID == "" {
		c.Status(404)
		return
	}
	c.JSON(200, route.ParseToDTO())
}
