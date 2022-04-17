package yoda

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	yoda_error "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const uri string = "http://smartnotes.example"

type comp func(expected, actual interface{}) bool

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest("GET", uri, nil)
	yReq, err := NewRequest(req, map[string]string{})

	assert.Nilf(err, "Unexpected Error: %s", err)
	assert.Lenf(yReq.Params, 0, "The length of YodaRequest.Params should be zero")

	req = httptest.NewRequest("GET", uri+"/person/4", nil)
	yReq, err = NewRequest(req, map[string]string{"id": "4"})

	assert.Nilf(err, "Unexpected Error: %s", err)
	assert.Lenf(yReq.Params, 1, "The length of YodaRequest.Params should be zero")
	param := yReq.Params["id"]
	assert.Equal("4", param, "The id parameter on yReq.Params should be set")

	body := strings.NewReader(`{request: "invalidJSON"}`)
	req = httptest.NewRequest("POST", uri+"/person/4", body)
	yReq, err = NewRequest(req, map[string]string{"id": "4"})
	assert.Nilf(yReq, "The YodaRequest should be nil. Received %#v", yReq)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("Should return invalid request error\n\tExpected: %T\n\tActual: %T\n", new(json.SyntaxError), err)
	}

	body = strings.NewReader(`{"request": "invalid"}`)
	req = httptest.NewRequest("POST", uri+"/person/4", body)
	yReq, err = NewRequest(req, map[string]string{"id": "4"})
	if _, ok := err.(ClientError); !ok {
		t.Fatalf("Should return invalid request error\n\tExpected: %T\n\tActual: %T\n", ClientError{}, err)
	}
	if yErr, _ := err.(ClientError); yErr.Code != yoda_error.BadRequest || yErr.Message != "The request is not a proper yoda request" {
		t.Errorf("Should return invalid request error\n\tExpected: %v\n\tActual: %v\n", nil, yReq)
	}
	assert.Nilf(yReq, "The YodaRequest should be nil\n\tExpected")

	body = strings.NewReader(`{"meta" : {"structure" : "object", "type" : "Person"}, "result" : {"firstname": "Dan", "lastname" : "Ashig"}}`)
	req = httptest.NewRequest("POST", uri+"/person/4", body)
	yReq, err = NewRequest(req, map[string]string{"id": "4"})

	assert.NotNil(yReq, "The YodaRequest should be set")
	assert.Nilf(err, "Unexpected Error: %s", err)

	expected := map[string]string{"structure": "object", "type": "Person"}
	if ok := reflect.DeepEqual(yReq.Meta, expected); !ok {
		t.Errorf("The YodaRequest should have metadata\n\tExpected: %v\n\tActual: %v\n", expected, yReq.Meta)
	}
	expectedR := map[string]interface{}{"firstname": "Dan", "lastname": "Ashig"}
	if ok := reflect.DeepEqual(yReq.Result, expectedR); !ok {
		t.Errorf("The YodaRequest should have correct result\n\tExpected: %#v\n\tActual: %#v\n", expectedR, yReq.Result)
	}

}

