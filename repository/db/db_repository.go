package db

import (
	"github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
)

type DbRepository interface {
	GetByID(string) (*access_token.AccessToken, *errors.RESTError)
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetByID(string) (*access_token.AccessToken, *errors.RESTError) {
	return nil, errors.NewInternalServerError("Database connection not implemented yet")
}
