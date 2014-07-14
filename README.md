# customerio #

`customerio` is a Go package for integration with the Customer.io
email service.

## Install ##

`go get github.com/joeshaw/customerio`

## API ##

The API is built on top of the
[Customer.io REST API](http://customer.io/docs/api/rest.html).

Full API docs are available on
[godoc](http://godoc.org/github.com/joeshaw/customerio).

```go
c := customerio.Client{
    SiteID: "my-site-id",
    APIKey: "my-api-key",
    HTTPClient: http.DefaultClient,
}

// Create a new customer with ID 5, email address customer@example.com,
// and a couple custom attributes
err := c.Identify("5", "customer@example.com", map[string]interface{}{
    "name": "Bob",
    "plan": "premium",
})
if err != nil {
    // uh oh
}

// Track an event named "purchased" on customer ID 5, with some custom
// event attributes.
err = c.Track("5", "purchased", map[string]interface{}{
    "price": 23.45,
}
if err != nil {
    // uh oh
}
```

You may pass in a `nil` map to either `Identify` or `Track`.
