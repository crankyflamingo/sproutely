package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"encoding/json"
)

// TODO: Also generate a random UDID Hash value to be used (see api.go)
type CouponConfig struct {
	User string
	Pass string
	Token string
}

func main() {

	login := flag.Bool("login", false,
		"Used to regenerate token, by logging in with username and password. Tokens are typically valid for months")
	update := flag.Bool("update", false,
		"Will log into site, gather coupons, and apply to account")

	flag.Parse()
	config := CouponConfig{}

	// just to make sure our config file doesn't have a handle to it
	func() {
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatal("Unable to load config.json", err.Error())
		}

		defer f.Close()
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&config)
		if err != nil {
			log.Fatal("Unable to load config as valid json", err.Error())
		}
	}()

	switch {
	case *login == true:
		fmt.Print("Logging in to get new token (don't need to do this often!)\n")
		token, err := DoLogin(config.User, config.Pass)
		if err != nil {
			log.Fatal("Could not get new token: ", err.Error())
		}
		log.Print("Obtained new token, writing new config back to file")
		f, err := os.OpenFile("config.json", os.O_RDWR|os.O_TRUNC, 0)
		if err != nil {
			log.Fatal(fmt.Sprintf("Unable to write new token %s back to config", token))
		}
		defer f.Close()
		config.Token = token
		encoder := json.NewEncoder(f)
		err = encoder.Encode(&config)
		if err != nil {
			log.Print("Unable to store new token in config!")
		}
		break
	case *update == true:
		fmt.Print("Updating account with new coupons\n")
		if config.Token == "" {
			log.Print("No token in config, logging in first")
			token, err := DoLogin(config.User, config.Pass)
			if err != nil {
				log.Fatal("Could not get new token: ", err.Error())
			}
			config.Token = token
		}
		success := DoAccountUpdate(config.Token)
		if success != true {
			log.Print("Did not complete account update! (Token expired?) ")
			os.Exit(1)		// so we can make this obvious when cron'd
		}
		break
	default:
		flag.PrintDefaults()
	}
}