package protocol

import (
	"bufio"
	"fmt"
	"io"
)

// RoundTrip sends a request to a kafka broker and returns the response.
//
// The function expects that there were no other concurrent requests served by
// the connection wrapped by rw, and therefore uses a zero correlation ID.
func RoundTrip(w io.Writer, r *bufio.Reader, apiVersion int16, clientID string, msg Message) (Message, error) {
	const correlationID = 0
	if err := WriteRequest(w, apiVersion, correlationID, clientID, msg); err != nil {
		return nil, err
	}
	id, res, err := ReadResponse(r, int16(msg.ApiKey()), apiVersion)
	if err != nil {
		return nil, err
	}
	if id != correlationID {
		return nil, fmt.Errorf("correlation id mismatch (expected=%d, found=%d)", correlationID, id)
	}
	return res, nil
}