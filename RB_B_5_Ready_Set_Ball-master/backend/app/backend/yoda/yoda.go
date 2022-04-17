package yoda

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"

	"io"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"github.com/julienschmidt/httprouter"
)

func paramsToMap(ps httprouter.Params) map[string]string {
	params := make(map[string]string, len(ps))
	for _, p := range ps {
		params[p.Key] = p.Value
	}
	return params
}

// SendClientError sends Yoda error to the client
func SendClientError(w http.ResponseWriter, err interface{}, status int) {
	log.Printf("Error: %v", err)
	if status < 100 {
		status = http.StatusBadRequest
	}
	w.WriteHeader(status)
	yErr := NewClientError(err)
	yRes := NewResponse(nil, yErr, true)
	json.NewEncoder(w).Encode(yRes)
}

// NewRequest creates a new Request
// Deprecated 0.0.1 Use [NewRequestWithSession]
func NewRequest(req *http.Request, params map[string]string) (*Request, error) {
	yReq := new(Request)
	if req.Method != "GET" {
		err := json.NewDecoder(req.Body).Decode(yReq)
		if err != nil {
			return nil, err
		}
		if yReq.Meta == nil || yReq.Result == nil {
			return nil, ClientError{yerr.BadRequest, "The request is not a proper yoda request"}
		}
	}

	yReq.Params = params
	return yReq, nil
}

// NewRequestFormData creates a new request with multipart data
func NewRequestFormData(req *http.Request, params map[string]string, sessVars map[interface{}]interface{}) (*Request, error) {
	yReq := new(Request)
	if req.Method != "GET" {
		req.ParseMultipartForm(constraints.MaxProfilePicSize * 2) //2x the max size to cover for other form data
		file, header, err := req.FormFile("file")
		if err != nil {
			return nil, err
		}

		yReq.Result = map[string]interface{}{
			"filesize": header.Size,
			"filename": header.Filename,
			"file":     file,
			"form":     map[string][]string(req.Form),
		}

		yReq.Meta = map[string]string{"type": "multipart/form"}
	}

	yReq.Session = sessVars
	yReq.Params = params
	return yReq, nil
}

// NewRequestWithSession creates a new request with session included. This should be used instead of [NewRequest]
func NewRequestWithSession(req *http.Request, params map[string]string, sessVars map[interface{}]interface{}) (*Request, error) {
	yReq, err := NewRequest(req, params)
	if err != nil {
		return nil, err
	}

	yReq.Session = sessVars
	return yReq, nil
}

// NewRequestWithSocket creates a new request with session and socket manager included
func NewRequestWithSocket(req *http.Request, params map[string]string, sessVars map[interface{}]interface{}, sm *sockets.Manager) (*Request, error) {
	yReq, err := NewRequest(req, params)
	if err != nil {
		return nil, err
	}

	yReq.Session = sessVars
	yReq.SocketManager = sm
	return yReq, nil
}

// NewResponse creates a new yoda response
func NewResponse(meta map[string]string, data interface{}, isError bool) *Response {
	if res, ok := data.(*Response); ok {
		return res
	}
	yRes := &Response{Meta: meta, Result: data}
	if isError {
		yRes.Meta = map[string]string{"Code": "200", "type": "Error"}
	}

	return yRes
}

// NewClientError creates a new yoda error
func NewClientError(err interface{}) *ClientError {
	if err == nil {
		return &ClientError{Code: 0, Message: "OK"}
	}

	var yErr ClientError

	switch v := err.(type) {
	case ClientError:
		yErr = err.(ClientError)
	case *ClientError:
		yErr = *v
	case *json.SyntaxError:
		yErr = ClientError{Code: yerr.BadRequest, Message: "The request wasn't understood"}
	default:
		switch v {
		case io.EOF:
			yErr = ErrInvalidRequest
		default:
			yErr = ClientError{Code: yerr.Unknown, Message: "An Unknown error was encountered"}
		}

	}

	return &yErr
}

func (yErr ClientError) Error() string {
	return fmt.Sprintf("Code: %d Message: %s", yErr.Code, yErr.Message)
}

// GetMeta returns a [map[string]string] of the metadata of a response
func (res Response) GetMeta() map[string]string {
	return res.Meta
}

// GetData returns a [map[string]interface{}] of the data of a response
func (res Response) GetData() interface{} {
	return res.Result
}

// GetMeta returns a [map[string]string] of the metadata of a request
func (req Request) GetMeta() map[string]string {
	return req.Meta
}

// GetData returns a [map[string]interface{}] of the data of a request
func (req Request) GetData() interface{} {
	return req.Result
}
