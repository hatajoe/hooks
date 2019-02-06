# hooks

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hatajoe/hooks)

hooks is a simple HTTP event dispatch handler.

## GitHub webhook server

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hatajoe/hooks"
	"github.com/hatajoe/hooks/github"
)

func main() {
	dispatcher := hooks.NewDispatcher(&github.EventParser{})

	verifier := github.NewVerifyMiddleware("GitHub webhook secret token")

	dispatcher.On("push", verifier.Verify(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("push event detected")
	}))

	if err := dispatcher.Listen("/webhooks", ":3000"); err != nil {
		log.Fatal(err)
	}
}
```

## Slack event API server

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hatajoe/hooks"
	"github.com/hatajoe/hooks/slack"
)

func main() {
	dispatcher := hooks.NewDispatcher(&slack.EventParser{
        ChallengeEventType: "challenge",
        VerifyToken: true,
        VerificationToken: "Slack event API secret token",    
    })

    dispatcher.On("challenge", slack.ChallengeHandler)
	dispatcher.On("app_mention", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("app_mention event detected")
	})

	if err := dispatcher.Listen("/webhooks", ":3000"); err != nil {
		log.Fatal(err)
	}
}
```
