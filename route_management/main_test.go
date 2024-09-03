package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"route_management/mocks"
	"route_management/model"
	"route_management/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoutesEnabled(t *testing.T) {
	assert.NotPanics(t, func() {
		Start()
	})
}

func TestPingRoute(t *testing.T) {
	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/routes/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPing(t *testing.T) {
	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/routes/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPostRouteHAPPY(t *testing.T) {
	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now().Add(time.Hour),
		PlannedEndDate:     time.Now().Add(time.Hour).Add(time.Hour),
	}

	mock.ExpectQuery("SELECT * FROM \"routes\" .*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mock.NewRows([]string{}))
	mock.ExpectCommit()
	mock.ExpectQuery("INSERT INTO routes").WithArgs(
		mockTrayecto.ID,
		mockTrayecto.FlightId,
		mockTrayecto.SourceAirportCode,
		mockTrayecto.DestinyAirportCode,
		mockTrayecto.SourceCountry,
		mockTrayecto.DestinyCountry,
		mockTrayecto.BagCost,
		mockTrayecto.PlannedStartDate,
		mockTrayecto.PlannedEndDate,
	).
		WillReturnRows(mock.NewRows([]string{mockTrayecto.ID, "1"}))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(mockTrayecto)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/routes", &buf)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestPostRouteNoToken(t *testing.T) {
	authenticator := mocks.GetInvalidAuthenticator(403, "what")
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}

	mock.ExpectQuery("SELECT * FROM \"routes\"*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mock.NewRows([]string{}))
	mock.ExpectCommit()
	mock.ExpectQuery("INSERT INTO routes").WithArgs(
		mockTrayecto.ID,
		mockTrayecto.FlightId,
		mockTrayecto.SourceAirportCode,
		mockTrayecto.DestinyAirportCode,
		mockTrayecto.SourceCountry,
		mockTrayecto.DestinyCountry,
		mockTrayecto.BagCost,
		mockTrayecto.PlannedStartDate,
		mockTrayecto.PlannedEndDate,
	).
		WillReturnRows(mock.NewRows([]string{mockTrayecto.ID, "1"}))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(mockTrayecto)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/routes", &buf)

	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

func TestPostRouteInvalidBody(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}

	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}

	mock.ExpectQuery("SELECT * FROM \"routes\"*").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mock.NewRows([]string{}))
	mock.ExpectCommit()
	mock.ExpectQuery("INSERT INTO routes").WithArgs(
		mockTrayecto.ID,
		mockTrayecto.FlightId,
		mockTrayecto.SourceAirportCode,
		mockTrayecto.DestinyAirportCode,
		mockTrayecto.SourceCountry,
		mockTrayecto.DestinyCountry,
		mockTrayecto.BagCost,
		mockTrayecto.PlannedStartDate,
		mockTrayecto.PlannedEndDate,
	).
		WillReturnRows(mock.NewRows([]string{mockTrayecto.ID, "1"}))
	mock.ExpectCommit()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(mockTrayecto)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/routes", &buf)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetRoutesHAPPY(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}
	query := `SELECT \* FROM "routes".*`

	mock.ExpectQuery(query).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			mock.NewRows([]string{
				"id",
				"flight_id",
				"source_airport_code",
				"destiny_airport_code",
				"source_country",
				"destiny_country",
				"bag_cost",
				"planned_start_date",
				"planned_end_date",
			}).
				AddRow(
					mockTrayecto.ID,
					mockTrayecto.FlightId,
					mockTrayecto.SourceAirportCode,
					mockTrayecto.DestinyAirportCode,
					mockTrayecto.SourceCountry,
					mockTrayecto.DestinyCountry,
					mockTrayecto.BagCost,
					mockTrayecto.PlannedStartDate,
					mockTrayecto.PlannedEndDate,
				))
	mock.ExpectCommit()

	req, _ := http.NewRequest("GET", "/routes?flight=1", nil)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	var response []model.Route
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, mockTrayecto.FlightId, response[0].FlightId)
}
func TestGetRoutesNoFlightIdHAPPY(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}
	query := `SELECT \* FROM "routes".*`

	mock.ExpectQuery(query).
		WithoutArgs().
		WillReturnRows(
			mock.NewRows([]string{
				"id",
				"flight_id",
				"source_airport_code",
				"destiny_airport_code",
				"source_country",
				"destiny_country",
				"bag_cost",
				"planned_start_date",
				"planned_end_date",
			}).
				AddRow(
					mockTrayecto.ID,
					mockTrayecto.FlightId,
					mockTrayecto.SourceAirportCode,
					mockTrayecto.DestinyAirportCode,
					mockTrayecto.SourceCountry,
					mockTrayecto.DestinyCountry,
					mockTrayecto.BagCost,
					mockTrayecto.PlannedStartDate,
					mockTrayecto.PlannedEndDate,
				))
	mock.ExpectCommit()

	req, _ := http.NewRequest("GET", "/routes", nil)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	var response []model.Route
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, mockTrayecto.FlightId, response[0].FlightId)
}

