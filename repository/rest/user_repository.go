package rest

import (
	"encoding/json"
	"time"

	"github.com/gvu0110/bookstore_oauth-api/domain/users"
	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

var (
	userRESTClient = rest.RequestBuilder{
		BaseURL: "localhost",
		Timeout: 100 * time.Microsecond,
	}
)

type RESTUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RESTError)
}

type usersRepository struct{}

func NewRepository() RESTUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RESTError) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := userRESTClient.Post("/users/login", request)
	// Timeout
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerRESTError("Invalid RESTClient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr errors.RESTError
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			return nil, errors.NewInternalServerRESTError("Invalid error interface then trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerRESTError("Error when trying to unmarshell user response")
	}
	return &user, nil
}
