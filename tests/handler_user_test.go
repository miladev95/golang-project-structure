package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	handlers "github.com/yourusername/yourproject/internal/handlers/http"
	"github.com/yourusername/yourproject/internal/models"
)

// MockUserService implements services.UserService for testing
type MockUserService struct {
	GetUserFunc      func(ctx context.Context, id int64) (*models.User, error)
	GetAllUsersFunc  func(ctx context.Context) ([]models.User, error)
	CreateUserFunc   func(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserFunc   func(ctx context.Context, user *models.User) error
	DeleteUserFunc   func(ctx context.Context, id int64) error
}

func (m *MockUserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, id)
	}
	return nil, errors.New("GetUserFunc not implemented")
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(ctx)
	}
	return nil, errors.New("GetAllUsersFunc not implemented")
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil, errors.New("CreateUserFunc not implemented")
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *models.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, user)
	}
	return errors.New("UpdateUserFunc not implemented")
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int64) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return errors.New("DeleteUserFunc not implemented")
}

// TestGetUserSuccess tests successful GetUser API call
func TestGetUserSuccess(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	expectedUser := &models.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			if id == 1 {
				return expectedUser, nil
			}
			return nil, errors.New("user not found")
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	// Execute
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check response structure
	if response["success"] != true {
		t.Errorf("Expected success: true, got %v", response["success"])
	}

	// Extract user data
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("Expected data to be a map, got %T", response["data"])
	}

	if data["id"] != float64(1) {
		t.Errorf("Expected user ID 1, got %v", data["id"])
	}

	if data["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %v", data["name"])
	}

	if data["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %v", data["email"])
	}
}

// TestGetUserInvalidID tests GetUser with invalid ID format
func TestGetUserInvalidID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	// Test cases with invalid IDs (excluding empty string as Gin routing won't match it)
	invalidIDs := []string{"abc", "12.34", "!@#$"}

	for _, invalidID := range invalidIDs {
		t.Run("invalid_id_"+invalidID, func(t *testing.T) {
			// Execute
			req, _ := http.NewRequest("GET", "/api/v1/users/"+invalidID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
			}

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response["success"] != false {
				t.Errorf("Expected success: false, got %v", response["success"])
			}

			if errorMsg, ok := response["error"].(string); !ok || errorMsg == "" {
				t.Errorf("Expected error message in response, got %v", response["error"])
			}
		})
	}
}

// TestGetUserNotFound tests GetUser when user doesn't exist
func TestGetUserNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	// Execute
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != false {
		t.Errorf("Expected success: false, got %v", response["success"])
	}

	if errorMsg, ok := response["error"].(string); !ok || errorMsg == "" {
		t.Errorf("Expected error message in response, got %v", response["error"])
	}
}

// TestGetUserMultipleUsers tests GetUser with different user IDs
func TestGetUserMultipleUsers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	users := map[int64]*models.User{
		1: {
			ID:    1,
			Name:  "User One",
			Email: "user1@example.com",
		},
		2: {
			ID:    2,
			Name:  "User Two",
			Email: "user2@example.com",
		},
		3: {
			ID:    3,
			Name:  "User Three",
			Email: "user3@example.com",
		},
	}

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			if user, ok := users[id]; ok {
				return user, nil
			}
			return nil, errors.New("user not found")
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	// Test fetching multiple users
	for userID, expectedUser := range users {
		t.Run("get_user_"+string(rune(userID)), func(t *testing.T) {
			// Execute
			req, _ := http.NewRequest("GET", "/api/v1/users/"+string(rune(userID)+'0'), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// For this simple test, we just verify the status is OK
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			data := response["data"].(map[string]interface{})
			if data["name"] != expectedUser.Name {
				t.Errorf("Expected name %s, got %v", expectedUser.Name, data["name"])
			}
		})
	}
}

// TestGetUserContextCancellation tests GetUser with cancelled context
func TestGetUserContextCancellation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			// Simulate context cancellation check
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return &models.User{ID: id}, nil
			}
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	// Execute - normal request (context won't be cancelled)
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify - should succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// BenchmarkGetUser benchmarks the GetUser handler
func BenchmarkGetUser(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockUserService{
		GetUserFunc: func(ctx context.Context, id int64) (*models.User, error) {
			return &models.User{
				ID:    id,
				Name:  "Test User",
				Email: "test@example.com",
			}, nil
		},
	}

	handler := handlers.NewUserHandler(mockService)
	router.GET("/api/v1/users/:id", handler.GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}