package access_token

import (
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {
	if expirationTime != 24 {
		t.Error("Expiration time should be 24 hours")
	}
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()
	if at.IsExpired() {
		t.Error("New access token should NOT be expired")
	}

	if at.AccessToken != "" {
		t.Error("New access token should NOT have undefined access token ID")
	}

	if at.UserId != 0 {
		t.Error("New access token should NOT have an associated user ID")
	}
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	if !at.IsExpired() {
		t.Error("Empty access token should be expired by default")
	}

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	if at.IsExpired() {
		t.Error("Access token expiring 3 hours from now should NOT be expired")
	}
}
