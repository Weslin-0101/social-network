package controllers

import (
	"backend/src/interfaces"
	"backend/src/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ============ Implementation of Mocks and Stubs =============

type MockUserRepository struct {
	users 		map[uint64]model.User
	nextId 		uint64
	failCreate 	bool
	failGet 	bool
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint64]model.User),
		nextId: 1,
	}
}

func (m *MockUserRepository) CreateUser(user model.User) (uint64, error) {
	if m.failCreate {
		return 0, fmt.Errorf("simulated create error")
	}

	if user.ID == 0 {
		user.ID = m.nextId
		m.nextId++
	}

	m.users[user.ID] = user
	return user.ID, nil
}

func setTestRepository(repo interfaces.UserRepositoryInterface) {
	userRepo = repo
}

func restoreRepository() {
	userRepo = nil
}

// ============ Test Cases =============

func TestCreateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "testuser",
		Nickname: "Test User",
		Email: "test@gmail.com",
		Password: "password123",
		CreatedAt: time.Now(),
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	expectedMsg := "User created successfully"
	if !strings.Contains(rr.Body.String(), expectedMsg) {
		t.Errorf("expected body to contain %q, got %q", expectedMsg, rr.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if _, exists := response["user_id"]; !exists {
		t.Error("response should contain user_id")
	}
}

// func TestGetAllUsers(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodGet, "/users", nil)
// 	rr := httptest.NewRecorder()

// 	controllers.GetAllUsers(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rr.Code)
// 	}
// 	expected := "All users retrieved successfully"
// 	if !strings.Contains(rr.Body.String(), expected) {
// 		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
// 	}
// }

// func TestGetUserByID(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
// 	rr := httptest.NewRecorder()

// 	controllers.GetUserByID(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rr.Code)
// 	}
// 	expected := "User retrieved successfully"
// 	if !strings.Contains(rr.Body.String(), expected) {
// 		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
// 	}
// }

// func TestUpdateUserByID(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodPut, "/users/1", nil)
// 	rr := httptest.NewRecorder()

// 	controllers.UpdateUserByID(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rr.Code)
// 	}
// 	expected := "User updated successfully"
// 	if !strings.Contains(rr.Body.String(), expected) {
// 		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
// 	}
// }

// func TestDeleteUserByID(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
// 	rr := httptest.NewRecorder()

// 	controllers.DeleteUserByID(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rr.Code)
// 	}
// 	expected := "User deleted successfully"
// 	if !strings.Contains(rr.Body.String(), expected) {
// 		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
// 	}
// }