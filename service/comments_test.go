package service_test

import (
	"errors"
	"newsfeed-service/domains"
	"newsfeed-service/models"
	"newsfeed-service/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentDB is a mock for the CommentDB interface
type MockCommentDB struct {
	mock.Mock
}

// Create mocks the Create method of CommentDB
func (m *MockCommentDB) Create(comment *models.Comment) (*models.Comment, error) {
	args := m.Called(comment)
	respVal := args.Get(0)
	if respVal != nil {
		return respVal.(*models.Comment), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete mocks the Delete method of CommentDB
func (m *MockCommentDB) Delete(userID string, commentID uint) error {
	args := m.Called(userID, commentID)
	return args.Error(1)
}

func TestCreateCommentSuccess(t *testing.T) {
	// Setup mock
	mockDB := new(MockCommentDB)
	expectedComment := &models.Comment{
		Creator: "user123",
		PostID:  10,
		Content: "This is a comment",
	}
	mockDB.On("Create", expectedComment).Return(expectedComment, nil)

	// Create service with mock
	commentService := service.NewCommentService(mockDB)

	// Create comment request
	req := domains.CreateCommentRequest{
		Content: "This is a comment",
	}

	// Call CreateComment
	commentID, err := commentService.CreateComment("user123", 10, req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedComment.ID, commentID)
	mockDB.AssertExpectations(t)
}

func TestCreateCommentError(t *testing.T) {
	// Setup mock
	mockDB := new(MockCommentDB)
	expectedError := errors.New("database error")
	mockDB.On("Create", mock.Anything).Return(nil, expectedError)

	// Create service with mock
	commentService := service.NewCommentService(mockDB)

	// Create comment request
	req := domains.CreateCommentRequest{
		Content: "This is a comment",
	}

	// Call CreateComment
	commentID, err := commentService.CreateComment("user123", 10, req)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, uint(0), commentID)
	mockDB.AssertExpectations(t)
}
