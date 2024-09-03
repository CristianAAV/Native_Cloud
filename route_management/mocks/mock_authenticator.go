package mocks

import (
	"route_management/utils"

	"github.com/gin-gonic/gin"
)

type MockAuthenticatorInvalid struct {
	Err     int
	Message string
}

func (m MockAuthenticatorInvalid) ValidateAuth(c *gin.Context, config utils.Config) bool {
	c.JSON(m.Err, m.Message)
	return false
}

func GetInvalidAuthenticator(err int, msg string) *MockAuthenticatorInvalid {
	return &MockAuthenticatorInvalid{Err: err, Message: msg}
}

type MockAuthenticatorValid struct{}

func (m MockAuthenticatorValid) ValidateAuth(c *gin.Context, config utils.Config) bool {
	return true
}

func GetValidAuthenticator() *MockAuthenticatorValid { return &MockAuthenticatorValid{} }
