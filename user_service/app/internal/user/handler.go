package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cadet-Blue/backend-go/user_service/internal/apperror"
	"github.com/Cadet-Blue/backend-go/user_service/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	// router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetUserByEmailAndPassword))
}

// func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) error {
// 	h.Logger.Info("GET USER BY EMAIL AND PASSWORD")
// 	w.Header().Set("Content-Type", "application/json")

// 	h.Logger.Debug("get email and password from URL")
// 	email := r.URL.Query().Get("email")
// 	password := r.URL.Query().Get("password")
// 	if email == "" || password == "" {
// 		return apperror.BadRequestError("invalid query parameters email or password")
// 	}

// 	user, err := h.UserService.GetByEmailAndPassword(r.Context(), email, password)
// 	if err != nil {
// 		return err
// 	}

// 	h.Logger.Debug("marshal user")
// 	userBytes, err := json.Marshal(user)
// 	if err != nil {
// 		return err
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write(userBytes)

// 	return nil
// }

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode create user dto")
	var crUser CreateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	if !h.UserService.CheckEmailExist(r.Context(), crUser.Email) {
		return apperror.BadRequestError("user with this email already exists")
	}

	userUUID, err := h.UserService.Create(r.Context(), crUser)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersURL, userUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}
