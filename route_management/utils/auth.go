package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthenticator interface {
	ValidateAuth(c *gin.Context, config Config) bool
}

type AuthenticatorInstance struct{}

func GetAuthenticator() *AuthenticatorInstance { return &AuthenticatorInstance{} }

func (a AuthenticatorInstance) ValidateAuth(c *gin.Context, config Config) bool {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Status(403)
		return false
	}

	authValid := hasValidAuth(config.User_url, token)
	if !authValid {
		c.Status(http.StatusUnauthorized)
		return false
	}
	return true
}

func hasValidAuth(url, token string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/me", url), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Errorf("Error with status: %v, %v", resp.Status, resp.Body)
		return false
	}

	return true
}
