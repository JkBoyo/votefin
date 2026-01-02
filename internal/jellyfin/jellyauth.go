package jellyfin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var JellyfinAuthHeaderTemp string = "MediaBrowser Token=\"%s\", Client=\"Votefin\", Version=\"%s\", DeviceId=\"%s\", Device=\"votefin server\""

func addJellyfinAuthHeader(r *http.Request, token, userName string) {
	r.Header.Add(
		"Authorization",
		fmt.Sprintf(JellyfinAuthHeaderTemp,
			token,
			os.Getenv("RELEASE_VERSION"),
			hex.EncodeToString([]byte(userName))+os.Getenv("SERVER_ID"),
		),
	)
	r.Header.Add("Content-Type", "application/json")
}

func AuthenticateUser(userName, password string, con context.Context) (string, error) {
	params := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{Username: userName, Pw: password}

	client := http.DefaultClient

	reqData, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Error marshaling params to json: %s", err.Error())
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/Users/authenticatebyname", os.Getenv("JELLYFIN_URL")), bytes.NewReader(reqData))
	if err != nil {

		return "", fmt.Errorf("Error creating request: %s", err.Error())
	}
	defer req.Body.Close()

	addJellyfinAuthHeader(req, "", userName)

	AuthResp := struct {
		AccessToken string
	}{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making the request: %s", err.Error())
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(respData, &AuthResp)
	if err != nil {
		return "", err
	}

	return AuthResp.AccessToken, nil
}
