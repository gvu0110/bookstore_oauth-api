package access_token

import (
	"strings"

	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
)

type Repository interface {
	GetByID(string) (*AccessToken, *errors.RESTError)
	CreateAccessToken(AccessToken) *errors.RESTError
	UpdateExpirationTime(AccessToken) *errors.RESTError
}

type Service interface {
	GetByID(string) (*AccessToken, *errors.RESTError)
	CreateAccessToken(AccessToken) *errors.RESTError
	UpdateExpirationTime(AccessToken) *errors.RESTError
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetByID(accessTokenID string) (*AccessToken, *errors.RESTError) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, errors.NewBadRequestRESTError("Invalid access token ID")
	}

	accessToken, err := s.repository.GetByID(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) CreateAccessToken(at AccessToken) *errors.RESTError {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.repository.CreateAccessToken(at)
}

func (s *service) UpdateExpirationTime(at AccessToken) *errors.RESTError {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.repository.UpdateExpirationTime(at)
}
