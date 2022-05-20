# go-web-framework

## Example

```go
package main

import (
	"net/http"
	"github.com/WengChaoxi/go-web-framework/framework"
)

func main() {
	core := framework.NewCore()
	core.Get("/", func(c *framework.Context) error {
		c.Json(200, "h3110 w0r1d")
		return nil
	})
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
```
