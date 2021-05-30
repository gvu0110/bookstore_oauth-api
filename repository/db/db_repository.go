package db

import (
	"github.com/gocql/gocql"
	"github.com/gvu0110/bookstore_oauth-api/clients/cassandra"
	"github.com/gvu0110/bookstore_oauth-api/domain/access_token"
	"github.com/gvu0110/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken       = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken    = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpirationTime = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetByID(string) (*access_token.AccessToken, rest_errors.RESTError)
	CreateAccessToken(access_token.AccessToken) rest_errors.RESTError
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RESTError
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetByID(id string) (*access_token.AccessToken, rest_errors.RESTError) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserID,
		&result.ClientID,
		&result.Expires,
	); err != nil {
		if err.Error() == gocql.ErrNotFound.Error() {
			return nil, rest_errors.NewNotFoundRESTError("No access token found with the given ID")
		}
		return nil, rest_errors.NewInternalServerRESTError("Error when trying to get access token by ID", err)
	}
	return &result, nil
}

func (r *dbRepository) CreateAccessToken(at access_token.AccessToken) rest_errors.RESTError {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserID,
		at.ClientID,
		at.Expires,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerRESTError("Error when trying to create access token", err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RESTError {
	if err := cassandra.GetSession().Query(queryUpdateExpirationTime,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerRESTError("Error when trying to update access token expiration time", err)
	}
	return nil
}
