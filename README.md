# hooks

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hatajoe/hooks)

hooks is a simple HTTP dispatcher like following.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/hatajoe/hooks"
	"github.com/hatajoe/hooks/github"
)

func main() {
	dispatcher := hooks.NewDispatcher(&github.Parser{})

    verifier := &github.VerifyMiddleware{
        secret: "webhook secret",
    }

	dispatcher.On("push", verifier.Verify(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("push event detected")
	}))

	if err := dispatcher.Listen("/webhooks", ":3000"); err != nil {
		panic(err)
	}
}
```
