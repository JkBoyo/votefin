package jellyfin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

var JellyfinAuthHeaderTemp string = "MediaBrowser Token=\"%s\", Client=\"Votefin\", Version=\"%s\", DeviceId=\"%s\", Device=\"votefin server\""

func addJellyfinAuthHeader(r *http.Request, token, userName string) {
	r.Header.Add(
		"Authorization",
		fmt.Sprintf(JellyfinAuthHeaderTemp,
			token,
			os.Getenv("RELEASE_VERSION"),
			hex.EncodeToString([]byte(userName))+"12345",
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

	data, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Error marshaling params to json: %s", err.Error())
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/Users/authenticatebyname", os.Getenv("JELLYFIN_URL")), bytes.NewReader(data))
	if err != nil {

		return "", fmt.Errorf("Error creating request: %s", err.Error())
	}

	addJellyfinAuthHeader(req, "", userName)

	reqDat, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", fmt.Errorf("Error dumping the request: %s", err.Error())
	}
	fmt.Println(string(reqDat))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making the request: %s", err.Error())
	}
	defer resp.Body.Close()

	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return "", fmt.Errorf("Error dumping the response: %s", err.Error())
	}
	fmt.Println(string(b))

	token := "not gotten yet"

	return token, nil
}
