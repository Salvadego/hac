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

	// loggers, err := client.Log.GetCurrentLoggers(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for _, logger := range loggers {
	// 	fmt.Printf("%s: %s\n", logger.Name, logger.EffectiveLevel.StandardLevel)
	// }

	// levels, err := client.Log.GetLogLevels(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for _, level := range levels {
	// 	fmt.Printf("%s: (%s)\n", level.StandardLevel, level.DeclaringClass)
	// }

	resp, err := client.Log.ChangeLogLevel(ctx, "root", hac.LogLevelAll)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: (%s)\n", resp.LoggerName, resp.LevelName)
}
