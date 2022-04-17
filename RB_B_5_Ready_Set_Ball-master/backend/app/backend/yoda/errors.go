package yoda

import (
	"bytes"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
)

// ErrNoLogin is sent when a login fails other than [ErrInvalidCredentials]
// As of right now, it happens if it can't fetch a session for a login request
// or if a logged in user tries to login
var ErrNoLogin = ClientError{Code: yerr.NoLogin,
	Message: "User is already logged in"}

// ErrInvalidCredentials is sent if there is an attempted login with invalid credentials
var ErrInvalidCredentials = ClientError{Code: yerr.InvalidUsernamePassword,
	Message: "Invalid username or password"}

// ErrInvalidSession is sent if the session cannot be fetched
// TODO: security wise, does this error have too much info?
var ErrInvalidSession = ClientError{Code: yerr.InvalidSession,
	Message: "Invalid session"}

// ErrNotLoggedIn is sent if the session is not authenticated
var ErrNotLoggedIn = ClientError{Code: yerr.NotLoggedIn,
	Message: "User not logged in"}

// ErrInvalidData is sent when [request.Result] cannot be read
// most likely resulting from not being able to convert it to a
// map[string]string
var ErrInvalidData = ClientError{Code: yerr.InvalidRequestBody,
	Message: "The request could not be read"}

// ErrInvalidArgument is sent when taking in arguments to a route req.Params["argument"]
//and the arguments is not found or not correct
var ErrInvalidArgument = ClientError{Code: yerr.InvalidArgument,
	Message: "The argument is incorrect for this route"}

// ErrInvalidFile is sent when the file being uploaded doesn't match the
// expected MIME type
var ErrInvalidFile = ClientError{Code: yerr.InvalidFile,
	Message: "The file is invalid"}

// ErrInvalidFileName is sent when there is an attempt to upload a file
// with an invalid filename
var ErrInvalidFileName = ClientError{Code: yerr.InvalidFileName,
	Message: "The filename is invalid"}

// ErrFileTooLarge is sent when there is an attempt to upload a file
// that is larger than the constraint
var ErrFileTooLarge = ClientError{Code: yerr.FileTooLarge,
	Message: "The file was too large"}

// ErrResourceNotFound is sent when there is no handler to handle to the request url
var ErrResourceNotFound = ClientError{Code: yerr.ResourceNotFound,
	Message: "The resource was not found"}

// ErrMethodNotAllowed is sent when there is no handler to handle the url with the requested method
var ErrMethodNotAllowed = ClientError{Code: yerr.MethodNotAllowed,
	Message: "The method is not allowed"}

// ErrInvalidGame is sent when a user requests a game that is invalid or they don't have access to
var ErrInvalidGame = ClientError{Code: yerr.InvalidID,
	Message: "The game is invalid or the user does not have permissions"}

// ErrGameNotFound is sent when a requested game is not found
// e.g. when the current game of a user is requested and the user doesn't
// have a current game
var ErrGameNotFound = ClientError{Code: yerr.GameNotFound,
	Message: "The game does not exist"}

// ErrInvalidRequest is sent when the body of a request cannot be unmarshaled to a yoda Request
var ErrInvalidRequest = ClientError{Code: yerr.InvalidRequest,
	Message: "The request wasn't understood"}

// ErrCannotUpgrade is sent when request could not be upgraded to a socket connection
var ErrCannotUpgrade = ClientError{Code: yerr.UnableToUpgrade,
	Message: "The request could not be upgraded"}

// ErrSameUserInvite is sent when a user tries to send or cancel an invite to/from itself
var ErrSameUserInvite = ClientError{Code: yerr.InvalidInvite,
	Message: "The user cannot send/cancel a request to itself"}

// ErrInvalidParameters is sent when the atleast one of the expected url paramters is invalid
// Could only be caused internally if [httprouter] fails
var ErrInvalidParameters = ClientError{Code: yerr.InvalidParameter,
	Message: "Invalid request parameters"}

// ErrMissingData is sent when [request.Result] doesn't contain
// a data value or it has the wrong type
func ErrMissingData(name, expectedType string) ClientError {
	var buffer bytes.Buffer
	buffer.WriteString("Data ")
	buffer.WriteString(name)
	buffer.WriteString(" of type ")
	buffer.WriteString(expectedType)
	buffer.WriteString(" was not found")
	return ClientError{Code: yerr.InvalidRequestBody,
		Message: buffer.String()}

}
