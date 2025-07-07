package turnstile

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func VerifyTurnstileToken(
	token string,
	secretKey string,
) (bool, error) {
	if token == "" || secretKey == "" {
		return false, nil
	}

	resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {secretKey},
		"response": {token},
	})
	if err != nil {
		return false, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	if !result.Success {
		return false, nil
	}

	return true, nil
}
