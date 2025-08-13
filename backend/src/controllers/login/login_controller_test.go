package login

import (
	"backend/src/interfaces"
	"backend/src/model"
	"backend/src/security"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ============ Implementation of Types and Structs =============

type LoginData struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ApiResponse struct {
	Message 	string 		`json:"message"`
	Token		string		`json:"token"`
}

// ============ Implementation of Mocks and Stubs =============

type MockLoginRepository struct {
	users  	map[string]model.LoginUser
	failGet bool
}

func NewMockLoginRepository() *MockLoginRepository {
	return &MockLoginRepository{
		users:  make(map[string]model.LoginUser),
	}
}

func (m *MockLoginRepository) GetUserByEmail(email string) (model.LoginUser, error) {
	if m.failGet {
		return model.LoginUser{}, nil
	}

	user, exists := m.users[email]
	if !exists {
		return model.LoginUser{}, nil
	
	}
	return user, nil
}

func (m *MockLoginRepository) AddUser(user model.LoginUser) {
	m.users[user.Email] = user
}

func (m *MockLoginRepository) SetFailGet(shouldFail bool) {
	m.failGet = shouldFail
}

func setTestRepository(repo interfaces.LoginRepositoryInterface) {
	loginRepo = repo
	repoErr = nil
}

func restoreRepository() {
	loginRepo = nil
	repoErr = nil
}

// ============ Test Cases =============

// First you need a real user added in the database to work with
func TestLogin_Success(t *testing.T) {
	mockRepo := NewMockLoginRepository()

	hashedPassword, err := security.HashPassword("arrozdoce")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	testUser := model.LoginUser{
		ID:        	1,
		Email: 		"teste@gmail.com",
		Password: 	string(hashedPassword),
	}

	mockRepo.AddUser(testUser)
	setTestRepository(mockRepo)
	defer restoreRepository()

	loginData := LoginData{
		Email:		"teste@gmail.com",
		Password:	"arrozdoce",
	}

	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("failed to marshal login data: %v", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	
	Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		t.Logf("Response body: %s", rr.Body.String())
	}

	expectedMessage := "Login successful"
	if !strings.Contains(rr.Body.String(), expectedMessage) {
		t.Errorf("Expected response body to contain %q, got %q", expectedMessage, rr.Body.String())
	}

	var response ApiResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Token == "" {
		t.Error("Expected token in response, got empty string")
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	mockRepo := NewMockLoginRepository()
	setTestRepository(mockRepo)
	defer restoreRepository()

	invalidJSON := []byte(`"email": "test@gmail.com", "password":}`)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}