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

	resp, err := client.Flex.Execute(ctx, hac.FlexQuery{
		SQLQuery: "SELECT item_t0.PK, item_t0.p_sapordercode FROM orders item_t0",
		MaxCount: 5,
	}, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("Headers:", resp.Headers)
	for _, row := range resp.ResultList {
		fmt.Println(row)
	}
}
