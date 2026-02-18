package utils

import "crypto/rand"

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenCode() string {
	const n = 8
	byts := make([]byte, 8)

	_, err := rand.Read(byts)
	if err != nil {
		return ""
	}

	for i := range n {
		byts[i] = characters[int(byts[i])%len(characters)]
	}

	return string(byts)
}
