// Package customerio is a wrapper around the Customer.io REST API,
// documented at http://customer.io/docs/api/rest.html
package customerio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context"
)

var urlPrefix = "https://track.customer.io/api/v1"

// Error represents an error from the Customer.io API, containing an
// HTTP status code returned by the request.  Errors are documented in
// the REST API here: http://customer.io/docs/api/rest.html#section-Errors
type Error struct {
	StatusCode int
}

// Returns the string representation of the error.  Implements the
// built-in error interface.
func (err *Error) Error() string {
	return fmt.Sprintf("customer.io API returned %d %s", err.StatusCode, http.StatusText(err.StatusCode))
}

// Client represents a client to the Customer.io REST API.
type Client struct {
	SiteID string
	APIKey string

	// The http.Client used to make connections to the Customer.io
	// REST API.  You may use your own or http.DefaultClient.
	HTTPClient *http.Client
}

// Identify creates or updates a customer.  id is a unique, non-email
// identifier for a customer.  The attrs map may be nil, or contain
// attributes which Customer.io can use to personalize triggered emails
// or affect the logic of who receives them.  The REST endpoint is
// documented here:
// http://customer.io/docs/api/rest.html#section-Creating_or_updating_customers
func (c *Client) Identify(id string, email string, attrs map[string]interface{}) error {
	if id == "" {
		return errors.New("id is required")
	}

	if attrs == nil {
		attrs = map[string]interface{}{}
	}

	if email != "" {
		attrs["email"] = email
	}

	data, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	return c.do("PUT", "/customers/"+id, data)
}

// Delete will remove a customer, and all their information from
// Customer.io.  The REST endpoint is documented here:
// http://customer.io/docs/api/rest.html#section-Deleting_customers
func (c *Client) Delete(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return c.do("DELETE", "/customers/"+id, nil)
}

// Track will send an event to Customer.io.  The attrs map may be nil,
// or contain any information to attach to this event.  These attributes
// can be used in triggers to control who should receive triggered
// email.  The REST endpoint is documented here:
// http://customer.io/docs/api/rest.html#section-Track_a_custom_event
func (c *Client) Track(id string, eventName string, attrs map[string]interface{}) error {
	if id == "" {
		return errors.New("id is required")
	}

	path := fmt.Sprintf("/customers/%s/events", id)

	return c.event(path, eventName, attrs)
}

// Sends an event email to a address that is not associated with a user in
// Customer.io.  The attrs may be nil, and recipient will be assigned to
// "recipient" within it.  The REST endpoint is documented here:
// http://customer.io/docs/invitation-emails.html
func (c *Client) TrackRecipient(recipient, eventName string, attrs map[string]interface{}) error {
	if attrs == nil {
		attrs = map[string]interface{}{}
	}

	if r, ok := attrs["recipient"]; ok && r != recipient {
		return errors.New("recipient would be overwritten in attrs")
	}

	attrs["recipient"] = recipient

	return c.event("/events", eventName, attrs)
}

func (c *Client) event(path, eventName string, attrs map[string]interface{}) error {
	jsonMap := map[string]interface{}{
		"name": eventName,
	}

	if attrs != nil {
		jsonMap["data"] = attrs
	}

	data, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}

	return c.do("POST", path, data)
}

func (c *Client) do(method, path string, data []byte) error {
	if c == nil {
		return nil
	}

	u := urlPrefix + path

	var reader io.Reader
	if data != nil {
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.SiteID, c.APIKey)

	if reader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &Error{resp.StatusCode}
	}

	return nil
}

type contextKey struct{}

// NewContext creates a new context.Context with the Client associated
func NewContext(ctx context.Context, c *Client) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

// FromContext returns the *Client associated with a context.Context,
// if any.  Otherwise a nil *Client is returned.
func FromContext(ctx context.Context) *Client {
	c, _ := ctx.Value(contextKey{}).(*Client)
	return c
}
