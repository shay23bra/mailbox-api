package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mailbox-api/api/router"
	"mailbox-api/config"
	"mailbox-api/logger"
	"mailbox-api/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *config.Config) {
	// Load config
	cfg, _ := config.Load()

	// For testing, override JWT secret
	cfg.Auth.JWTSecret = "test-secret-key"

	// Create logger
	log := logger.NewLogger()

	// Create service
	mailboxService := service.NewMailboxService(testMailboxRepo, testDepartmentRepo)

	// Create router
	r := router.SetupRouter(cfg, log, mailboxService)

	return r.GetEngine(), cfg
}

func TestGetToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _ := setupTestRouter()

	// Test CEO token endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/token/ceo", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])

	// Test CTO token endpoint
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/token/cto", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}

// These are currently disabled, should populate data before testing these.
// func TestGetMailboxes(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	router, cfg := setupTestRouter()

// 	// Generate CEO and CTO tokens
// 	ceoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCEO)
// 	ctoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCTO)

// 	// Test CEO access - should see all mailboxes
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/api/mailboxes", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ceoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var ceoResponse model.MailboxResponse
// 	err := json.Unmarshal(w.Body.Bytes(), &ceoResponse)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, ceoResponse.Data)
// 	assert.NotNil(t, ceoResponse.Pagination)

// 	// CEO should see all mailboxes
// 	assert.Equal(t, 5, ceoResponse.Pagination.TotalItems)

// 	// Test CTO access - should see only their sub-organization
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/api/mailboxes", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var ctoResponse model.MailboxResponse
// 	err = json.Unmarshal(w.Body.Bytes(), &ctoResponse)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, ctoResponse.Data)
// 	assert.NotNil(t, ctoResponse.Pagination)

// 	// CTO should see only their direct and indirect reports
// 	// (This number might need adjustment based on your test data)
// 	assert.Less(t, ctoResponse.Pagination.TotalItems, ceoResponse.Pagination.TotalItems)
// }

// func TestGetMailbox(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	router, cfg := setupTestRouter()

// 	// Generate CEO and CTO tokens
// 	ceoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCEO)
// 	ctoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCTO)

// 	// Test CEO access to any mailbox
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/api/mailboxes/cto@example.com", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ceoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Test CTO access to their own mailbox
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/api/mailboxes/cto@example.com", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Test CTO access to a mailbox in their sub-organization
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/api/mailboxes/dev1@example.com", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Test CTO access to a mailbox outside their sub-organization
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/api/mailboxes/marketing@example.com", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusForbidden, w.Code)
// }

// func TestCalculateOrgMetrics(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	router, cfg := setupTestRouter()

// 	// Generate CEO and CTO tokens
// 	ceoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCEO)
// 	ctoToken, _ := middleware.GenerateToken(cfg, middleware.RoleCTO)

// 	// Test CEO access to calculate metrics
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/api/mailboxes/calculate-metrics", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ceoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Test CTO access to calculate metrics (should be forbidden)
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("POST", "/api/mailboxes/calculate-metrics", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ctoToken))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusForbidden, w.Code)
// }
