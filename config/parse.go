package config

import (
	"encoding/json"
	"math/rand"
	"os"
)

func Parse(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	return cfg, err
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func WriteDefault(path string) (*Config, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		JWTKey:       randSeq(16),
		APIKey:       randSeq(16),
		CookieDomain: "127.0.0.1",
		CookieSecure: false,
		CookiePath:   "/",
	}
	err = json.NewEncoder(f).Encode(cfg)
	return cfg, err
}
