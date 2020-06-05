package utils

import (
	"testing"
	"time"
)

func assertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}

const samplePrincipalId = "65fc9e38-fc5d-4f0c-bc4d-e1a05493b4a8"
const sampleToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI1MmVkMGQxOGZhNWUzMzJmNmU2N2ZhZWM3NzQzYzE4ODYwYmMzODQ0Iiwic3ViIjoiNjVmYzllMzgtZmM1ZC00ZjBjLWJjNGQtZTFhMDU0OTNiNGE4IiwiYXVkIjoiNzUxMjNjNzktZjVlYS00Yjg1LThhMGQtNDViNGU0Zjg4ZGJiIiwiaWF0IjoxNTkxMTk0Mzg4LCJuYmYiOjE1OTExOTQzODgsImV4cCI6MTU5MTE5Nzk4OCwic2NvcGUiOiJhZ2VudCIsInNoYXJkIjoxfQ.rnK5RQkOKc1YnKAuyfRTMqLIJwAo7AG2Hy6fZC910D0"
const sampleExpiry = 1591197988

func TestGetTokenPrincipal(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		principalId, err := GetAccessTokenPrincipalId(sampleToken)
		if err != nil {
			t.Errorf("Got error: %q", err)
		}
		assertEquals(t, samplePrincipalId, principalId)
	})

	t.Run("invalid token", func(t *testing.T) {
		principalId, _ := GetAccessTokenPrincipalId("not-a token")
		assertEquals(t, "", principalId)
	})
}

func TestGetTokenExpiry(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		expiry, err := GetAccessTokenExpiry(sampleToken)
		if err != nil {
			t.Errorf("Got error: %q", err)
		}
		assertEquals(t, time.Unix(sampleExpiry, 0), expiry)
	})

	t.Run("invalid token", func(t *testing.T) {
		expiry, _ := GetAccessTokenExpiry("not-a token")
		assertEquals(t, time.Unix(0, 0), expiry)
	})
}