func TestGetRouteHAPPY(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}
	query := `SELECT \* FROM "routes".*`

	mock.ExpectQuery(query).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			mock.NewRows([]string{
				"id",
				"flight_id",
				"source_airport_code",
				"destiny_airport_code",
				"source_country",
				"destiny_country",
				"bag_cost",
				"planned_start_date",
				"planned_end_date",
			}).
				AddRow(
					mockTrayecto.ID,
					mockTrayecto.FlightId,
					mockTrayecto.SourceAirportCode,
					mockTrayecto.DestinyAirportCode,
					mockTrayecto.SourceCountry,
					mockTrayecto.DestinyCountry,
					mockTrayecto.BagCost,
					mockTrayecto.PlannedStartDate,
					mockTrayecto.PlannedEndDate,
				))
	mock.ExpectCommit()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/routes/%s", uuid.NewString()), nil)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response model.Route
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, mockTrayecto.FlightId, response.FlightId)
}

func TestGetRoutenonUIID(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()

	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}
	query := `SELECT \* FROM "routes".*`

	mock.ExpectQuery(query).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			mock.NewRows([]string{
				"id",
				"flight_id",
				"source_airport_code",
				"destiny_airport_code",
				"source_country",
				"destiny_country",
				"bag_cost",
				"planned_start_date",
				"planned_end_date",
			}).
				AddRow(
					mockTrayecto.ID,
					mockTrayecto.FlightId,
					mockTrayecto.SourceAirportCode,
					mockTrayecto.DestinyAirportCode,
					mockTrayecto.SourceCountry,
					mockTrayecto.DestinyCountry,
					mockTrayecto.BagCost,
					mockTrayecto.PlannedStartDate,
					mockTrayecto.PlannedEndDate,
				))
	mock.ExpectCommit()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/routes/%s", "uuid.NewString()"), nil)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestDeleteRouteHAPPY(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	mockResultId := uuid.NewString()
	mockTrayecto := model.Route{
		ID:                 mockResultId,
		FlightId:           "1",
		SourceAirportCode:  "test-source-airport",
		DestinyAirportCode: "test-destiny-airport",
		SourceCountry:      "test-source-country",
		DestinyCountry:     "test-destiny-country",
		BagCost:            12,
		PlannedStartDate:   time.Now(),
		PlannedEndDate:     time.Now().Add(time.Hour),
	}

	query := `SELECT \* FROM "routes".*`

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(mock.NewRows([]string{
			"id",
			"flight_id",
			"source_airport_code",
			"destiny_airport_code",
			"source_country",
			"destiny_country",
			"bag_cost",
			"planned_start_date",
			"planned_end_date",
		}).
			AddRow(
				mockTrayecto.ID,
				mockTrayecto.FlightId,
				mockTrayecto.SourceAirportCode,
				mockTrayecto.DestinyAirportCode,
				mockTrayecto.SourceCountry,
				mockTrayecto.DestinyCountry,
				mockTrayecto.BagCost,
				mockTrayecto.PlannedStartDate,
				mockTrayecto.PlannedEndDate,
			))
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/routes/%s", uuid.NewString()), nil)
	req.Header.Add("Authorization", uuid.NewString())

	router.ServeHTTP(w, req)

	type Response struct {
		Msg string
	}
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "el trayecto fue eliminado", response.Msg)
}

func TestResetRoutesHAPPY(t *testing.T) {

	authenticator := mocks.GetValidAuthenticator()
	mockConfig := utils.Config{
		User_url:      "test",
		Authenticator: authenticator,
	}
	mockDb, mock, _ := sqlmock.New()
	defer mockDb.Close()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	router := setupRouter(db, mockConfig)

	w := httptest.NewRecorder()

	query := `DELETE \*`

	mock.ExpectQuery(query).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(
			mock.NewRows([]string{}),
		)
	mock.ExpectCommit()

	req, _ := http.NewRequest("POST", "/routes/reset", nil)
	req.Header.Add("Authorization", "1234")

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	type Response struct {
		Msg string
	}
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "todos los datos fueron eliminados", response.Msg)
}
