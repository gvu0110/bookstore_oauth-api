package access_token

import (
	"net/http"
	"strings"

	"github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_oauth-api/repository/db"
	"github.com/gvu0110/bookstore_oauth-api/repository/rest"
	"github.com/gvu0110/bookstore_utils-go/rest_errors"
)

type Service interface {
	GetByID(string) (*access_token.AccessToken, rest_errors.RESTError)
	CreateAccessToken(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RESTError)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RESTError
	DeleteAccessToken(string, access_token.AccessTokenRequest) rest_errors.RESTError
}

type service struct {
	restUsersRepo rest.RESTUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RESTUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, rest_errors.RESTError) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, rest_errors.NewBadRequestRESTError("Invalid access token ID")
	}

	accessToken, err := s.dbRepo.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) CreateAccessToken(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RESTError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// TODO: Support both grant types: password and client_credentials
	// Authenticate the user against the Users API
	user, err := s.restUsersRepo.LoginUser(request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token
	at := access_token.GetNewAccessToken()
	at.Generate()
	at.UserID = user.ID

	// Save the new access token in Cassandra
	if err := s.dbRepo.CreateAccessToken(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RESTError {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}

func (s *service) DeleteAccessToken(accessTokenID string, request access_token.AccessTokenRequest) rest_errors.RESTError {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return rest_errors.NewBadRequestRESTError("Invalid access token ID")
	}

	if err := request.Validate(); err != nil {
		return err
	}

	user, err := s.restUsersRepo.LoginUser(request.Email, request.Password)
	if err != nil {
		return err
	}

	at, err := s.GetByID(accessTokenID)
	if err != nil {
		return err
	}

	if user.ID == at.UserID {
		if err := s.dbRepo.DeleteAccessToken(accessTokenID); err != nil {
			return err
		}
		return nil
	} else {
		return rest_errors.NewRESTError("Invalid credentials", http.StatusUnauthorized, "unauthorized", []interface{}{"This access token doesn't belong to the given credentials"})
	}
}
