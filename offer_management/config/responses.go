package config

import "time"

type OfferResponse struct {
	Id          string    `json:"id"`
	PostId      string    `json:"postId"`
	UserId      string    `json:"userId"`
	Description string    `json:"description"`
	Size        string    `json:"size"`
	Fragile     bool      `json:"fragile"`
	Offer       float64   `json:"offer"`
	CreatedAt   time.Time `json:"createdAt"`
}
