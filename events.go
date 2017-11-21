package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	logger "github.com/blendlabs/go-logger"
)

const (
	// Flag is a logger event flag.
	Flag logger.Flag = "request"
	// FlagResponse is a logger event flag.
	FlagResponse logger.Flag = "request.response"
)

// EventOutgoing is a logger event for outgoing requests.
type EventOutgoing struct {
	ts  time.Time
	req *Meta
}

// Flag returns the event flag.
func (eo EventOutgoing) Flag() logger.Flag {
	return Flag
}

// Timestamp returns the event timestamp.
func (eo EventOutgoing) Timestamp() time.Time {
	return eo.ts
}

// Request returns the request meta.
func (eo EventOutgoing) Request() *Meta {
	return eo.req
}

// WriteText writes an outgoing request as text to a given buffer.
func (eo EventOutgoing) WriteText(tf logger.TextFormatter, buf *bytes.Buffer) error {
	buf.WriteString(fmt.Sprintf("%s %s", eo.req.Verb, eo.req.URL.String()))
	if len(eo.req.Body) > 0 {
		buf.WriteRune(logger.RuneNewline)
		buf.WriteString("request body")
		buf.WriteRune(logger.RuneNewline)
		buf.Write(eo.req.Body)
	}
	return nil
}

// MarshalJSON marshals an outgoing request event as json.
func (eo EventOutgoing) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"flag": eo.Flag(),
		"ts":   eo.ts,
		"req":  eo.req,
	})
}

// EventResponse is a response to outgoing requests.
type EventResponse struct {
	ts   time.Time
	req  *Meta
	res  *ResponseMeta
	body []byte
}

// Flag returns the event flag.
func (er EventResponse) Flag() logger.Flag {
	return FlagResponse
}

// Timestamp returns the event timestamp.
func (er EventResponse) Timestamp() time.Time {
	return er.ts
}

// Request returns the request meta.
func (er EventResponse) Request() *Meta {
	return er.req
}

// Response returns the response meta.
func (er EventResponse) Response() *ResponseMeta {
	return er.res
}

// Body returns the outgoing request body.
func (er EventResponse) Body() []byte {
	return er.body
}

// WriteText writes the event to a text writer.
func (er EventResponse) WriteText(tf logger.TextFormatter, buf *bytes.Buffer) error {
	buf.WriteString(fmt.Sprintf("%s %s %s", tf.ColorizeStatusCode(er.res.StatusCode), er.req.Verb, er.req.URL.String()))
	if len(er.body) > 0 {
		buf.WriteRune(logger.RuneNewline)
		buf.WriteString("response body")
		buf.WriteRune(logger.RuneNewline)
		buf.Write(er.body)
	}
	return nil
}

// MarshalJSON marshals an outgoing request event as json.
func (er EventResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"flag": er.Flag(),
		"ts":   er.ts,
		"req":  er.req,
		"res":  er.res,
		"body": er.body,
	})
}
