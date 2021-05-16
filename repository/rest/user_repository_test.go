package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromAPI(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@example.com","password":"password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Invalid RESTClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInteraface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@example.com","password":"password"}`,
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     `{"status_code":"404","message":"Invalid login credentials","error":"NOT FOUND"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Invalid error interface then trying to login user", err.Message)
}

func TestLoginUserInvalidUserCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@example.com","password":"password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"status_code": "404", "message": "Invalid login credentials", "error": "NOT FOUND"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	assert.EqualValues(t, "Invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJSONResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@example.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "1", "first_name": "Adam", "last_name": "Vu", "email": "adam.vu@gmail.com", "date_created": "2006-01-02 15:04:05", "status": "active"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Error when trying to unmarshell user response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"email@example.com","password":"password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 1, "first_name": "Adam", "last_name": "Vu", "email": "adam.vu@gmail.com", "date_created": "2006-01-02 15:04:05", "status": "active"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "Adam", user.FirstName)
	assert.EqualValues(t, "Vu", user.LastName)
	assert.EqualValues(t, "adam.vu@gmail.com", user.Email)
	assert.EqualValues(t, "2006-01-02 15:04:05", user.DateCreated)
	assert.EqualValues(t, "active", user.Status)
}
