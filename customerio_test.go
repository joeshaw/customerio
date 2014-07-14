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
}
