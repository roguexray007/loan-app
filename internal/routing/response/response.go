package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/internal/provider"
	"github.com/roguexray007/loan-app/internal/trace"
)

// Response interface for response struct
type Response interface {
	SetError(err error) Response
	SetResponse(interface{}) Response
	SetHeaders(map[string]string) Response
	SetStatusCode(int) Response
	SetResponseType(string) Response
	AppendHeaders(map[string]string) Response
	Log()
	StatusCode() int
	ResponseType() string
	Body() interface{}
	ErrorBody() interface{}
	Error() error
	Headers() map[string]string
	GetResponseBody() interface{}
	String() string
	ModifyError(ctx context.Context)
}

// response struct will manage the response which has to be sent
type response struct {
	Response
	err          error
	body         interface{}
	statusCode   int
	responseType string
	headers      map[string]string
	ctx          *gin.Context
}

// NewResponse will create a new response instance
func NewResponse(ctx *gin.Context) Response {
	statusCode := http.StatusOK

	return &response{
		statusCode: statusCode,
		ctx:        ctx,
	}
}

// SetHeader will set the headers which has to be sent in the response
func (r *response) SetHeaders(headers map[string]string) Response {
	r.headers = headers

	return r
}

// AppendHeaders will append the given headers with the existing headers
func (r *response) AppendHeaders(headers map[string]string) Response {
	for key, val := range headers {
		r.headers[key] = val
	}

	return r
}

// Header will give back the headers which has to be set for the response
func (r *response) Headers() map[string]string {
	return r.headers
}

// SetStatusCode will set the code which will be displayed in with the error message
func (r *response) SetStatusCode(code int) Response {
	r.statusCode = code

	return r
}

// SetResponseType will defined the response format
func (r *response) SetResponseType(responseType string) Response {
	r.responseType = responseType

	return r
}

// SetError will set the error attribute with the given error
func (r *response) SetError(err error) Response {
	r.err = err

	return r
}

// SetResponse will set the response attribute with the given response
func (r *response) SetResponse(response interface{}) Response {
	r.body = response

	return r
}

// StatusCode will give the status code which has to be sent as response
func (r *response) StatusCode() int {
	return r.statusCode
}

// Body will provide the response body which was set earlier
func (r *response) Body() interface{} {
	return r.body
}

func (r *response) ErrorBody() interface{} {
	if provider.GetConfig(nil).App.Debug {
		return r.err
	}

	return r.err.Error()
}

// ResponseType will give the response format which has to be followed for the current response
func (r *response) ResponseType() string {
	return r.responseType
}

// Error will give the error instance set
func (r *response) Error() error {
	return r.err
}

// GetResponseBody will provide the response body which has to set to the client
// If there is no error set then response body will be same as Body
// In case of error response will the formatted map of error
func (r *response) GetResponseBody() interface{} {
	if r.err != nil {
		return r.ErrorBody()
	}

	return r.Body()
}

// Log will log the response with headers which will be sent to the client
func (r *response) Log() {
	data := map[string]interface{}{
		constants.StatusCode: r.StatusCode(),
		constants.Response:   r.GetResponseBody(),
		constants.Headers:    r.Headers(),
	}

	if r.ctx.Request.URL.String() != "/metrics" && r.ctx.Request.URL.String() != "/commit.txt" && r.ctx.Request.URL.String() != "/status" {
		fmt.Println(trace.Response+": %m", data)
	}
}

func (r *response) String() string {
	out, _ := json.Marshal(r)
	return string(out)
}

func (r *response) ModifyError(ctx context.Context) {
}
