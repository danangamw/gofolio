package auth

import "net/http"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// TODO: render login.html template
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: process login form, create session, redirect to /admin
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: destroy session, clear cookie, redirect to /login
}