func TestNewResponse(t *testing.T) {
	type args struct {
		meta    map[string]string
		data    interface{}
		isError bool
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		// TODO: Add test cases.
		{
			name: "Error meta is created",
			args: args{nil, &ClientError{Code: yoda_error.OK, Message: "Error"}, true},
			want: &Response{
				Meta:   map[string]string{"Code": "200", "type": "Error"},
				Result: &ClientError{Code: yoda_error.OK, Message: "Error"},
			},
		},
		{
			name: "Result is returned when there's no error",
			args: args{map[string]string{"type": "text"}, "Response to user", false},
			want: &Response{
				Meta:   map[string]string{"type": "text"},
				Result: "Response to user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(tt.args.meta, tt.args.data, tt.args.isError); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientError(t *testing.T) {
	type args struct {
		err interface{}
	}
	tests := []struct {
		name string
		args args
		want *ClientError
	}{
		// TODO: Add test cases.
		{
			name: "An yoda_error.OK error is created from a nil error",
			args: args{nil},
			want: &ClientError{Code: yoda_error.OK, Message: "OK"},
		},
		{
			name: "The same Yoda Error is created from a pointer to a Yoda Error",
			args: args{&ClientError{Code: yoda_error.InvalidID, Message: "Invalid ID"}},
			want: &ClientError{Code: yoda_error.InvalidID, Message: "Invalid ID"},
		},
		{
			name: "The same Yoda Error is created from a Yoda Error",
			args: args{ClientError{Code: yoda_error.InvalidID, Message: "Invalid ID"}},
			want: &ClientError{Code: yoda_error.InvalidID, Message: "Invalid ID"},
		},
		{
			name: "Yoda errors are created for JSON syntax error",
			args: args{&json.SyntaxError{}},
			want: &ClientError{Code: yoda_error.BadRequest, Message: "The request wasn't understood"},
		},
		{
			name: "The default Yoda Error is used for unknown errors",
			args: args{errors.New("An unknown error")},
			want: &ClientError{Code: yoda_error.Unknown, Message: "An Unknown error was encountered"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientError_Error(t *testing.T) {
	type fields struct {
		Code    int
		Message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "Error message is formatted correctly",
			fields: fields{Code: yoda_error.OK, Message: "This is an error message"},
			want:   "Code: 0 Message: This is an error message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yErr := ClientError{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
			}
			if got := yErr.Error(); got != tt.want {
				t.Errorf("ClientError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendClientError(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		err    interface{}
		status int
	}
	tests := []struct {
		name     string
		args     args
		response interface{}
		compare  comp
	}{
		// Test cases.
		{
			name: "Bad status codes default to http.yoda_error.BadRequest",
			args: args{
				w:      httptest.NewRecorder(),
				err:    ClientError{Code: yoda_error.InvalidRequest, Message: "Invalid status code"},
				status: 99,
			},
			response: &httptest.ResponseRecorder{Code: http.StatusBadRequest},
			compare: func(expected, actual interface{}) bool {
				expectedRes := expected.(*httptest.ResponseRecorder)
				actualRes := actual.(*httptest.ResponseRecorder)
				return expectedRes.Code == actualRes.Code
			},
		},
		{
			name: "Error is body of response",
			args: args{
				w:      httptest.NewRecorder(),
				err:    ClientError{Code: yoda_error.InvalidRequest, Message: "Resource Not Found"},
				status: 404,
			},
			response: ClientError{Code: yoda_error.InvalidRequest, Message: "Resource Not Found"},
			compare: func(expected, actual interface{}) bool {
				expectedErr := expected.(ClientError)
				expectedToMap := map[string]interface{}{"code": float64(expectedErr.Code), "message": expectedErr.Message}

				actualRes := actual.(*httptest.ResponseRecorder)
				yRes := &Response{}
				err := json.NewDecoder(actualRes.Body).Decode(yRes)
				if err != nil {
					t.Fatalf("Error while decoding response: %v\n", err)
				}

				return reflect.DeepEqual(yRes.Result, expectedToMap)
			},
		},
		{
			name: "Response is created as an error (contains error metadata)",
			args: args{
				w:      httptest.NewRecorder(),
				err:    ClientError{Code: yoda_error.InvalidRequest, Message: "Resource not found"},
				status: 404,
			},
			response: map[string]string{"Code": "200", "type": "Error"},
			compare: func(expected, actual interface{}) bool {
				expectedRes := expected.(map[string]string)
				actualRes := actual.(*httptest.ResponseRecorder)

				yRes := &Response{}
				err := json.NewDecoder(actualRes.Body).Decode(yRes)
				if err != nil {
					t.Fatalf("Error while decoding response: %v\n", err)
				}

				return reflect.DeepEqual(expectedRes, yRes.Meta)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := tt.args.w
			SendClientError(wr, tt.args.err, tt.args.status)
			if ok := tt.compare(tt.response, tt.args.w); !ok {
				t.Errorf("SendClientError wrote response = %v, \nExpected: %v", tt.args.w, tt.response)
			}
		})
	}
}

func Test_paramsToMap(t *testing.T) {
	type args struct {
		ps httprouter.Params
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		//Test cases
		{
			name: "Params are correctly converted to maps",
			args: args{
				ps: []httprouter.Param{
					httprouter.Param{Key: "id", Value: "1"},
					httprouter.Param{Key: "username", Value: "johndoe"},
					httprouter.Param{Key: "firstName", Value: "John"},
				},
			},
			want: map[string]string{"id": "1", "username": "johndoe", "firstName": "John"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paramsToMap(tt.args.ps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paramsToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
