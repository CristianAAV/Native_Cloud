package main

import (
	"fmt"
	"net/http"
	"os"
	"route_management/model"
	"route_management/routes"
	"route_management/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB, config utils.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/routes", func(c *gin.Context) {
		routes.GetRoutes(c, db, config)
	})
	router.GET("/routes/:id", func(c *gin.Context) {
		routes.GetRoute(c, db, config)
	})

	router.POST("/routes", func(c *gin.Context) {
		routes.PostRoute(c, db, config)
	})

	router.DELETE("/routes/:id", func(c *gin.Context) {
		routes.DeleteRoute(c, db, config)
	})

	// health
	router.GET("/routes/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/routes/reset", func(c *gin.Context) {
		res := db.Exec("DELETE FROM routes")
		if res.Error != nil {
			fmt.Printf("Error: %v\n", res.Error)
		}
		c.JSON(200, gin.H{"msg": "todos los datos fueron eliminados"})
	})
	return router
}

func CreateUUID4(db *gorm.DB) error {
	return db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
}

func Start() *gin.Engine {
	var (
		postgresql_host      = os.Getenv("DB_HOST")
		postgresql_password  = os.Getenv("DB_PASSWORD")
		postgresql_user      = os.Getenv("DB_USER")
		postgresql_port, err = strconv.Atoi(os.Getenv("DB_PORT"))
		postgresql_dbname    = os.Getenv("DB_NAME")
		user_url             = os.Getenv("USERS_PATH")
	)

	config := utils.Config{
		User_url:      user_url,
		Authenticator: utils.GetAuthenticator(),
	}

	if err != nil {
		postgresql_port = 5432
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresql_host, postgresql_port, postgresql_user, postgresql_password, postgresql_dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{}) // sql.Open("postgres", psqlInfo)

	if err := db.Migrator().DropTable(&model.Route{}); err != nil {
		fmt.Printf("Error deleting Routes table: %v", err)
	}

	if err != nil {
		fmt.Printf("Failed to connect to db: " + err.Error())
	}
	err = CreateUUID4(db)
	if err != nil {
		fmt.Printf("Failed to create uuid: " + err.Error())
	}

	err = db.AutoMigrate(model.Route{})
	if err != nil {
		fmt.Printf("Failed to automigrate: " + err.Error())
	}

	router := setupRouter(db, config)
	return router
}

func main() {
	router := Start()
	err := router.Run(":3002")
	if err != nil {
		fmt.Printf("Failed to start router: " + err.Error())
	}
}
