package jellyfin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var JellyfinAuthHeaderTemp string = "MediaBrowser Token=\"%s\", Client=\"Votefin\", Version=\"%s\", DeviceId=\"%s\", Device=\"votefin server\""

var JellyfinAuthError error = errors.New("User Not Authenticated")

type JellyfinUser struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
}

type JellyfinAuthResp struct {
	AccessToken string       `json:"AccessToken"`
	User        JellyfinUser `json:"User"`
}

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

func ValidateToken(token string) (JellyfinUser, error) {
	client := http.DefaultClient

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/Users/Me", os.Getenv("Jellyfin_URL")), nil)
	if err != nil {
		return JellyfinUser{}, fmt.Errorf("Error creating request: %v", err)
	}
	defer req.Body.Close()

	addJellyfinAuthHeader(req, token, "")

	resp, err := client.Do(req)
	if err != nil {
		return JellyfinUser{}, fmt.Errorf("Error performing request: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return JellyfinUser{}, JellyfinAuthError
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return JellyfinUser{}, fmt.Errorf("Error reading in response data: %v", err)
	}

	jellyUser := JellyfinUser{}

	err = json.Unmarshal(respData, &jellyUser)
	if err != nil {
		return JellyfinUser{}, fmt.Errorf("Error unmarshaling Json: %v", err)
	}

	return jellyUser, nil
}

func AuthenticateUser(userName, password string, con context.Context) (JellyfinAuthResp, error) {
	params := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{Username: userName, Pw: password}

	client := http.DefaultClient

	reqData, err := json.Marshal(params)
	if err != nil {
		return JellyfinAuthResp{}, fmt.Errorf("Error marshaling params to json: %s", err.Error())
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/Users/authenticatebyname", os.Getenv("JELLYFIN_URL")), bytes.NewReader(reqData))
	if err != nil {

		return JellyfinAuthResp{}, fmt.Errorf("Error creating request: %s", err.Error())
	}
	defer req.Body.Close()

	addJellyfinAuthHeader(req, "", userName)

	AuthResp := JellyfinAuthResp{}

	resp, err := client.Do(req)
	if err != nil {
		return JellyfinAuthResp{}, fmt.Errorf("Error making the request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return JellyfinAuthResp{}, JellyfinAuthError
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return JellyfinAuthResp{}, err
	}

	err = json.Unmarshal(respData, &AuthResp)
	if err != nil {
		return JellyfinAuthResp{}, err
	}

	return AuthResp, nil
}
