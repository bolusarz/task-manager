package api

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var validate = validator.New()

type User struct {
	Password string `validate:"strong"`
}

func TestIsPasswordStrong(t *testing.T) {
	tests := []struct {
		name    string
		data    User
		wantErr bool
	}{
		{
			name: "Invalid case: Length",
			data: User{
				Password: "234",
			},
			wantErr: true,
		},
		{
			name: "Invalid case: No Lowercase",
			data: User{
				Password: "BOLUWATIFE123",
			},
			wantErr: true,
		},
		{
			name: "Invalid case: No Uppercase",
			data: User{
				Password: "boluwatife@123",
			},
			wantErr: true,
		},
		{
			name: "Invalid case: No Digit",
			data: User{
				Password: "Boluwatife@",
			},
			wantErr: true,
		},
		{
			name: "Invalid case: No Special character",
			data: User{
				Password: "Boluwatife123",
			},
			wantErr: true,
		},
		{
			name: "Valid case",
			data: User{
				Password: "Boluwatife@123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type Payload struct {
	Email string `validate:"required,alpha,min=3,max=10,email,strong"`
}

func TestTransformValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		data          any
		errorMessages []string
	}{
		{
			name: "Invalid case: 'alpha'",
			data: struct {
				Name string `validate:"alpha"`
			}{
				Name: "12",
			},
			errorMessages: []string{
				"Name can only contain non numeric strings",
			},
		},
		{
			name: "Invalid case: 'min'",
			data: struct {
				Name string `validate:"min=3"`
			}{
				Name: "ab",
			},
			errorMessages: []string{
				"Name requires a min length of 3",
			},
		},
		{
			name: "Invalid case: 'strong'",
			data: struct {
				Password string `validate:"strong"`
			}{
				Password: "abcxyz.com",
			},
			errorMessages: []string{
				"Password is not strong enough",
			},
		},
		{
			name: "Invalid case: 'email'",
			data: struct {
				Email string `validate:"email"`
			}{
				Email: "abcxyz.com",
			},
			errorMessages: []string{
				"Email is not a valid email address",
			},
		},
		{
			name: "Invalid case: 'required'",
			data: struct {
				FirstName string `validate:"required"`
			}{
				FirstName: "",
			},
			errorMessages: []string{
				"FirstName is a required field",
			},
		},
		{
			name: "Invalid case: 'max'",
			data: struct {
				FirstName string `validate:"max=10"`
			}{
				FirstName: "boluwatifeadewusi",
			},
			errorMessages: []string{
				"FirstName exceeds the maximum length of 10",
			},
		},
		{
			name: "No Errors",
			data: struct {
				Email string `validate:"required,min=3,max=10,email"`
			}{
				Email: "abc@x.com",
			},
			errorMessages: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.data)
			errorMessages := TransformValidationErrors(err)

			if len(tt.errorMessages) == 0 {
				require.Empty(t, errorMessages)
				return
			}
			require.Equal(t, len(tt.errorMessages), len(errorMessages))

			for _, msg := range tt.errorMessages {
				require.Contains(t, errorMessages, msg)
			}
		})
	}
}
