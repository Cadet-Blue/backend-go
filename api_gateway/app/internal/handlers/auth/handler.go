package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Cadet-Blue/backend-go/api_gateway/internal/apperror"
	"github.com/Cadet-Blue/backend-go/api_gateway/internal/client/user_service"
	"github.com/Cadet-Blue/backend-go/api_gateway/pkg/jwt"
	"github.com/Cadet-Blue/backend-go/api_gateway/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	signInURL = "/api/signin"
	signUpURL = "/api/signup"
)

type Handler struct {
	Logger      logging.Logger
	UserService user_service.UserService
	JWTHelper   jwt.Helper
}

func (h *Handler) Register(router *httprouter.Router) {
	// router.HandlerFunc(http.MethodPost, signInURL, apperror.Middleware(h.SignIn))
	// router.HandlerFunc(http.MethodPut, signInURL, apperror.Middleware(h.SignIn))
	router.HandlerFunc(http.MethodPost, signUpURL, apperror.Middleware(h.SignUp))
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var dto user_service.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data")
	}

	err := h.UserService.Create(r.Context(), dto)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

// func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
// 	w.Header().Set("Content-Type", "application/json")

// 	var token []byte
// 	var err error
// 	switch r.Method {
// 	case http.MethodPost:
// 		defer r.Body.Close()
// 		var dto user_service.SigninUserDTO
// 		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 			return apperror.BadRequestError("failed to decode data")
// 		}
// 		u, err := h.UserService.GetByEmailAndPassword(r.Context(), dto.Email, dto.Password)
// 		if err != nil {
// 			return err
// 		}
// 		token, err = h.JWTHelper.GenerateAccessToken(u)
// 		if err != nil {
// 			return err
// 		}
// 	case http.MethodPut:
// 		defer r.Body.Close()
// 		var rt jwt.RT
// 		if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
// 			return apperror.BadRequestError("failed to decode data")
// 		}
// 		token, err = h.JWTHelper.UpdateRefreshToken(rt)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(token)

// 	return err
// }
