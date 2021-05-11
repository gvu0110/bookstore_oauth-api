package access_token

import (
	"strings"
	"time"

	"github.com/gvu0110/bookstore_oauth-api/utils/errors"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() *errors.RESTError {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestRESTError("Invalid access token ID")
	}
	if at.UserID <= 0 {
		return errors.NewBadRequestRESTError("Invalid user ID")
	}
	if at.ClientID <= 0 {
		return errors.NewBadRequestRESTError("Invalid client ID")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestRESTError("Invalid expiration time")
	}
	return nil
}

func GetNewAccessToken() *AccessToken {
	return &AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

// Web frontend - ClientID: 123
// Android app - ClientID: 234
