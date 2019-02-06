package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack/slackevents"
)

// ChallengeHandler handle Slack Event API challenge request
// If you want to make Slack Event API hook server, you can use this handler
func ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Read request body failed: %v", err), http.StatusBadRequest)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse event error: %v", err), http.StatusBadRequest)
		return
	}

	if eventsAPIURLVerificationEvent, ok := eventsAPIEvent.Data.(*slackevents.EventsAPIURLVerificationEvent); !ok {
		http.Error(w, "Unexpected event type detected: *slackevents.EventsAPIURLVerificationEvent is expected", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(eventsAPIURLVerificationEvent.Challenge))
	}
	return
}
