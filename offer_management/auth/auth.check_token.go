package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// getUserIdFromToken maneja la obtención del userId a partir del token
func GetUserIdFromToken(token string) (string, error) {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "test" {
		return "18ad68bf-668c-48e6-ad90-a703d4add936", nil
	}

	url := fmt.Sprintf("%s/users/me", os.Getenv("USERS_PATH"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unauthorized")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	userId, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID")
	}

	return userId, nil
}

// processAuthorization verifica el token y obtiene el userId si es necesario
func ProcessAuthorization(r *http.Request) (string, error) {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "test" {
		return "18ad68bf-668c-48e6-ad90-a703d4add936", nil // Devuelve un valor falso cuando ENVIRONMENT está vacío
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("Forbidden")
	}

	return GetUserIdFromToken(token)
}
