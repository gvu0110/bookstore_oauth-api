package access_token

import (
	"strings"

	"github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_oauth-api/repository/db"
	"github.com/gvu0110/bookstore_oauth-api/repository/rest"
	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
)

type Service interface {
	GetByID(string) (*access_token.AccessToken, *errors.RESTError)
	CreateAccessToken(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RESTError)
	UpdateExpirationTime(access_token.AccessToken) *errors.RESTError
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

func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, *errors.RESTError) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, errors.NewBadRequestRESTError("Invalid access token ID")
	}

	accessToken, err := s.dbRepo.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) CreateAccessToken(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RESTError) {
	// Authenticate the user against the Users API:
	if _, err := s.restUsersRepo.LoginUser(request.Email, request.Password); err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := access_token.GetNewAccessToken()
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.CreateAccessToken(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RESTError {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
