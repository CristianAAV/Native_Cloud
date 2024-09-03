package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"offer_management/auth"
	"offer_management/config"
	"offer_management/db"
	"offer_management/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// PostOfferHandler maneja la creación de nuevas ofertas
func PostOfferHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entrando a PostOfferHandler")

	var offer models.Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Println("Solicitud decodificada correctamente")

	// Validar campos obligatorios
	if offer.PostId == "" || offer.Description == "" || offer.Size == "" {
		log.Println("Validación de campos fallida: campos obligatorios faltantes o inválidos")
		http.Error(w, "Precondition Failed", http.StatusBadRequest)
		return
	}

	if offer.Offer < 0 {
		log.Println("Oferta negativa, la oferta debe ser positiva")
		http.Error(w, "Precondition Failed", http.StatusPreconditionFailed)
		return
	}

	// Validar el campo `Size`
	validSizes := map[string]bool{"LARGE": true, "MEDIUM": true, "SMALL": true}
	if !validSizes[offer.Size] {
		log.Println("Validación de tamaño fallida: tamaño del paquete no es válido")
		http.Error(w, "Precondition Failed", http.StatusPreconditionFailed)
		return
	}
	log.Println("Campos obligatorios y tamaño del paquete validados correctamente")

	// Verificar el token de autorización y obtener el userId
	userId, err := auth.ProcessAuthorization(r)
	if err != nil {
		log.Printf("Error en la autorización: %v\n", err)
		if err.Error() == "Forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden) // 403 Forbidden
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // 401 Unauthorized
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError) // 500 Internal Server Error para otros errores
		}
		return
	}
	log.Printf("userId obtenido exitosamente: %s\n", userId)

	// Asignar el userId al offer
	offer.UserId = userId
	log.Printf("Asignado userId '%s' a la oferta\n", userId)

	// Guardar la oferta en la base de datos
	if err := db.DB.Create(&offer).Error; err != nil {
		log.Printf("Error al guardar la oferta en la base de datos: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Oferta guardada exitosamente en la base de datos")

	// Preparar y enviar la respuesta con la oferta creada
	responseOffer := config.OfferResponse{
		Id:          offer.ID,
		PostId:      offer.PostId,
		UserId:      offer.UserId,
		Description: offer.Description,
		Size:        offer.Size,
		Fragile:     offer.Fragile,
		Offer:       offer.Offer,
		CreatedAt:   offer.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(responseOffer); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v\n", err)
	}
	log.Println("Respuesta enviada exitosamente")
}

// GetOffersHandler maneja la recuperación de todas las ofertas
func GetOffersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entrando a GetOffersHandler")

	// Verificar el token de autorización y obtener el userId
	userId, err := auth.ProcessAuthorization(r)
	if err != nil {
		log.Printf("Error en la autorización: %v\n", err)
		if err.Error() == "Forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden) // 403 Forbidden
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // 401 Unauthorized
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError) // 500 Internal Server Error para otros errores
		}
		return
	}
	log.Printf("userId obtenido exitosamente: %s\n", userId)

	var offers []models.Offer
	query := db.DB

	if postID := r.URL.Query().Get("post"); postID != "" {
		query = query.Where("post_id = ?", postID)
	}

	if owner := r.URL.Query().Get("owner"); owner != "" {
		if owner == "me" {
			query = query.Where("user_id = ?", userId)
		} else {
			query = query.Where("user_id = ?", owner)
		}
	}

	if err := query.Find(&offers).Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseOffers := []config.OfferResponse{}

	for _, offer := range offers {
		responseOffers = append(responseOffers, config.OfferResponse{
			Id:          offer.ID,
			PostId:      offer.PostId,
			UserId:      offer.UserId,
			Description: offer.Description,
			Size:        offer.Size,
			Fragile:     offer.Fragile,
			Offer:       offer.Offer,
		})
	}

	json.NewEncoder(w).Encode(responseOffers)
}

// GetOfferHandler maneja la recuperación de una oferta específica por ID
func GetOfferHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entrando a GetOfferHandler")

	// Verificar el token de autorización y obtener el userId
	_, err := auth.ProcessAuthorization(r)
	if err != nil {
		log.Printf("Error en la autorización: %v\n", err)
		if err.Error() == "Forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden) // 403 Forbidden
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // 401 Unauthorized
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError) // 500 Internal Server Error para otros errores
		}
		return
	}
	log.Println("userId obtenido exitosamente")

	params := mux.Vars(r)
	id := params["id"]

	// Validar que el ID sea un UUID válido
	if _, err := uuid.Parse(id); err != nil {
		log.Printf("ID inválido: %v\n", err)
		http.Error(w, "Bad Request - Invalid ID", http.StatusBadRequest)
		return
	}

	var offer models.Offer
	if err := db.DB.First(&offer, "id = ?", id).Error; err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseOffer := config.OfferResponse{
		Id:          offer.ID,
		PostId:      offer.PostId,
		UserId:      offer.UserId,
		Description: offer.Description,
		Size:        offer.Size,
		Fragile:     offer.Fragile,
		Offer:       offer.Offer,
	}

	json.NewEncoder(w).Encode(responseOffer)
}

// DeleteOfferHandler maneja la eliminación de una oferta específica por ID
func DeleteOfferHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entrando a DeleteOfferHandler")

	// Verificar el token de autorización y obtener el userId
	_, err := auth.ProcessAuthorization(r)
	if err != nil {
		log.Printf("Error en la autorización: %v\n", err)
		if err.Error() == "Forbidden" {
			http.Error(w, "Forbidden", http.StatusForbidden) // 403 Forbidden
		} else if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // 401 Unauthorized
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError) // 500 Internal Server Error para otros errores
		}
		return
	}
	log.Println("userId obtenido exitosamente")

	id := mux.Vars(r)["id"]

	// Buscar la oferta por ID
	var offer models.Offer
	if err := db.DB.First(&offer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		log.Printf("Error al buscar la oferta: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	// Eliminar la oferta
	if err := db.DB.Unscoped().Delete(&offer).Error; err != nil {
		log.Printf("Error al eliminar la oferta: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "la oferta fue eliminada"})
}

// PingHandler devuelve un simple "pong" para verificar el estado del servidor
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// ResetDatabaseHandler elimina todas las ofertas de la base de datos
func ResetDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.DB.Exec("DELETE FROM offers").Error; err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "Todos los datos fueron eliminados"})
}
