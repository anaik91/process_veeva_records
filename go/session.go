package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SessionRespose struct {
	SessionId string
}

func getSessionId(apiurl, apiversion, username, password string, target interface{}) error {
	url := apiurl + "/api/" + apiversion + "/auth"
	fmt.Println(url)
	method := "POST"
	payload := strings.NewReader("username=" + username + "&password=" + password)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}
