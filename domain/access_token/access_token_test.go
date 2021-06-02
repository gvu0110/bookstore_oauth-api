package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "Expiration time should be 24 hours")
	assert.EqualValues(t, "password", grantTypePassword)
	assert.EqualValues(t, "client_credentials", grantTypeClientCredentials)
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	assert.False(t, at.IsExpired(), "New access token should NOT be expired")
	assert.EqualValues(t, "", at.AccessToken, "New access token should NOT have undefined access token ID")
	assert.True(t, at.UserID == 0, "New access token should NOT have an associated user ID")
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.IsExpired(), "Empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "Access token expiring 3 hours from now should NOT be expired")
}

func TestValidateAccessTokenInvalidID(t *testing.T) {
	at := AccessToken{
		AccessToken: "",
		UserID:      123,
		ClientID:    234,
		Expires:     345,
	}
	err := at.Validate()
	assert.EqualValues(t, "Invalid access token ID", err.Message())
	assert.EqualValues(t, 400, err.StatusCode())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestValidateAccessTokenInvalidUserID(t *testing.T) {
	at := AccessToken{
		AccessToken: "123abc",
		UserID:      -1,
		ClientID:    234,
		Expires:     345,
	}
	err := at.Validate()
	assert.EqualValues(t, "Invalid user ID", err.Message())
	assert.EqualValues(t, 400, err.StatusCode())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestValidateAccessTokenInvalidClientID(t *testing.T) {
	at := AccessToken{
		AccessToken: "123abc",
		UserID:      123,
		ClientID:    -1,
		Expires:     345,
	}
	err := at.Validate()
	assert.EqualValues(t, "Invalid client ID", err.Message())
	assert.EqualValues(t, 400, err.StatusCode())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestValidateAccessTokenInvalidExpirationTime(t *testing.T) {
	at := AccessToken{
		AccessToken: "123abc",
		UserID:      123,
		ClientID:    234,
		Expires:     -1,
	}
	err := at.Validate()
	assert.EqualValues(t, "Invalid expiration time", err.Message())
	assert.EqualValues(t, 400, err.StatusCode())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestValidateAccessTokenNoError(t *testing.T) {
	at := AccessToken{
		AccessToken: "123abc",
		UserID:      123,
		ClientID:    234,
		Expires:     345,
	}
	err := at.Validate()
	assert.Nil(t, err)
}

func TestValidateAccessTokenRequestInvalidGrantType(t *testing.T) {
	request := AccessTokenRequest{
		GrantType: "abc",
	}
	err := request.Validate()
	assert.EqualValues(t, "Invalid grant type parameter", err.Message())
	assert.EqualValues(t, 400, err.StatusCode())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestValidateAccessTokenRequestNoError(t *testing.T) {
	request := AccessTokenRequest{
		GrantType: grantTypePassword,
	}
	err := request.Validate()
	assert.Nil(t, err)

	request = AccessTokenRequest{
		GrantType: grantTypeClientCredentials,
	}
	err = request.Validate()
	assert.Nil(t, err)
}

func TestGenerate(t *testing.T) {
	at := AccessToken{
		UserID:  123,
		Expires: 345,
	}
	at.Generate()
	assert.EqualValues(t, "0f21432ecbf12034c99350948c7a0726", at.AccessToken)
}
