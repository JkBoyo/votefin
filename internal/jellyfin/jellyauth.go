package jellyfin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

func AuthenticateUser(userName, password string, con context.Context) (string, error) {
	params := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{Username: userName, Pw: password}

	client := http.DefaultClient

	data, err := json.Marshal(params)
	if err != nil {
		//TODO: Handle error from marshaling the json
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/Users/AuthenticateByName", os.Getenv("JELLYFIN_URL")), bytes.NewReader(data))
	if err != nil {
		//TODO: Handle error from making the request
	}

	resp, err := client.Do(req)
	if err != nil {
		//TODO: Handle client response errors
	}

	b, err := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))

	token := "not gotten yet"

	return token, nil
}
