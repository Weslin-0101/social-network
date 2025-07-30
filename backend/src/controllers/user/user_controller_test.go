package controllers

import (
	"backend/src/exceptions"
	"backend/src/interfaces"
	"backend/src/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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

func (m *MockUserRepository) CreateUser(user model.User) (model.User, error) {
	if m.failCreate {
		return model.User{}, fmt.Errorf("simulated create error")
	}

	user.ID = m.nextId
	m.users[m.nextId] = user
	m.nextId++
	return user, nil
}

func (m *MockUserRepository) GetAllUsers() ([]model.User, error) {
	if m.failGet {
		return nil, fmt.Errorf("simulated get error")
	}

	var users []model.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) GetUserByID(userID uint64) (model.User, error) {
	if m.failGet {
		return model.User{}, fmt.Errorf("simulated get error")
	}

	user, exists := m.users[userID]
	if !exists {
		return model.User{}, exceptions.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetUserByNickname(nickname string) (model.User, error) {
	if m.failGet {
		return model.User{}, fmt.Errorf("simulated get error")
	}

	for _, user := range m.users {
		if user.Nickname == nickname {
			return user, nil
		}
	}

	return model.User{}, exceptions.ErrUserNotFound
}

func (m *MockUserRepository) UpdateUserByID(userID uint64, user model.User) (model.User, error) {
	if m.failGet {
		return model.User{}, fmt.Errorf("simulated update error")
	}

	existingUser, exists := m.users[userID]
	if !exists {
		return model.User{}, exceptions.ErrUserNotFound
	}

	existingUser.Username = user.Username
	existingUser.Nickname = user.Nickname
	existingUser.Email = user.Email
	m.users[userID] = existingUser

	return existingUser, nil
}

func (m *MockUserRepository) DeleteUserByID(userID uint64) error {
	if m.failGet {
		return fmt.Errorf("simulated delete error")
	}

	_, exists := m.users[userID]
	if !exists {
		return exceptions.ErrUserNotFound
	}

	delete(m.users, userID)
	return nil
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

	if _, exists := response["user"]; !exists {
		t.Error("response should contain user data")
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "",
		Nickname: "Test User",
		Email: "teste@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	expectedMsg := "username is required"
	if !strings.Contains(rr.Body.String(), expectedMsg) {
		t.Errorf("expected body to contain %q, got %q", expectedMsg, rr.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if response["error"] == nil {
		t.Error("response should contain error message")
	}
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "Tester User",
		Nickname: "Test User",
		Email: "teste@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)
	
	req = httptest.NewRequest("GET", "/users", nil)
	rr = httptest.NewRecorder()

	GetAllUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var users []model.User
	err := json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(users) == 0 {
		t.Error("expected at least one user, got none")
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "Test User",
		Nickname: "Test User",
		Email: "email@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	var createResponse map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &createResponse)
	if err != nil {
		t.Errorf("failed to unmarshal create response: %v", err)
	}

	userID, exists := createResponse["user_id"]
	if !exists {
		t.Fatal("user_id not found in create response")
	}

	userIDStr := fmt.Sprintf("%.0f", userID.(float64));

	req = httptest.NewRequest("GET", "/users/"+userIDStr, nil)
	req = mux.SetURLVars(req, map[string]string{"userID": userIDStr})

	rr = httptest.NewRecorder()

	GetUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	expected := "User retrieved successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["user"] == nil {
		t.Error("response should contain user data")
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	userID := "999"
	req := httptest.NewRequest("GET", "/users/"+userID, nil)
	req = mux.SetURLVars(req, map[string]string{"userID": userID})
	rr := httptest.NewRecorder()

	GetUserByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetUserByNickname_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "Test User",
		Nickname: "testuser",
		Email: "test@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	nickname := user.Nickname
	req = httptest.NewRequest("GET", "/users/nickname/"+nickname, nil)
	req = mux.SetURLVars(req, map[string]string{"nickname": nickname})
	rr = httptest.NewRecorder()

	GetUserByNickname(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
		t.Logf("Response body: %s", rr.Body.String())
	}

	expected := "User retrieved successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["user"] == nil {
		t.Error("response should contain user data")
	}
}

func TestGetUserByNickname_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	nickname := "nonexistent"
	req := httptest.NewRequest("GET", "/users/nickname/"+nickname, nil)
	req = mux.SetURLVars(req, map[string]string{"nickname": nickname})
	rr := httptest.NewRecorder()

	GetUserByNickname(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUpdateUserByID_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "test",
		Nickname: "Test User",
		Email: "test@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	var createResponse map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &createResponse)
	if err != nil {
		t.Errorf("failed to unmarshal create response: %v", err)
	}

	userID, exists := createResponse["user_id"]
	if !exists {
		t.Fatal("user_id not found in create response")
	}

	userIDStr := fmt.Sprintf("%.0f", userID.(float64))
	
	updatedUser := model.User{
		Username: "updateduser",
		Nickname: "updatednick",
		Email: "test@gmail.com",
		Password: "password123",
	}
	updatedUserJSON, _ := json.Marshal(updatedUser)

	req = httptest.NewRequest("PUT", "/users/"+userIDStr, bytes.NewBuffer(updatedUserJSON))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"userID": userIDStr})
	rr = httptest.NewRecorder()

	UpdateUserByID(rr, req)
	
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestDeleteUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	user := model.User{
		Username: "testuser",
		Nickname: "Test User",
		Email: "test@gmail.com",
		Password: "password123",
	}

	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreateUser(rr, req)

	var createResponse map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &createResponse)
	if err != nil {
		t.Errorf("failed to unmarshal create response: %v", err)
	}
	userID, exists := createResponse["user_id"]
	if !exists {
		t.Fatal("user_id not found in create response")
	}

	userIDStr := fmt.Sprintf("%.0f", userID.(float64))
	req = httptest.NewRequest("DELETE", "/users/"+userIDStr, nil)
	req = mux.SetURLVars(req, map[string]string{"userID": userIDStr})
	rr = httptest.NewRecorder()

	DeleteUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "User deleted successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}