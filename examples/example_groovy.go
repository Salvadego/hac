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
print 'Hello from HAC Go client'
return 42
`

	resp, err := client.Groovy.Execute(ctx, hac.GroovyRequest{
		Script:     script,
		ScriptType: hac.ScriptGroovy,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Output:", resp.Output)
	fmt.Println("Result:", resp.Result)
	fmt.Println("Stacktrace:", resp.Stacktrace)
}
