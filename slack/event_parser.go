package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack/slackevents"
)

// EventParser parse Slack event from http.Request
// EventParser is `EventParser` implementation
type EventParser struct {
	// ChallengeEventType is pattern string for receiving the Slack event API Challenge request
	// Slack event API send this special request due to verify server at the first step of Slack App registration
	// See also https://api.slack.com/events-api#begin
	ChallengeEventType string

	// VerifyToken is a flag that token verification enable or not
	// If it's true set, to verify token before event parsing
	VerifyToken bool

	// VerificationToken is Slack event API verification token
	// If VerifyToken is false, you don't have to set this
	VerificationToken string
}

// Parse returns the event string from http.Request
// If VerifyToken is true, verify token before event parsing
func (p EventParser) Parse(r *http.Request) (string, error) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	opts := []slackevents.Option{}
	if p.VerifyToken {
		opts = append(opts, slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.VerificationToken}))
	} else {
		opts = append(opts, slackevents.OptionNoVerifyToken())
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), opts...)
	if err != nil {
		return "", err
	}

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		return eventsAPIEvent.InnerEvent.Type, nil
	case slackevents.URLVerification:
		return p.ChallengeEventType, nil
	case slackevents.AppRateLimited:
		return "", fmt.Errorf("It has been reached to application rate limit")
	default:
		return "", fmt.Errorf("Unkown Slack event type detected: %s", eventsAPIEvent.Type)
	}
}
