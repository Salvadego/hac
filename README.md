# HAC SDK Client for Go

A Go client for **SAP Hybris Administration Console (HAC)** that allows you to
authenticate, run FlexibleSearch queries, execute scripts, and handle Impex
imports/exports.

## Features

* Authentication (login/logout)
* FlexibleSearch query execution
* Groovy/JavaScript/Beanshell script execution
* Impex import/export and type/attribute fetching
* PK analysis

## Installation

```bash
go get github.com/Salvadego/hac/hac
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/Salvadego/hac/hac"
)

func main() {
	ctx := context.Background()

	cfg := &hac.Config{
		BaseURL:       "https://localhost:9002/hac",
		Username:      "admin",
		Password:      "nimda",
		SkipTLSVerify: true,
	}

	client := hac.NewClient(cfg)

	if err := client.Auth.Login(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Logged in successfully!")
}
```

---

## FlexibleSearch Example

```go
resp, err := client.Flex.Execute(ctx, hac.FlexQuery{
	// SQLQuery: "SELECT item_t0.PK, item_t0.p_sapordercode FROM orders item_t0", // Works too
	FlexibleSearchQuery: "SELECT {pk}, {sapOrderCode} FROM {Order}",
	User:     cfg.Username,
	MaxCount: 10,
}, nil)

if err != nil {
	panic(err)
}

fmt.Println("Raw SQL: ", resp.Query) // you can also see the raw SQL query
for _, row := range resp.ResultList {
	fmt.Println(row)
}
```

---

## Groovy Script Example

```go
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
```

---

## Impex Import Example

```go
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
```

---

## Impex Export Example

```go
exportScript := `
INSERT_UPDATE Customer;uid[unique=true];name
"#% impex.exportItemsFlexibleSearch(""SELECT {c.pk} FROM {Customer AS c}"");"
`

result, downloadURL, err := client.Impex.Export(ctx, hac.ImpexExportRequest{
	ScriptContent: exportScript,
})

if err != nil {
	panic(err)
}

fmt.Println("Export result:", result)
fmt.Println("Download URL:", downloadURL)

// You can also download the export zip file and save it to a file
zipBytes, err := client.Impex.DownloadExportZip(downloadURL)
if err != nil {
	panic(err)
}
fmt.Printf("Zip bytes: %d\n", len(zipBytes))

// save to file
filename := "export.zip"
err = os.WriteFile(filename, zipBytes, 0644)
if err != nil {
	panic(err)
}

```

---

## PK Analyzer Example

```go
pkResp, err := client.PKA.Analyze(ctx, hac.PKAnalyzeRequest{
	PKString: "1",
})

if err != nil {
	panic(err)
}

fmt.Printf("PK details: %+v\n", pkResp)
```

---

## Notes

* The client automatically handles CSRF tokens and session cookies.
- For more examples, you can check the `examples/` folder.

---
