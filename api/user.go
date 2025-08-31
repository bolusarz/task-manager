package api

import (
	"fmt"
	"log"
	"net/http"

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
		log.Println(err)
		render.Render(w, r, ErrInternalServer())
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
		log.Println(err)
		render.Render(w, r, ErrInternalServer())
		return
	}

	render.Render(w, r, SuccessfulResponse(newCreateAccountResponse(user)))
}
