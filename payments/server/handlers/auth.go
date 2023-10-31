package handlers

import (
	"encoding/json"
	"net/http"
	"payments/models"
	"payments/service"
	"strconv"

	"github.com/zeebo/errs"
)

const (
	authorizationHeader = "Authorization-Bearer"
	userIdHeader        = "UserID"
)

var authErr = errs.Class("authorization handlers")

type AuthorizationHandlers struct {
	service service.Authorization
}

func NewAuthorizationHandlers(s service.Authorization) *AuthorizationHandlers {
	return &AuthorizationHandlers{service: s}
}

// SignUp parses a user info from the request, validates it, and creates the new user.
// Returns the id of the new user.
func (h *AuthorizationHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusBadRequest)
		return
	}

	if err := input.SignUpValidation(); err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusBadRequest)
		return
	}
	userId, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": userId,
	})
}

// Login parses email and password of the user from the request, validates, and
// returns a new token which contains the user ID.
func (h *AuthorizationHandlers) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusBadRequest)
		return
	}

	if err := input.LoginValidation(); err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusBadRequest)
		return
	}
	token, err := h.service.GenerateToken(r.Context(), input)
	if err != nil {
		NewErrorResponse(w, authErr.Wrap(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set(authorizationHeader, token)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}

func (h *AuthorizationHandlers) UserIdentity(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			NewErrorResponse(w, authErr.New("empty authorization header"), http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ParseToken(header)
		if err != nil {
			NewErrorResponse(w, authErr.Wrap(err), http.StatusUnauthorized)
			return
		}

		r.Header.Set(userIdHeader, strconv.Itoa(userId))
		handler(w, r)
	}
}
