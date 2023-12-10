package api

import (
	"bytes"
	"dev.hackerman.me/artheon/l7-shared-launcher/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Login authenticates user with the API
func Login(email, password string) (string, error) {
	var (
		requestBody []byte
		err         error
	)

	requestBody, err = json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	if err != nil {
		log.Fatalln(err)
	}

	url := fmt.Sprintf("%s/auth/login", config.ApiUrl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var v map[string]string
	if err = json.Unmarshal(body, &v); err != nil {
		return "", err
	}

	if v["status"] == "error" {
		return "", fmt.Errorf("authentication error %d: %s\n", resp.StatusCode, v["message"])
	} else if v["status"] == "ok" {
		return v["data"], nil
	}

	return "", fmt.Errorf("authentication error %d: %s\n", resp.StatusCode, v["message"])
}

func FetchLatestAppRelease(token string, app string, platform string, configuration string, target string) ([]byte, error) {
	url := fmt.Sprintf("%s/releases/latest?app-id=%s&platform=%s&target=%s&configuration=%s", config.GetApi2Url(), app, platform, target, configuration)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the latest release: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to fetch the latest release, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
