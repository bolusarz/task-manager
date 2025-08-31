package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/bolusarz/task-manager/util"
	"github.com/go-chi/render"
)

type createAccountPayload struct {
	FirstName string `json:"firstName" validate:"required,alpha,min=3,max=50"`
	LastName  string `json:"lastName" validate:"required,alpha,min=3,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,strong"`
}

type createAccountResponse struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	ProfilePicture  string `json:"profilePicture"`
}

func newCreateAccountResponse(user db.User) createAccountResponse {
	return createAccountResponse{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		ProfilePicture:  user.ProfilePictureUrl,
	}
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload createAccountPayload

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := s.validate.Struct(payload); err != nil {
		fieldErrors := TransformValidationErrors(err)
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("%s", fieldErrors[0])))
		return
	}

	hashedPassword, err := util.HashPassword(util.SanitizeInput(payload.Password))

	if err != nil {
		render.Render(w, r, ErrInternalServer())
		return
	}

	arg := db.CreateUserParams{
		FirstName:    util.SanitizeInput(payload.FirstName),
		LastName:     util.SanitizeInput(payload.LastName),
		Email:        util.SanitizeInput(payload.Email),
		PasswordHash: hashedPassword,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			render.Render(w, r, ErrInvalidRequestWithCode(fmt.Errorf("email already exists"), http.StatusConflict))
			return
		}
		render.Render(w, r, ErrInternalServer())
		return
	}

	render.Render(w, r, SuccessfulResponse(newCreateAccountResponse(user)))
}

type loginPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	ProfilePicture  string    `json:"profilePicture"`
	Token           string    `json:"token"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload loginPayload

	if err := render.DecodeJSON(r.Body, &payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := s.validate.Struct(payload); err != nil {
		fieldErrors := TransformValidationErrors(err)
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("%s", fieldErrors[0])))
		return
	}

	user, err := s.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid login credentials")))
			return
		}
		render.Render(w, r, ErrInternalServer())
		return
	}

	err = util.ComparePassword(payload.Password, user.PasswordHash)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid login credentials")))
		return
	}

	token, tokenPayload, err := s.tokenMaker.CreateToken(user.Email, s.config.AccessTokenDuration)
	if err != nil {
		render.Render(w, r, ErrInternalServer())
		return
	}

	response := loginResponse{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		ProfilePicture:  user.ProfilePictureUrl,
		Token:           token,
		ExpiresAt:       tokenPayload.ExpiresAt,
	}

	render.Render(w, r, SuccessfulResponse(response))
}
