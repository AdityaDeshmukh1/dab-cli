package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginCommand() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Login to your DAB account",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "email", Aliases: []string{"e"}, Required: true},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Required: true},
		},

		Action: func(c *cli.Context) error {
			email := c.String("email")
			password := c.String("password")

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

			// Grabbing that cookie!
			for _, cookie := range resp.Cookies() {
				if cookie.Name == "session" {
					f, err := os.Create(".session")
					if err != nil {
						return fmt.Errorf("failed to write session: %v", err)
					}

					defer f.Close()
					f.WriteString(cookie.Value)
					fmt.Println("Login Successful!")
					return nil
				}
			}

			return fmt.Errorf("no session cookie found")
		},
	}
}
