package constraints

// Application contraints
const (
	//User credential coontraints
	MinUsernameLength = 4
	MaxUsernameLength = 30

	MinNameLength = 2
	MaxNameLength = 64

	MinPasswordLength = 10
	MaxPasswordength  = 128

	MaxPasswordDictionaryLength = 3

	MaxLoginAttempts = 5

	MaxEmailLength = 128

	// File contraints
	MaxProfilePicSize = 2 * 1024 * 1024

	// Game constraints
	JoinCodeLen = 7

	GameNameMin = 4
	GameNameMax = 50

	MinPlayerAge = 13
	MaxPlayerAge = 120

	MinLatitude = -90.00
	MaxLatitude = 90.00

	MinLongitude = -180.00
	MaxLongitude = 180.00

	MinGameRating = 1.0
	MaxGameRating = 5.0
)
