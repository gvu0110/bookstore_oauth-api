package rest

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/gvu0110/bookstore_oauth-api/domain/users"
	"github.com/gvu0110/bookstore_utils-go/rest_errors"
)

const (
	UserLoginAPIEndpoint = "http://localhost:8081/users/login"
)

var (
	restClient = resty.New()
)

type RESTUsersRepository interface {
	LoginUser(string, string) (*users.User, rest_errors.RESTError)
}

type usersRepository struct{}

func NewRepository() RESTUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, rest_errors.RESTError) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(UserLoginAPIEndpoint)

	// Timeout
	if err != nil {
		return nil, rest_errors.NewInternalServerRESTError("Invalid RESTClient response when trying to login user", err)
	}

	if response.StatusCode() > 299 {
		restErr, err := rest_errors.NewRESTErrorFromBytes(response.Body())
		if err != nil {
			return nil, rest_errors.NewInternalServerRESTError("Invalid error interface then trying to login user", err)
		}
		return nil, restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, rest_errors.NewInternalServerRESTError("Error when trying to unmarshall user response", err)
	}
	return &user, nil
}
