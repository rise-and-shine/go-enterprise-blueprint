package domain

// package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	"github.com/your-org/project/internal/domain"
// 	"github.com/your-org/project/internal/domain/mocks"
// )

// // =============================================================================
// // Test Setup
// // =============================================================================

// func setupUserService(t *testing.T) (*UserService, *mocks.MockUserRepository) {
// 	// NewMockUserRepository automatically registers cleanup with t.Cleanup()
// 	mockRepo := mocks.NewMockUserRepository(t)
// 	service := NewUserService(mockRepo)
// 	return service, mockRepo
// }

// // =============================================================================
// // GetUser Tests
// // =============================================================================

// func TestUserService_GetUser_Success(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	expectedUser := &domain.User{
// 		ID:    "user-123",
// 		Email: "john@example.com",
// 		Name:  "John Doe",
// 	}

// 	// EXPECT style (type-safe, recommended)
// 	mockRepo.EXPECT().
// 		GetByID(ctx, "user-123").
// 		Return(expectedUser, nil).
// 		Once()

// 	// Act
// 	user, err := service.GetUser(ctx, "user-123")

// 	// Assert
// 	require.NoError(t, err)
// 	assert.Equal(t, expectedUser.ID, user.ID)
// 	assert.Equal(t, expectedUser.Email, user.Email)
// }

// func TestUserService_GetUser_NotFound(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	mockRepo.EXPECT().
// 		GetByID(ctx, "nonexistent").
// 		Return(nil, nil). // repo returns nil, nil when not found
// 		Once()

// 	user, err := service.GetUser(ctx, "nonexistent")

// 	assert.Nil(t, user)
// 	assert.ErrorIs(t, err, ErrUserNotFound)
// }

// func TestUserService_GetUser_RepositoryError(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()
// 	dbError := errors.New("database connection failed")

// 	mockRepo.EXPECT().
// 		GetByID(ctx, "user-123").
// 		Return(nil, dbError).
// 		Once()

// 	user, err := service.GetUser(ctx, "user-123")

// 	assert.Nil(t, user)
// 	assert.ErrorContains(t, err, "database connection failed")
// }

// func TestUserService_GetUser_EmptyID(t *testing.T) {
// 	service, _ := setupUserService(t)
// 	ctx := context.Background()

// 	// No mock expectations - repo should not be called
// 	user, err := service.GetUser(ctx, "")

// 	assert.Nil(t, user)
// 	assert.ErrorContains(t, err, "id cannot be empty")
// }

// // =============================================================================
// // CreateUser Tests
// // =============================================================================

// func TestUserService_CreateUser_Success(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	// First call: check if email exists
// 	mockRepo.EXPECT().
// 		GetByEmail(ctx, "new@example.com").
// 		Return(nil, nil). // user doesn't exist
// 		Once()

// 	// Second call: save the new user
// 	// Using mock.MatchedBy for flexible argument matching
// 	mockRepo.EXPECT().
// 		Save(ctx, mock.MatchedBy(func(u *domain.User) bool {
// 			return u.Email == "new@example.com" && u.Name == "New User"
// 		})).
// 		Return(nil).
// 		Once()

// 	user, err := service.CreateUser(ctx, "new@example.com", "New User")

// 	require.NoError(t, err)
// 	assert.Equal(t, "new@example.com", user.Email)
// 	assert.Equal(t, "New User", user.Name)
// }

// func TestUserService_CreateUser_AlreadyExists(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	existingUser := &domain.User{
// 		ID:    "existing-123",
// 		Email: "existing@example.com",
// 	}

// 	mockRepo.EXPECT().
// 		GetByEmail(ctx, "existing@example.com").
// 		Return(existingUser, nil).
// 		Once()

// 	// Save should NOT be called
// 	// (mockery will fail if unexpected calls happen)

// 	user, err := service.CreateUser(ctx, "existing@example.com", "Name")

// 	assert.Nil(t, user)
// 	assert.ErrorContains(t, err, "already exists")
// }

