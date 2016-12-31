package main

import (
	"log"
	"encoding/json"
	"fmt"
)

type CouponInstance struct {
	CouponId int `json:"coupon_id"`
	IsActive int `json:"is_active"`
	Featured int `json:"featured"`
	Targeted int `json:"targeted"`
	Brand string `json:"brand"`
	Department string `json:"department"`
	Tags []string `json:"tags"`
	ImageUrl string `json:"image_url"`
	DealType string `json:"deal_type"`
	OfferValue float32 `json:"offer_value"`
	OfferHeadline string `json:"offer_headline"`
	ShortDescription string `json:"short_description"`
	LongDescription string `json:"long_description"`
	StartTs string `json:"start_ts"`	// ISO8601 and golang are not friends
	EndTs string `json:"end_ts"`	// ISO8601 and golang are not friends
	ExpiryTs string `json:"expiry_ts"`	// ISO8601 and golang are not friends
}

type CouponResponse struct {
	Success int `json:"success"`
	Timestamp string `json:"timestamp"`	// ISO8601 and golang are not friends
	CouponCount int `json:"couponCount"`
	TotalResults int `json:"totalResults"`
	Coupons []CouponInstance `json:"coupons"`
}

func parseCoupons(jstr []byte) (CouponResponse, error) {
	var response CouponResponse
	err := json.Unmarshal(jstr, &response)
	if err != nil {
		log.Print("Can't parse coupon response")
		return response, err
	}
	fmt.Print("Couponcount ", response.CouponCount, response.TotalResults, "\n")
	return response, nil
}

// Queries available coupons, returns map of coupon id and description
func GetCoupons(api * SproutApi) *CouponResponse {

	response, err := api.Get(offers, nil)

	if err != nil {
		log.Print("Couldn't get API offers: ", err.Error())
		return nil
	}

	parsed, err := parseCoupons(response)
	if err != nil {
		log.Print("issue parsing coupons", err.Error())
		return nil
	}
	return &parsed
}

func redeemCoupon(api * SproutApi, id int) bool {
	response, err := api.Post(claim, fmt.Sprintf("coupon_ids=%d", id))
	if err != nil {
		log.Print("Problem submitting coupon ", id, " ", err.Error())
		return false
	}
	log.Print("Response ", string(response))
	return true
}


func DoAccountUpdate(token string) bool {
	api := SproutApi{token: token}

	var savings float32
	coupons := GetCoupons(&api)
	if coupons == nil {
		log.Print("Error fetching coupons")
		return false
	}
	log.Printf("Fetched %d coupons", coupons.CouponCount)
	if coupons.Success  == 1 && coupons.CouponCount > 0 {
		for _, coupon := range coupons.Coupons {
			if redeemCoupon(&api, coupon.CouponId) {
				log.Print("Redeemed coupon for ", coupon.OfferHeadline)
				savings += coupon.OfferValue
			}
		}
	}
	log.Print("Total savings: ", savings)
	return true
}