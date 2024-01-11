package DAO

type Settings struct {
	Proxy            string `yml:"proxy"`
	Cookie           string `yml:"cookie"`
	Agelimit         bool   `yml:"r-18" `
	Downloadposition string `yml:"downloadposition"`
	LikeLimit        int64  `yml:"minlikelimit"`
	Retry429         int    `yml:"retry429"`
	Downloadinterval int    `yml:"downloadinterval"`
	Retryinterval    int    `yml:"retryinterval"`
}
