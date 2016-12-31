package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"errors"
	"strings"
)

var (
	// API endpoints we care about
	login = "https://api.mercatustechnologies.com/v1/auth/token"
	offers = "https://api.mercatustechnologies.com/v1/coupon/getoffers"
	claim = "https://api.mercatustechnologies.com/v1/coupon/savecoupon"
)

var (
	// Headers required(?) for the app
	CouponHeaders = map[string]string{
	"X-Mct-Apikey": "f9c073403c65640f7bddae0996f746b8",	// This is same for various devices for specific version of App
	"X-UDID-Hash" : "406e9a66b47d28d5c01a458591358a22",  	// This is "unique" to the app install it seems. If you're reading this, you may want to generate a random md5 and replace it, or better yet, intercept your real phone's traffic and use those values.
	"User-Agent" : "sprouts/3.7.7 (iPhone OS, 10.2, iPhone, Screen/320x568/2.00)",
	"Srn-Auth-Token" : "",
	}

	postHeaders = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
	}
)

type SproutApi struct {
	token string
}

func (s * SproutApi) Get (endpoint string, params *map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		add := req.URL.Query()
		for h, i := range *params {
			add.Add(h, i)
		}
		req.URL.RawQuery = add.Encode()
	}

	CouponHeaders["Srn-Auth-Token"] = s.token

	for k,v := range CouponHeaders {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Problem with request")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Print("Response not 200", resp.Status)
		return nil, errors.New("Status code not 200")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Problem reading response %s", err.Error())
		fmt.Print(resp.Body)
	}
	return data, nil
}


func (s * SproutApi) Post(endpoint string, body string) ([]byte, error) {

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	CouponHeaders["Srn-Auth-Token"] = s.token
	// for logging in to get another token
	if s.token == "" {
		delete(CouponHeaders, "Srn-Auth-Token")
	}

	for k,v := range postHeaders {
		req.Header.Add(k, v)
	}
	for k,v := range CouponHeaders {
		req.Header.Add(k, v)
	}

	// the app doesn't use/send content length, so we'll look a bit differet traffic-wise,
	// but we seemingly can't avoid it
	// Potential solution: proxy and filter the outgoing traffic of the HTTPS transaction
	// would need to share our TLS structure
	// req.ContentLength = 0

	// remove adding of 'Accept-Encoding: gzip' to POST
	client := &http.Client{}
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	transport.DisableCompression = true
	client.Transport = transport

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("Problem with request")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Print("Response not 200", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Problem reading response %s", err.Error())
		return nil, errors.New(fmt.Sprint("Cannot read response: ", resp.Status))
	}
	return data, nil
}