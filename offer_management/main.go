package main

import (
	"log"
	"net/http"
	"offer_management/db"
	"offer_management/migrations"
	"offer_management/models"
	"offer_management/routes"

	"github.com/gorilla/mux"
)

// SetupRouter configura y devuelve el router
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Ruta para verificar la salud del servicio
	router.HandleFunc("/offers/ping", routes.PingHandler).Methods("GET")

	// Ruta para restablecer la base de datos
	router.HandleFunc("/offers/reset", routes.ResetDatabaseHandler).Methods("POST")

	// Rutas de gestión de ofertas
	router.HandleFunc("/offers", routes.PostOfferHandler).Methods("POST")
	router.HandleFunc("/offers", routes.GetOffersHandler).Methods("GET")
	router.HandleFunc("/offers/{id}", routes.GetOfferHandler).Methods(http.MethodGet)
	router.HandleFunc("/offers/{id}", routes.DeleteOfferHandler).Methods("DELETE")

	return router
}

func main() {
	db.DBconnection()

	if err := db.DB.Migrator().DropTable(&models.Offer{}); err != nil {
		log.Fatalf("Error deleting Offer table: %v", err)
	}

	if err := migrations.DropEnumTypes(db.DB); err != nil {
		log.Fatalf("Error deleting enum types: %v", err)
	}

	if err := migrations.DropUUID4(db.DB); err != nil {
		log.Fatalf("Error deleting UUID: %v", err)
	}

	if err := migrations.CreateEnumTypes(db.DB); err != nil {
		log.Fatalf("Error creating enum types: %v", err)
	}

	if err := migrations.CreateUUID4(db.DB); err != nil {
		log.Fatalf("Error creating UUID: %v", err)
	}

	db.DB.AutoMigrate(models.Offer{})
	router := SetupRouter()

	http.ListenAndServe(":3003", router)
}
