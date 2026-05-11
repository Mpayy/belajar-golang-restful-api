package test

import (
	"belajar-golang-restful-api/app"
	"belajar-golang-restful-api/controller"
	"belajar-golang-restful-api/helper"
	"belajar-golang-restful-api/middleware"
	"belajar-golang-restful-api/model/domain"
	"belajar-golang-restful-api/repository"
	"belajar-golang-restful-api/service"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func setupTestDb() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/belajar-golang-restful-api-test")
	helper.PanicIfError(err)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":"Gadget"}`)
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/category", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}
func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name":""}`)
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/category", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	requestBody := strings.NewReader(`{"name":"Gadget"}`)
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}
func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	requestBody := strings.NewReader(`{"name":""}`)
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget", responseBody["data"].(map[string]interface{})["name"])
}
func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/category/404", nil)
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	request, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/category/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}
func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/category/404", nil)
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestListCategorySuccess(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category1 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	category2 := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Computer",
	})
	tx.Commit()

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/category", nil)
	request.Header.Add("X-API-KEY", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	var categories = responseBody["data"].([]interface{})
	categoriesResponse1 := categories[0].(map[string]interface{})
	categoriesResponse2 := categories[1].(map[string]interface{})
	assert.Equal(t, category1.Id, int(categoriesResponse1["id"].(float64)))
	assert.Equal(t, "Gadget", categoriesResponse1["name"])
	assert.Equal(t, category2.Id, int(categoriesResponse2["id"].(float64)))
	assert.Equal(t, "Gadget", categoriesResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDb()
	truncateCategory(db)
	router := setupRouter(db)

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/category", nil)
	request.Header.Add("X-API-KEY", "SALAH")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	_ = json.Unmarshal(body, &responseBody)
	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZE", responseBody["status"])
}
