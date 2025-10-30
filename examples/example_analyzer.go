package main

import (
	"context"
	"fmt"

	"github.com/Salvadego/hac/hac"
)

func main() {
	ctx := context.Background()

	client := hac.NewClient(&hac.Config{
		BaseURL:       "https://localhost:9002/hac",
		Username:      "admin",
		Password:      "nimda",
		SkipTLSVerify: true,
	})

	client.Auth.Login(ctx)

	resp, err := client.PKA.Analyze(ctx, hac.PKAnalyzeRequest{
		PKString: "1",
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("PK details: %+v\n", resp)
}
