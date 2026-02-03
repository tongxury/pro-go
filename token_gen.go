package main

import (
	"fmt"
	"store/pkg/krathelper"
)

func main() {
	// Generating a token for user ID 1
	token, err := krathelper.GenerateToken(1)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}
	fmt.Println(token)
}
