package main

import (
	"fmt"
	"log"
	"errors"
	"encoding/json"
)

type LoginResponse struct {
	Success int
	Timestamp string	// ISO8601 and golang are not friends
	Authentication string	// token we want
	Expires string
}

func DoLogin(user string, pass string) (string, error) {
	var response LoginResponse

	api := SproutApi{""}
	resp, err := api.Post(login, fmt.Sprintf("identifier=%s&secret=%s", user, pass))
	if err != nil {
		return "", errors.New("Problem attempting to reach login API\n")
	}

	log.Print("Response: ", string(resp))
	// try unmarshal response, which may be a 400 or other error
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot interpret login response. [%s]", err.Error()))
	}
	if response.Success == 1 {
		log.Printf("Login successful. Token %s valid until %s", response.Authentication, response.Expires)
		return response.Authentication, nil
	}
	return "", errors.New("Login returned unsuccess (grammar is hard)")

}