package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/bolusarz/task-manager/db/mock"
	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/bolusarz/task-manager/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.Comparepassword(e.password, arg.PasswordHash)
	if err != nil {
		return false
	}

	e.arg.PasswordHash = arg.PasswordHash
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password, err := util.RandomPassword(12)
	require.NoError(t, err)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		FirstName:    util.RandomString(10),
		LastName:     util.RandomString(10),
		PasswordHash: hashedPassword,
		Email:        util.RandomEmail(),
	}
	return
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		payload       map[string]any
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			payload: map[string]any{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(args, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "BadRequest",
			payload: map[string]any{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest: InvalidEmail",
			payload: map[string]any{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     util.RandomString(10),
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest: Invalid Password",
			payload: map[string]any{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"password":  util.RandomString(7),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest: User Already Exists",
			payload: map[string]any{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).Return(user, db.ErrUniqueViolation)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			payload: map[string]any{
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonBody, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var responseData SuccessResponse
	err = json.Unmarshal(data, &responseData)
	require.NoError(t, err)

	fmt.Println(responseData)

	gotUser, ok := responseData.Data.(map[string]any)
	require.True(t, ok)

	require.Equal(t, user.FirstName, gotUser["firstName"])
	require.Equal(t, user.LastName, gotUser["lastName"])
	require.Equal(t, user.Email, gotUser["email"])
	require.Equal(t, false, gotUser["isEmailVerified"])
	require.NotContains(t, gotUser, "passwordHash")
}