// // =============================================================================
// // Table-Driven Tests Example
// // =============================================================================

// func TestUserService_GetUser_TableDriven(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		userID      string
// 		setupMock   func(*mocks.MockUserRepository)
// 		wantUser    *domain.User
// 		wantErr     error
// 		errContains string
// 	}{
// 		{
// 			name:   "success",
// 			userID: "user-1",
// 			setupMock: func(m *mocks.MockUserRepository) {
// 				m.EXPECT().
// 					GetByID(mock.Anything, "user-1").
// 					Return(&domain.User{ID: "user-1", Email: "test@example.com"}, nil)
// 			},
// 			wantUser: &domain.User{ID: "user-1", Email: "test@example.com"},
// 			wantErr:  nil,
// 		},
// 		{
// 			name:   "not found",
// 			userID: "nonexistent",
// 			setupMock: func(m *mocks.MockUserRepository) {
// 				m.EXPECT().
// 					GetByID(mock.Anything, "nonexistent").
// 					Return(nil, nil)
// 			},
// 			wantUser: nil,
// 			wantErr:  ErrUserNotFound,
// 		},
// 		{
// 			name:   "empty id",
// 			userID: "",
// 			setupMock: func(m *mocks.MockUserRepository) {
// 				// No expectations - repo should not be called
// 			},
// 			wantUser:    nil,
// 			errContains: "id cannot be empty",
// 		},
// 		{
// 			name:   "repository error",
// 			userID: "user-1",
// 			setupMock: func(m *mocks.MockUserRepository) {
// 				m.EXPECT().
// 					GetByID(mock.Anything, "user-1").
// 					Return(nil, errors.New("db error"))
// 			},
// 			wantUser:    nil,
// 			errContains: "db error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			service, mockRepo := setupUserService(t)
// 			tt.setupMock(mockRepo)

// 			got, err := service.GetUser(context.Background(), tt.userID)

// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, err, tt.wantErr)
// 			} else if tt.errContains != "" {
// 				assert.ErrorContains(t, err, tt.errContains)
// 			} else {
// 				require.NoError(t, err)
// 			}

// 			assert.Equal(t, tt.wantUser, got)
// 		})
// 	}
// }

// // =============================================================================
// // Advanced Patterns
// // =============================================================================

// // Testing with mock.Anything for context or complex args
// func TestUserService_WithAnyContext(t *testing.T) {
// 	service, mockRepo := setupUserService(t)

// 	mockRepo.EXPECT().
// 		GetByID(mock.Anything, "user-1"). // Accept any context
// 		Return(&domain.User{ID: "user-1"}, nil)

// 	user, err := service.GetUser(context.Background(), "user-1")

// 	require.NoError(t, err)
// 	assert.Equal(t, "user-1", user.ID)
// }

// // Testing call counts
// func TestUserService_CallCounts(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	mockRepo.EXPECT().
// 		GetByID(ctx, "user-1").
// 		Return(&domain.User{ID: "user-1"}, nil).
// 		Times(3) // Expect exactly 3 calls

// 	// Call 3 times
// 	for i := 0; i < 3; i++ {
// 		_, _ = service.GetUser(ctx, "user-1")
// 	}
// 	// Assertions happen automatically via t.Cleanup()
// }

// // Testing with RunAndReturn for complex logic
// func TestUserService_RunAndReturn(t *testing.T) {
// 	service, mockRepo := setupUserService(t)
// 	ctx := context.Background()

// 	callCount := 0
// 	mockRepo.EXPECT().
// 		GetByID(mock.Anything, mock.Anything).
// 		RunAndReturn(func(ctx context.Context, id string) (*domain.User, error) {
// 			callCount++
// 			return &domain.User{ID: id, Name: "User " + id}, nil
// 		})

// 	_, _ = service.GetUser(ctx, "user-1")
// 	_, _ = service.GetUser(ctx, "user-2")

// 	assert.Equal(t, 2, callCount)
// }
