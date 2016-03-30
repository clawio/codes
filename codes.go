package codes

import (
	"fmt"
	"net/http"
	"net/url"
)

// A Code is an unsigned 32-bit error code.
type Code uint32

const (
	// To add new coded always add them in the end, to not break iota

	// Sucess indicates no error.
	Success Code = iota

	// InvalidToken is returned when the auth token is invalid or has expired
	InvalidToken

	// Unauthenticated is returned when authentication is needed for execution.
	Unauthenticated

	// BadAuthenticationData is returned when the authentication fails.
	BadAuthenticationData

	// BadInputData is returned when the input parameters are not valid.
	BadInputData

	// Internal is returned when there is an unexpected/undesired problem
	Internal
)

// String returns a string representation of the Code
func (c Code) String() string {
	switch c {
	case InvalidToken:
		return "Invalid or expired token"
	case Unauthenticated:
		return "Unauthenticated request"
	case BadAuthenticationData:
		return "Bad authentication data"
	case BadInputData:
		return "Bad input data"
	case Internal:
		return "Internal error. Please submit a query to the support team"
	default:
		return "FIXME: this should be a helpful message"
	}
}

// Response is a GitHub API response.  This wraps the standard http.Response
// returned from GitHub and provides convenient access to future things like
// pagination links.
type Response struct {
	*http.Response
}

// NewResponse creates a new Response for the provided http.Response.
func NewResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// An ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response `json:"-"` // HTTP response that caused this error
	*Err     `json:"error"` // more detail on individual errors
}

func NewErrorResponse(res *http.Response, e *Err) *ErrorResponse {
	response := &ErrorResponse{Response: res, Err: e}
	return response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Error)
}

// An Err reports more details on an individual error in an ErrorResponse.
type Err struct {
	Message string `json:"message"`
	Code    Code   `json:"code"`
}

// NewErr is a usefull function to create Errs with the corresponding Code message.
// If no message is passed, the detault code message will be used.
func NewErr(c Code, msg string) *Err {
	if msg == "" {
		msg = c.String()
	}
	return &Err{msg, c}
}

// Error() implements the Error interface.
func (e *Err) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Code.String())
}

// sanitizeURL redacts the token parameter from the URL which may be
// exposed to the user, specifically in the ErrorResponse error message.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("token")) > 0 {
		params.Set("token", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}
