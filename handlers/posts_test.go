package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"newsfeed-service/domains"
	"newsfeed-service/handlers"
	"newsfeed-service/middlewares"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(userID string, req domains.PostCreateRequest) (uint, error) {
	args := m.Called(userID, req)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockPostService) GetPostsWithComments(userID string, cursor *uint, pageSize int) (*domains.PostsPagedResult, error) {
	args := m.Called(userID, cursor, pageSize)
	return args.Get(0).(*domains.PostsPagedResult), args.Error(1)
}

func TestCreatePostSuccess(t *testing.T) {
	// Setup
	userID := "a user"
	mockService := new(MockPostService)
	mockService.On("CreatePost", userID, domains.PostCreateRequest{Caption: "test caption"}).Return(uint(10), nil)
	h := handlers.PostsHandler{PostService: mockService}
	router := gin.Default()
	router.Use(middlewares.BearerTokenAuth())
	router.POST("/posts", h.CreatePost)

	// Request body
	reqBody, _ := json.Marshal(domains.PostCreateRequest{Caption: "test caption"})

	// Test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userID)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		ID uint `json:"ID"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(10), response.ID)
	mockService.AssertExpectations(t)
}

func TestCreatePostBadRequest(t *testing.T) {
	// Setup
	mockService := new(MockPostService)
	h := handlers.PostsHandler{PostService: mockService}
	router := gin.Default()
	router.Use(middlewares.BearerTokenAuth())
	router.POST("/posts", h.CreatePost)

	// Invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	userID := "a user"
	req.Header.Set("Authorization", "Bearer "+userID)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CreatePost")
}
