package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"offer_management/config"
	"offer_management/db"
	"offer_management/migrations"
	"offer_management/models"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() func() {
	os.Setenv("ENVIRONMENT", "test")
	host, user, password, dbName, port := config.GetTestDBConfig()

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port + " sslmode=disable"
	var err error

	// Conectar a la base de datos
	db.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Obtener la conexión SQL cruda
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get raw database connection: %v", err)
	}

	// Eliminar todas las tablas en la base de datos
	_, err = sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatalf("Failed to clean database: %v", err)
	}

	// Crear la extensión uuid-ossp
	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		log.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}

	// Crear los tipos enumerados
	if err := migrations.CreateEnumTypes(db.DB); err != nil {
		log.Fatalf("Error creating enum types: %v", err)
	}

	// Migrar el esquema
	if err := db.DB.AutoMigrate(&models.Offer{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Retornar la función de limpieza
	return func() {
		sqlDB, _ := db.DB.DB()
		sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	}
}

func TestPostOfferHandler(t *testing.T) {
	// Setup: prepara una base de datos de prueba
	cleanup := setupTestDB()
	defer cleanup()

	// Datos de prueba
	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
	}

	// Serializa la oferta a JSON
	body, _ := json.Marshal(offer)

	// Crea una solicitud HTTP POST con el body y agrega el token de autenticación al header
	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy-token") // Usa un token de prueba

	// Crea un ResponseRecorder para capturar la respuesta del handler
	w := httptest.NewRecorder()

	// Llama al handler
	PostOfferHandler(w, req)

	// Valida la respuesta
	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var responseOffer models.Offer
	json.NewDecoder(res.Body).Decode(&responseOffer)

	// Valida que la respuesta sea correcta
	assert.Equal(t, offer.ID, responseOffer.ID)
	assert.Equal(t, offer.PostId, responseOffer.PostId)
	assert.Equal(t, offer.Description, responseOffer.Description)
	assert.Equal(t, offer.Size, responseOffer.Size)
	assert.Equal(t, offer.Fragile, responseOffer.Fragile)
	assert.Equal(t, offer.Offer, responseOffer.Offer)
	assert.NotEmpty(t, responseOffer.UserId)
}

func TestPostOfferWithoutTokenHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()
	os.Setenv("ENVIRONMENT", "")

	// Datos de prueba
	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
	}

	body, _ := json.Marshal(offer)

	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))

	w := httptest.NewRecorder()

	PostOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestPostOfferWithBadSize(t *testing.T) {
	// Setup: prepara una base de datos de prueba
	cleanup := setupTestDB()
	defer cleanup()

	// Datos de prueba
	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "NONE",
		Fragile:     true,
		Offer:       10.00,
	}

	// Serializa la oferta a JSON
	body, _ := json.Marshal(offer)

	// Crea una solicitud HTTP POST con el body y agrega el token de autenticación al header
	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy-token") // Usa un token de prueba

	w := httptest.NewRecorder()

	PostOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusPreconditionFailed, res.StatusCode)
}

func TestPostOfferWithMissingrRequeredData(t *testing.T) {
	// Setup: prepara una base de datos de prueba
	cleanup := setupTestDB()
	defer cleanup()

	// Datos de prueba
	offer := &models.Offer{
		ID:      "550e8400-e29b-41d4-a716-446655440000",
		Size:    "NONE",
		Fragile: true,
		Offer:   10.00,
	}

	// Serializa la oferta a JSON
	body, _ := json.Marshal(offer)

	// Crea una solicitud HTTP POST con el body y agrega el token de autenticación al header
	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy-token") // Usa un token de prueba

	w := httptest.NewRecorder()

	PostOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestPostOfferWithNegativeOffer(t *testing.T) {
	// Setup: prepara una base de datos de prueba
	cleanup := setupTestDB()
	defer cleanup()

	// Datos de prueba
	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       -10.00,
	}

	// Serializa la oferta a JSON
	body, _ := json.Marshal(offer)

	// Crea una solicitud HTTP POST con el body y agrega el token de autenticación al header
	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer dummy-token") // Usa un token de prueba

	w := httptest.NewRecorder()

	PostOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusPreconditionFailed, res.StatusCode)
}

func TestGetOffersHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}

	db.DB.Create(offer)

	req := httptest.NewRequest(http.MethodGet, "/offers", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseOffers []models.Offer
	json.NewDecoder(res.Body).Decode(&responseOffers)

	assert.Len(t, responseOffers, 1)
	assert.Equal(t, offer.ID, responseOffers[0].ID)
}

func TestGetOffersWithOutTokenHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()
	os.Setenv("ENVIRONMENT", "")

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}

	db.DB.Create(offer)

	req := httptest.NewRequest(http.MethodGet, "/offers", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestGetOfferWithIcorrectOwner(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
	}

	db.DB.Create(offer)

	req := httptest.NewRequest(http.MethodGet, "/offers?owner=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseOffers []models.Offer
	json.NewDecoder(res.Body).Decode(&responseOffers)

	assert.Len(t, responseOffers, 0)
	assert.Empty(t, responseOffers)
}

func TestGetOfferHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}

	db.DB.Create(offer)

	// Caso exitoso
	req := httptest.NewRequest(http.MethodGet, "/offers/550e8400-e29b-41d4-a716-446655440000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "550e8400-e29b-41d4-a716-446655440000"})
	w := httptest.NewRecorder()

	GetOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseOffer models.Offer
	json.NewDecoder(res.Body).Decode(&responseOffer)

	assert.Equal(t, offer.ID, responseOffer.ID)
	assert.Equal(t, offer.PostId, responseOffer.PostId)
	assert.Equal(t, offer.Description, responseOffer.Description)
	assert.Equal(t, offer.Size, responseOffer.Size)
	assert.Equal(t, offer.Fragile, responseOffer.Fragile)
	assert.Equal(t, offer.Offer, responseOffer.Offer)
	assert.Equal(t, offer.UserId, responseOffer.UserId)

	// Caso de oferta no encontrada
	req = httptest.NewRequest(http.MethodGet, "/offers/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w = httptest.NewRecorder()

	GetOfferHandler(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteOfferHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}

	db.DB.Create(offer)

	// Caso exitoso
	req := httptest.NewRequest(http.MethodDelete, "/offers/550e8400-e29b-41d4-a716-446655440000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "550e8400-e29b-41d4-a716-446655440000"})
	w := httptest.NewRecorder()

	DeleteOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "la oferta fue eliminada", response["msg"])

	var checkOffer models.Offer
	err := db.DB.First(&checkOffer, "id = ?", offer.ID).Error
	assert.Error(t, err)
}

func TestDeleteOfferWithOutTokenHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()
	os.Setenv("ENVIRONMENT", "")

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}

	db.DB.Create(offer)

	// Caso exitoso
	req := httptest.NewRequest(http.MethodDelete, "/offers/550e8400-e29b-41d4-a716-446655440000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "550e8400-e29b-41d4-a716-446655440000"})
	w := httptest.NewRecorder()

	DeleteOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/offers/ping", nil)
	w := httptest.NewRecorder()

	PingHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body := w.Body.String()
	assert.Equal(t, "pong", body)
}

func TestResetDatabaseHandler(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6b",
		PostId:      "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o7c",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o8d",
	}

	db.DB.Create(offer)

	req := httptest.NewRequest(http.MethodPost, "/offers/reset", nil)
	w := httptest.NewRecorder()

	ResetDatabaseHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]string
	json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "Todos los datos fueron eliminados", response["msg"])

	var count int64
	db.DB.Model(&models.Offer{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestPostOfferHandler_InvalidData(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	invalidOffer := &models.Offer{
		PostId:      "", // Invalid data
		Description: "",
		Size:        "",
		Offer:       -1.00,
	}

	body, _ := json.Marshal(invalidOffer)
	req := httptest.NewRequest(http.MethodPost, "/offers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	PostOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetOffersHandler_WithQueryParams(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	offer := &models.Offer{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		PostId:      "550e8400-e29b-41d4-a716-446655440000",
		Description: "Test Offer",
		Size:        "LARGE",
		Fragile:     true,
		Offer:       10.00,
		UserId:      "550e8400-e29b-41d4-a716-446655440000",
	}
	db.DB.Create(offer)

	req := httptest.NewRequest(http.MethodGet, "/offers?post=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	GetOffersHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var responseOffers []models.Offer
	json.NewDecoder(res.Body).Decode(&responseOffers)

	assert.Len(t, responseOffers, 1)
	assert.Equal(t, offer.ID, responseOffers[0].ID)
}

func TestGetOfferHandler_WithBadUid(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/offers/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	GetOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteOfferHandler_StatusBadRequest(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/offers", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	DeleteOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteOfferHandler_NotFound(t *testing.T) {
	cleanup := setupTestDB()
	defer cleanup()

	req := httptest.NewRequest(http.MethodDelete, "/offers", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "550e8400-e29b-41d4-a716-446655440000"})
	w := httptest.NewRecorder()

	DeleteOfferHandler(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
