package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/gvu0110/bookstore_utils-go/encryption"
	"github.com/gvu0110/bookstore_utils-go/rest_errors"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	ClientID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

// Web frontend - ClientID: 123
// Android app - ClientID: 234

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`

	// Used for "grant_type": "password"
	Email    string `json:"email"`
	Password string `json:"password"`

	// Used for "grant_type": "client_credentials"
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessToken) Validate() rest_errors.RESTError {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.NewBadRequestRESTError("Invalid access token ID")
	}
	if at.UserID <= 0 {
		return rest_errors.NewBadRequestRESTError("Invalid user ID")
	}
	if at.ClientID <= 0 {
		return rest_errors.NewBadRequestRESTError("Invalid client ID")
	}
	if at.Expires <= 0 {
		return rest_errors.NewBadRequestRESTError("Invalid expiration time")
	}
	return nil
}

func (at *AccessTokenRequest) Validate() rest_errors.RESTError {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return rest_errors.NewBadRequestRESTError("Invalid grant type parameter")
	}
	return nil
}

func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = encryption.GetMD5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
