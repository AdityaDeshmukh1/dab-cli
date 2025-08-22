package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(email, password string) error {
	payload := LoginPayload{Email: email, Password: password}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://dab.yeet.su/api/auth/login", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", string(body))
	}
	// Grab that session cookie!
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			f, err := os.Create(".session")
			if err != nil {
				return fmt.Errorf("failed to write session: %v", err)
			}
			defer f.Close()

			if _, err := f.WriteString(cookie.Value); err != nil {
				return fmt.Errorf("failed to write session value: %v", err)
			}

			fmt.Println("Login Successful!")
			return nil
		}
	}
	return fmt.Errorf("no session cookie found")
}
