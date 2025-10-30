package main

import (
	"context"
	"fmt"

	"github.com/Salvadego/hac/hac"
)

func main() {
	// This example shows how to login to HAC.
	// This is required before you can use any other HAC API.

	ctx := context.Background()

	cfg := &hac.Config{
		BaseURL:  "https://localhost:9002/hac",
		Username: "admin",
		Password: "nimda",
		SkipTLSVerify: true,
	}

	client := hac.NewClient(cfg)

	if err := client.Auth.Login(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Logged in successfully!")
}
