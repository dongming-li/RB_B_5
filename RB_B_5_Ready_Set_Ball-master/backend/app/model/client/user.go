package client

import (
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"gopkg.in/mgo.v2/bson"
)

// Model represents an object type sent back to the user
type Model interface{}

// User represents a reduced [model.User] sent to the client
type User struct {
	Username       string                 `json:"username,omitempty"`
	Firstname      string                 `json:"firstname,omitempty"`
	Lastname       string                 `json:"lastname,omitempty"`
	PersonalRating float32                `json:"pRating,omitempty"`
	TeamRating     float64                `json:"rating"`
	City           string                 `json:"city,omitempty"`
	Friends        []User                 `json:"friends"`
	FriendRequests []User                 `json:"friendRequests"`
	FriendStatus   FriendStatus           `json:"friendStatus"`
	GameRequests   []model.GameRequest    `json:"gameRequests"`
	Avatar         []byte                 `json:"avatar,omitempty" bson:"avatar,omitempty"`
	PrevGames      []Game                 `json:"prevGames,omitempty" bson:"prevGames"`
	Preferences    map[string]interface{} `json:"preferences,omitempty"`
}

// FriendStatus represents the friendship status that the logged in user has with another user
type FriendStatus string

const (
	areFriends     FriendStatus = "areFriends"
	sentR                       = "sentRequest"
	receivedR                   = "receivedRequest"
	isUser                      = "isUser"
	noFriendStatus              = ""
)

type sessionUser struct {
	id, username string
}

// newUser sets up a client appropriate [model.User]
func newUser(user *model.User) *User {
	return &User{
		Username:       user.Username,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		PersonalRating: user.PersonalRating,
		TeamRating:     user.TeamRating,
		GameRequests:   user.GameRequests,
		Avatar:         user.Avatar,
	}
}

func (u *User) populatePrivateInfo(c model.Collection, user *model.User) error {
	u.City = user.City
	u.GameRequests = user.GameRequests
	friends, err := user.GetFriendRequests(c)
	if err != nil {
		return err
	}

	u.FriendRequests = make([]User, len(friends))
	for i, f := range friends {
		u.FriendRequests[i] = *newUser(&f)
	}

	return nil
}

// populate fills the exported properties based on the unexported properties
// by making the required database calls
// e.g. [user.Friends] gets filled from db calls based on [user.friends]
func (u *User) populate(friends []model.User) {
	u.Friends = make([]User, len(friends))

	// Convert each [model.User] friend to [User] and add to the friends list
	for i, friend := range friends {
		cFriend := newUser(&friend)
		u.Friends[i] = *cFriend
	}
}

func (u *User) populatePrevGame(games []model.Game) {
	u.PrevGames = make([]Game, len(games))

	// Convert each [model.Game] in to [Game] and add to the prevGames list
	for i := range games {
		cGame := newGame(&games[i])
		u.PrevGames[i] = *cGame
	}
}

func getUser(s model.Store, username string, owner sessionUser, populate bool) (*User, error) {
	c := s.GetUsers()
	user, err := model.GetUserByUsername(c, username)
	if err != nil {
		return nil, err
	}

	u := newUser(user)
	if pref, ok := user.Meta["client_pref"].(map[string]interface{}); ok {
		u.Preferences = pref
	}
	// add private user info if this is the owner of the profile
	if username == owner.username {
		//TODO handle other errors besides just game not found
		user.UpdateGameRequests(s)
		u.populatePrivateInfo(c, user)
	}

	// add friendStatus with the logged in user
	u.FriendStatus = getFriendStatus(user, owner)

	if populate {
		friends, err := user.GetFriends(c)
		if err != nil {
			return nil, err
		}

		// populate user with friends
		u.populate(friends)

		games, err := user.GetPrevGames(s.GetCollection("prevgame"))
		if err != nil {
			return nil, err
		}

		// populate user with prevgames
		u.populatePrevGame(games)

	}
	return u, err
}

// GetUser returns a client appropriate [User] object
// Note: It returns null for properties referenced by [bson.ObjectId]
func GetUser(s model.Store, username, owner, ownerID string) (*User, error) {
	return getUser(s, username, sessionUser{username: owner, id: ownerID}, false)
}

// GetUserAndPopulate returns a client appropriate [User] object
// This makes an additional db call to get the [User]'s friends
func GetUserAndPopulate(s model.Store, username, owner, ownerID string) (*User, error) {
	return getUser(s, username, sessionUser{username: owner, id: ownerID}, true)
}

// getFriendStatus returns the [FriendStatus] between the [user] and the currently logged in user
func getFriendStatus(user *model.User, owner sessionUser) FriendStatus {
	if user.Username == owner.username {
		return isUser
	}
	for _, req := range user.SentRequests {
		if req == bson.ObjectIdHex(owner.id) {
			return receivedR
		}
	}
	for _, req := range user.FriendRequests {
		if req == bson.ObjectIdHex(owner.id) {
			return sentR
		}
	}
	for _, req := range user.Friends {
		if req == bson.ObjectIdHex(owner.id) {
			return areFriends
		}
	}
	return noFriendStatus
}
