package controllers

import (
	"backend/src/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rr := httptest.NewRecorder()

	controllers.CreateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	
	expected := "User created successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}

func TestGetAllUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	controllers.GetAllUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	expected := "All users retrieved successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}

func TestGetUserByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()

	controllers.GetUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	expected := "User retrieved successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}

func TestUpdateUserByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/users/1", nil)
	rr := httptest.NewRecorder()

	controllers.UpdateUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	expected := "User updated successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}

func TestDeleteUserByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()

	controllers.DeleteUserByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	expected := "User deleted successfully"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("expected body to contain %q, got %q", expected, rr.Body.String())
	}
}