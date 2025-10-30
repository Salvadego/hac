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

	script := `
INSERT_UPDATE Product;code[unique=true];name[lang=en]
; testProduct ; Test Product
`

	result, err := client.Impex.Import(ctx, hac.ImpexImportRequest{
		ScriptContent: script,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Impex import result:", result)
}
