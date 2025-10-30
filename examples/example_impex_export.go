package main

import (
	"context"
	"fmt"
	"os"

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
INSERT_UPDATE Customer;uid[unique=true];name;sessionLanguage(isocode);sessionCurrency(isocode);groups(uid);creationtime[dateformat=yyyy-MM-dd HH:mm:ss];lastLogin[dateformat=yyyy-MM-dd HH:mm:ss];customerID;sapContactID
"#% impex.exportItemsFlexibleSearch(""SELECT {c.pk} FROM {Customer AS c JOIN Order AS o ON {o.user}={c.pk}} WHERE {o.creationtime} >= '2025-03-01 00:00:00' GROUP BY {c.pk} ORDER BY {c.creationtime} DESC"");"
`

	result, downloadURL, err := client.Impex.Export(ctx, hac.ImpexExportRequest{
		ScriptContent: script,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Export result:", result)
	fmt.Println("Download URL:", downloadURL)

	if downloadURL != "" {
		zipBytes, err := client.Impex.DownloadExportZip(downloadURL)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("export.zip", zipBytes, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Export ZIP saved as export.zip")
	}
}
