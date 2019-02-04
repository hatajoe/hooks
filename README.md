# hooks

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hatajoe/hooks)

hooks is a simple HTTP dispatcher like following.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/hatajoe/hooks"
)

func main() {
	dispatcher := hooks.NewDispatcher(func(r *http.Request) (string, error) {
		event := r.Header.Get("X-GitHub-Event")
		if event == "" {
			return "", fmt.Errorf("X-GitHub-Event is empty")
		}
		return event, nil
	})

	dispatcher.On("push", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("push event detected")
	})

	if err := dispatcher.Listen("/webhooks", ":3000"); err != nil {
		panic(err)
	}
}
```
