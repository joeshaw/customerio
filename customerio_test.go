package customerio

import (
	"log"
	"net/http"
	"os"
	"testing"
)

var siteID, apiKey string

const testID = "customerio-go-test"
const testEmail = "customerio@example.com"

func init() {
	siteID = os.Getenv("TEST_CUSTOMER_SITE_ID")
	apiKey = os.Getenv("TEST_CUSTOMER_API_KEY")

	if siteID == "" || apiKey == "" {
		log.Fatal("TEST_CUSTOMER_SITE_ID and TEST_CUSTOMER_API_KEY environment variables must be set to run the tests")
	}
}

func TestAPI(t *testing.T) {
	c := Client{
		SiteID:     siteID,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}

	err := c.Identify(testID, testEmail, nil)
	if err != nil {
		t.Fatalf("c.Identify(testID, testEmail, nil) failed: %s", err)
	}

	defer func() {
		err := c.Delete(testID)
		if err != nil {
			t.Fatalf("c.Delete(testID) failed: %s", err)
		}
	}()

	err = c.Identify(testID, testEmail, map[string]interface{}{"foo": 42})
	if err != nil {
		t.Fatalf("c.Identify(testID, testEmail, map[string]interface{}{\"foo\": 42}) failed: %s", err)
	}

	err = c.Track(testID, "nil-attrs", nil)
	if err != nil {
		t.Fatalf("c.Track(testID, \"nil-attrs\", nil) failed: %s", err)
	}

	err = c.Track(testID, "with-attrs", map[string]interface{}{"bar": "baz", "quux": 3.14159})
	if err != nil {
		t.Fatalf("c.Track(testID, \"with-attrs\", map[string]interface{}{\"bar\": \"baz\", \"quux\": 3.14159}) failed: %s", err)
	}

	err = c.TrackRecipient(testEmail, "nil-attrs", nil)
	if err != nil {
		t.Fatalf("c.TrackRecipient(testEmail, \"nil-attrs\", nil) failed: %s", err)
	}

	err = c.TrackRecipient(testEmail, "with-attrs", map[string]interface{}{"bar": "baz", "quux": 3.14159})
	if err != nil {
		t.Fatalf("c.TrackRecipient(testEmail, \"with-attrs\", map[string]interface{}{\"bar\": \"baz\", \"quux\": 3.14159}) failed: %s", err)
	}
}

func TestNilClient(t *testing.T) {
	var c *Client = nil

	if err := c.Identify(testID, testEmail, nil); err != nil {
		t.Fatalf("Error on nil client c.Identify: %s", err)
	}

	if err := c.Track(testID, "test-event", nil); err != nil {
		t.Fatalf("Error on nil client c.Track: %s", err)
	}

	if err := c.Delete(testID); err != nil {
		t.Fatalf("Error on nil client c.Delete: %s", err)
	}

	if err := c.TrackRecipient(testEmail, "test-event", nil); err != nil {
		t.Fatalf("Error on nil client c.TrackRecipient: %s", err)
	}
}

func TestInvalidAPIArgs(t *testing.T) {
	c := Client{
		SiteID:     siteID,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}

	idRequired := "id is required"

	err := c.Identify("", testEmail, nil)
	if err == nil || err.Error() != idRequired {
		t.Fatalf("c.Identify expected error '%s'", idRequired)
	}

	err = c.Delete("")
	if err == nil || err.Error() != idRequired {
		t.Fatalf("c.Delete expected error '%s'", idRequired)
	}

	err = c.Track("", "test-event", nil)
	if err == nil || err.Error() != idRequired {
		t.Fatalf("c.Track expected error '%s'", idRequired)
	}

	err = c.TrackRecipient(testEmail, "test-event", map[string]interface{}{"recipient": "x" + testEmail})
	if err == nil || err.Error() != "recipient would be overwritten in attrs" {
		t.Fatalf("c.TrackRecipient expected error 'recipient would be overwritten in attrs'")
	}
}
