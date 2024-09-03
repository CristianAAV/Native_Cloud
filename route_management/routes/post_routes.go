package routes

import (
	"fmt"
	"route_management/model"
	"route_management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRoute(c *gin.Context, db *gorm.DB, config utils.Config) {

	isValid := config.Authenticator.ValidateAuth(c, config)
	if !isValid {
		return
	}

	var trayecto model.Route
	if err := c.BindJSON(&trayecto); err != nil {
		fmt.Printf("Failed parsing body to trayecto %v\n", c.Request.Body)
		fmt.Printf("An error ocurred: %v \n", err)
		c.Status(400)
		return
	}

	err := utils.ValidateTrayecto(trayecto, db, c)
	if err != nil {
		fmt.Printf("An error ocurred: %v \n", err)
		return
	}

	createdTrayecto := trayecto
	db.Create(&createdTrayecto)

	if createdTrayecto.ID == "" {
		c.Status(500)
	}

	currentTime := time.Now()
	c.JSON(201, gin.H{"id": createdTrayecto.ID, "createdAt": currentTime})
}
