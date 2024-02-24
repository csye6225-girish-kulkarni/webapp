package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"webapp/controller"
	"webapp/db"
	"webapp/middleware"
	"webapp/service"
	"webapp/types"
	"webapp/utils"
)

func setupDB() *db.PostgreSQL {
	connString := os.Getenv("POSTGRES_CONN_STR_TEST")
	connString = connString + "?sslmode=disable"
	postgresObj := db.NewPostgreSQL(connString)
	fmt.Println(connString)
	err := postgresObj.Ping(nil)

	if err != nil {
		log.Fatalf("Unable to ping to DB err : %v", err)
	}

	fmt.Println("Successfully Connected to DB")
	return postgresObj
}

func teardownDB(postgresObj *db.PostgreSQL) {
	_ = postgresObj.DB.Exec("DELETE FROM USERS")

	err := postgresObj.Close()
	if err != nil {
		log.Fatalf("Unable to close the DB connection err : %v", err)
	}
	fmt.Println("Successfully Closed the DB Connection")
}

func setupTestRouter(postgresObj *db.PostgreSQL) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	router.Use(middleware.SetNoCacheHeader())

	healthService := service.NewHealthService(postgresObj)
	userService := service.NewUserService(postgresObj)
	userController := controller.NewUserController(userService)
	healthController := controller.NewHealthController(healthService)

	router.GET("/healthz", middleware.CheckNoAuthEndpoints(), healthController.GetHealth)
	router.Use(func(context *gin.Context) {
		if context.Request.URL.Path == "/healthz" && context.Request.Method != http.MethodGet {
			context.Status(http.StatusMethodNotAllowed)
			context.Abort()
		}
	})

	router.GET("/v1/user", middleware.BasicAuth(userService), userController.GetUser)
	router.PUT("/v1/user/self", middleware.BasicAuth(userService), userController.UpdateUser)

	router.POST("/v1/user", middleware.CheckNoAuthEndpoints(), userController.CreateUser)
	router.NoRoute(func(context *gin.Context) {
		context.Data(http.StatusNotFound, "text/plain", []byte{})
		context.Abort()
	})

	return router
}

func TestCreateAccount(t *testing.T) {
	postgresObj := setupDB()
	defer teardownDB(postgresObj)
	router := setupTestRouter(postgresObj)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Creating a new user account with the following details
	createAccReq := types.UserRequest{
		Username:  "test@user.com",
		Password:  "testpassword",
		FirstName: "testfirst",
		LastName:  "testlast",
	}
	r, _ := json.Marshal(createAccReq)

	req, err := http.NewRequest("POST", ts.URL+"/v1/user", strings.NewReader(string(r)))
	if err != nil {
		t.Fatalf("Error while creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error while making request: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, but got: %v", res.StatusCode)
	}
	time.Sleep(1 * time.Second)

	getUserReq, err := http.NewRequest("GET", ts.URL+"/v1/user", nil)
	if err != nil {
		t.Fatalf("Error while creating request: %v", err)
	}

	authHeader := utils.CreateBasicAuth("test@user.com", "testpassword")
	getUserReq.Header.Set("Authorization", authHeader)

	res, err = http.DefaultClient.Do(getUserReq)
	if err != nil {
		t.Fatalf("Error while making request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, but got: %v", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error while reading response body: %v", err)
	}
	var user types.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		t.Fatalf("Error while unmarshalling response body: %v", err)
	}

	// Verifying the provided user details
	if user.Username != "test@user.com" || user.FirstName != "testfirst" || user.LastName != "testlast" {
		t.Fatalf("Expected user details to be {Username:test@user.com, FirstName:testfirst, LastName:testlast}, but got: %v", user)
	}
}

func TestUpdateAccount(t *testing.T) {
	postgresObj := setupDB()
	defer teardownDB(postgresObj)
	router := setupTestRouter(postgresObj)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Creating a new user account
	createAccReq := types.UserRequest{
		Username:  "test@user.com",
		Password:  "testpassword",
		FirstName: "testfirst",
		LastName:  "testlast",
	}
	r, _ := json.Marshal(createAccReq)

	req, err := http.NewRequest("POST", ts.URL+"/v1/user", strings.NewReader(string(r)))
	if err != nil {
		t.Fatalf("Error while creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error while making request: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code 201, but got: %v", res.StatusCode)
	}
	time.Sleep(1 * time.Second)

	// Updating the user account
	updateAccReq := types.UpdateUserRequest{
		Password:  "testpassword",
		FirstName: "updatedfirst",
		LastName:  "updatedlast",
	}
	r, _ = json.Marshal(updateAccReq)

	req, err = http.NewRequest("PUT", ts.URL+"/v1/user/self", strings.NewReader(string(r)))
	if err != nil {
		t.Fatalf("Error while creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	authHeader := utils.CreateBasicAuth("test@user.com", "testpassword")
	req.Header.Set("Authorization", authHeader)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error while making request: %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code 200, but got: %v", res.StatusCode)
	}
	time.Sleep(1 * time.Second)

	// Getting the updated user account
	getUserReq, err := http.NewRequest("GET", ts.URL+"/v1/user", nil)
	if err != nil {
		t.Fatalf("Error while creating request: %v", err)
	}
	getUserReq.Header.Set("Authorization", authHeader)

	res, err = http.DefaultClient.Do(getUserReq)
	if err != nil {
		t.Fatalf("Error while making request: %v", err)
	}
	if res.StatusCode != http.StatusBadGateway {
		t.Fatalf("Expected status code 200, but got: %v", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error while reading response body: %v", err)
	}
	var user types.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		t.Fatalf("Error while unmarshalling response body: %v", err)
	}

	// Verifying the updated user details
	if user.Username != "test@user.com" || user.FirstName != "updatedfirst" || user.LastName != "updatedlast" {
		t.Fatalf("Expected user details to be {Username:test@user.com, FirstName:updatedfirst, LastName:updatedlast}, but got: %v", user)
	}
}
