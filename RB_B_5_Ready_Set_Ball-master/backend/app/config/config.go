package config

// Server-wide config constants
const (
	//session vars
	SessionName      = "session_var"
	PostSessionName  = "todo_change_this"
	SessionAuth      = "Super-secret-from-environment-var"
	SessionStoreName = "session"

	//environment settings
	Dev  = true
	Beta = false

	//database
	BetaDatabase  = "beta"
	LocalDatabase = "temp"
	MongoDBHosts  = "ds044709.mlab.com:44709"
	AuthDatabase  = "beta"
	AuthUserName  = "ginger"
	AuthPassword  = "ginger"

	//server constraints
	MaxSessionAttempts = 10
)

// DBNameMain is the name of the database to be connected to depending on whether [Beta] is true
var DBNameMain = map[bool]string{true: BetaDatabase, false: LocalDatabase}[Beta]
