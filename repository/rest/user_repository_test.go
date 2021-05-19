package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	os.Exit(m.Run())
}

func TestUserLoginAPIEndpointConst(t *testing.T) {
	assert.EqualValues(t, "https://localhost:8080/users/login", UserLoginAPIEndpoint)
}

func TestLoginUserTimeoutFromAPI(t *testing.T) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	defer httpmock.DeactivateAndReset()
	mockURL := UserLoginAPIEndpoint
	httpmock.RegisterResponder("POST", mockURL, nil)

	repository := usersRepository{}
	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Invalid RESTClient response when trying to login user", err.Message)
	assert.EqualValues(t, "INTERNAL SERVER ERROR", err.Error)
}

func TestLoginUserInvalidErrorInteraface(t *testing.T) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	defer httpmock.DeactivateAndReset()
	responseBody := `{"status_code":"404","message":"Invalid login credentials","error":"NOT FOUND"}`
	responder := httpmock.NewStringResponder(http.StatusInternalServerError, responseBody)
	mockURL := UserLoginAPIEndpoint
	httpmock.RegisterResponder("POST", mockURL, responder)

	repository := usersRepository{}
	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Invalid error interface then trying to login user", err.Message)
	assert.EqualValues(t, "INTERNAL SERVER ERROR", err.Error)
}

func TestLoginUserInvalidUserCredentials(t *testing.T) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	defer httpmock.DeactivateAndReset()
	responseBody := `{"status_code": 404, "message": "Invalid login credentials", "error": "NOT FOUND"}`
	responder := httpmock.NewStringResponder(http.StatusNotFound, responseBody)
	mockURL := UserLoginAPIEndpoint
	httpmock.RegisterResponder("POST", mockURL, responder)

	repository := usersRepository{}
	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	assert.EqualValues(t, "Invalid login credentials", err.Message)
	assert.EqualValues(t, "NOT FOUND", err.Error)
}

func TestLoginUserInvalidUserJSONResponse(t *testing.T) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	defer httpmock.DeactivateAndReset()
	responseBody := `{"id": "1", "first_name": "Adam", "last_name": "Vu", "email": "adam.vu@gmail.com", "date_created": "2006-01-02 15:04:05", "status": "active"}`
	responder := httpmock.NewStringResponder(http.StatusOK, responseBody)
	mockURL := UserLoginAPIEndpoint
	httpmock.RegisterResponder("POST", mockURL, responder)

	repository := usersRepository{}
	user, err := repository.LoginUser("email@example.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "Error when trying to unmarshall user response", err.Message)
	assert.EqualValues(t, "INTERNAL SERVER ERROR", err.Error)
}

func TestLoginUserNoError(t *testing.T) {
	httpmock.ActivateNonDefault(restClient.GetClient())
	defer httpmock.DeactivateAndReset()
	responseBody := `{"id": 1, "first_name": "Adam", "last_name": "Vu", "email": "adam.vu@gmail.com", "date_created": "2006-01-02 15:04:05", "status": "active"}`
	responder := httpmock.NewStringResponder(http.StatusOK, responseBody)
	mockURL := UserLoginAPIEndpoint
	httpmock.RegisterResponder("POST", mockURL, responder)

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
