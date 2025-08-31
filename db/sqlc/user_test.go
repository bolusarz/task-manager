package db

import (
	"testing"

	"github.com/bolusarz/task-manager/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		FirstName:    util.RandomString(10),
		LastName:     util.RandomString(10),
		Email:        util.RandomEmail(),
		PasswordHash: util.RandomString(12),
	}

	user, err := testQueries.CreateUser(t.Context(), arg)

	require.NoError(t, err)
	require.Equal(t, user.FirstName, arg.FirstName)
	require.Equal(t, user.LastName, arg.LastName)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.PasswordHash, arg.PasswordHash)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserById(t *testing.T) {
	user := createRandomUser(t)

	fetchedUser, err := testQueries.GetUserById(t.Context(), user.ID)
	require.NoError(t, err)

	require.Equal(t, user, fetchedUser)

	_, err = testQueries.GetUserById(t.Context(), -1)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	fetchedUser, err := testQueries.GetUserByEmail(t.Context(), user.Email)
	require.NoError(t, err)

	require.Equal(t, user, fetchedUser)

	_, err = testQueries.GetUserByEmail(t.Context(), "not-exists@gmail.com")
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserParams{
		ID:                user.ID,
		FirstName:         util.RandomString(11),
		LastName:          user.LastName,
		Email:             user.Email,
		ProfilePictureUrl: user.ProfilePictureUrl,
		IsEmailVerified:   user.IsEmailVerified,
	}

	updatedUser, err := testQueries.UpdateUser(t.Context(), arg)
	require.NoError(t, err)

	require.Equal(t, updatedUser.FirstName, arg.FirstName)
	require.Equal(t, updatedUser.LastName, arg.LastName)
	require.Equal(t, updatedUser.Email, arg.Email)
	require.Equal(t, updatedUser.PasswordHash, user.PasswordHash)

	arg = UpdateUserParams{
		ID:                -1,
		FirstName:         util.RandomString(11),
		LastName:          user.LastName,
		Email:             user.Email,
		ProfilePictureUrl: user.ProfilePictureUrl,
		IsEmailVerified:   user.IsEmailVerified,
	}

	updatedUser, err = testQueries.UpdateUser(t.Context(), arg)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	arg := UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: util.RandomString(20),
	}

	err := testQueries.UpdateUserPassword(t.Context(), arg)
	require.NoError(t, err)

	updatedUser, _ := testQueries.GetUserById(t.Context(), user.ID)

	require.NotEqual(t, user.PasswordHash, updatedUser.PasswordHash)
}

func TestDeleteUserPassword(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(t.Context(), user.ID)
	require.NoError(t, err)

	_, err = testQueries.GetUserById(t.Context(), user.ID)

	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
