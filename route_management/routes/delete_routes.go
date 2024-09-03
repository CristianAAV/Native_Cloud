package routes

import (
	"route_management/model"
	"route_management/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteRoute(c *gin.Context, db *gorm.DB, config utils.Config) {
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

	route := model.Route{ID: parsedId.String()}
	res := db.
		Delete(&route)

	if res.RowsAffected != 1 {
		c.Status(404)
		return
	}

	c.JSON(200, gin.H{"msg": "el trayecto fue eliminado"})
}
