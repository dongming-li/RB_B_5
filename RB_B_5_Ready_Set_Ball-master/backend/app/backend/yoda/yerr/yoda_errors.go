package yerr

/**TODO:
* Make an error.go file
* Change name from ClientError to Error
 */
// YodaError error codes
const (
	OK = iota
	BadRequest
	UserNotFound
	InvalidRequest
	InvalidUpdateParameters //TODO: not being used
	InvalidID
	InvalidOperation
	InvalidUsernamePassword //TODO: change to InvalidCredentials
	InvalidRequestBody
	NoLogin
	NotLoggedIn
	InvalidSession
	InvalidParameter
	InviteAlreadySent
	InvalidGameName
	InvalidAgeRange
	InvalidTimeRange
	InvalidLocation
	InvalidSport
	InvalidHost
	Internal
	InvalidInvite
	AlreadyInGame
	InvalidArgument
	UserAlreadyExists
	InvalidFile
	InvalidFileName
	FileTooLarge
	ResourceNotFound
	MethodNotAllowed
	GameNotFound
	InvalidRating
	UnableToUpgrade
	InvalidPrivateGame
	GameHasStarted
	Unknown = 100
)
