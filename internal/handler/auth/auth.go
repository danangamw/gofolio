package auth

import (
	"net/http"
	"time"

	"go-cms/internal/repository"
	"go-cms/internal/service"
	"go-cms/internal/session"
)

const sessionCookieName = "go_cms_session"

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
}

// Handler handles authentication HTTP requests (login page, login POST, logout).
type Handler struct {
	tmpl     Renderer
	userRepo *repository.UserRepository
	sessions session.Store
}

// New creates an AuthHandler with the given dependencies.
func New(tmpl Renderer, userRepo *repository.UserRepository, store session.Store) *Handler {
	return &Handler{
		tmpl:     tmpl,
		userRepo: userRepo,
		sessions: store,
	}
}

// loginViewModel is the data passed to the login template.
type loginViewModel struct {
	Title       string
	Description string
	ActiveMenu  string
	CSRFToken   string // placeholder — CSRF to be wired in Phase 5
	Error       string
}

// LoginPage renders the GET /login page.
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	vm := loginViewModel{Title: "Login — Go CMS"}
	h.tmpl.Render(w, "login", vm)
}

// Login handles POST /login: validates credentials, creates a session, sets cookie.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	renderError := func(msg string) {
		vm := loginViewModel{Title: "Login — Go CMS", Error: msg}
		w.WriteHeader(http.StatusUnauthorized)
		h.tmpl.Render(w, "login", vm)
	}

	user, err := h.userRepo.FindByUsername(r.Context(), username)
	if err != nil || user == nil {
		renderError("Username atau password salah.")
		return
	}

	ok, err := service.VerifyPassword(password, user.PasswordHash)
	if err != nil || !ok {
		renderError("Username atau password salah.")
		return
	}

	token, err := session.GenerateToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.sessions.Set(r.Context(), token, user.ID.String(), 0); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// Logout handles POST /admin/logout: destroys the session and clears the cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		h.sessions.Delete(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
