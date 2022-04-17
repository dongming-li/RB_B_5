package model

import (
	"log"

	"strings"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/analytics"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"gopkg.in/mgo.v2/bson"
)

// User represents a User
type User struct {
	ID             bson.ObjectId          `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string                 `json:"username,omitempty" bson:"username,omitempty"`
	Password       string                 `json:"password,omitempty" bson:"password,omitempty"`
	Firstname      string                 `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname       string                 `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email          string                 `json:"email,omitempty" bson:"email,omitempty"`
	City           string                 `json:"city,omitempty" bson:"city,omitempty"`
	PersonalRating float32                `json:"pRating,omitempty" bson:"pRating,omitempty"`
	TeamRating     float64                `json:"rating,omitempty" bson:"rating,omitempty"`
	Friends        []bson.ObjectId        `json:"friends,omitempty" bson:"friends,omitempty"`
	CurrentGame    bson.ObjectId          `json:"currentGame,omitempty" bson:"currentGame,omitempty"`
	FriendRequests []bson.ObjectId        `json:"friendRequests,omitempty" bson:"friendRequests,omitempty"`
	SentRequests   []bson.ObjectId        `json:"sentRequests,omitempty" bson:"sentRequests,omitempty"`
	PrevGames      map[string]float64     `json:"prevGames,omitempty" bson:"prevGames,omitempty"`
	GameRequests   []GameRequest          `json:"gameRequests,omitempty" bson:"gameRequests,omitempty"`
	Avatar         []byte                 `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Meta           map[string]interface{} `json:"meta,omitempty" bson:"meta,omitempty"`
}

// newUser represents a data structure to hold new user data before it is processed
type newUser struct {
	username    string //`json:"username,omitempty"`
	password    string //`json:"password,omitempty"`
	confirmPass string //`json:"confirmPass,omitempty"`
	firstname   string //`json:"firstname,omitempty"`
	lastname    string //`json:"lastname,omitempty"`
	email       string
	city        string
}

// identifier for a [User]'s [username]
type username string

//GameRequest stores the information about a game request
type GameRequest struct {
	From string `json:"from,omitempty" bson:"from,omitempty"`
	Game string `json:"game,omitempty" bson:"game,omitempty"`
}

type populateType int

// Represents the invite populate type
const (
	Friend populateType = iota
	FriendRequest
	SentRequest
	PrevGames
)

func (un username) Identify() map[string]interface{} {
	return map[string]interface{}{"username": string(un)}
}

func (u *newUser) OK() error {
	if !validation.IsValidName(u.firstname) {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid firstname"}
	}
	if !validation.IsValidName(u.lastname) {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid lastname"}
	}
	if !validation.IsValidUserName(u.username) {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username"}
	}
	u.username = strings.ToLower(u.username)
	if !validation.IsValidPassword([]rune(u.password)) {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid password"}
	}
	if !validation.IsValidEmail(u.email) {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid email"}
	}
	u.email = validation.CleanEmail(u.email)
	if u.password != u.confirmPass {
		return yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Passwords do not match"}
	}
	if !validation.IsValidName(u.city) {
		return yoda.ClientError{Code: yerr.InvalidParameter, Message: "Invalid city"}
	}

	return nil
}

// GetUsers returns an array of the users and an error
func GetUsers(store Store, usernames ...string) ([]User, error) {
	var result []User
	c := store.GetUsers()
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// getUser returns a user if the [identifier] is present, or nil and an [error] if an error occurs.
// It uses the collection from [GetUsers].
// It returns an [error] if an invalid identifier is passed in
// We can export this and get rid of [GetUserByUsername] and the likes but the user will have to
// create their own identify before calling this method, for a username example, it'll just be:
// u := username{username: uName}
// user, err := getUser(c, u)
func getUser(c Collection, identifier Identity) (*User, error) {
	user := new(User)
	query := c.Find(identifier.Identify())

	err := query.One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername returns a [*User] if the [username] is present or or nil and an [error] if an error occurs
func GetUserByUsername(c Collection, uName string) (*User, error) {
	if !validation.IsValidUserName(uName) {
		return nil, &yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username/password"}
	}
	u := username(uName)
	return getUser(c, u)
}

// CreateUser creates a user and returns an error if unsuccessful
func CreateUser(store Store, data map[string]string) error {
	//validate user

	nUser, err := makeNewUserEntry(data)
	if err != nil {
		return err
	}

	//check if user exists by using store.Find
	present, err := store.UserIsPresent(nUser.username)
	if err != nil {
		return err
	}
	if present {
		return yoda.ClientError{Code: yerr.UserAlreadyExists, Message: "Invalid Username"}
	}

	//create user struct
	_, err = nUser.createUser(store.GetUsers())

	return err
}

// GetFriends returns a slice of [User]s representing the [u]'s friends
func (u *User) GetFriends(c Collection) ([]User, error) {
	return u.getFriends(c, Friend)
}

// GetFriendRequests returns a slice of [User]s representing the [u]'s friend requests
func (u *User) GetFriendRequests(c Collection) ([]User, error) {
	return u.getFriends(c, FriendRequest)
}

//UpdateGameRequests will update the users gamerequests
func (u *User) UpdateGameRequests(s Store) {
	var err error
	for i, gr := range u.GameRequests {
		_, err = GetGameFromID(s.GetCollection("game"), gr.Game)
		if err != nil {
			// TODO: account for other error? log?
			ClearGame(s, gr.Game)
			u.GameRequests = append(u.GameRequests[:i], u.GameRequests[i+1:]...)
		}
	}
}

// getFriends returns a slice of [User]s representing
// the [u]'s friends, friend requests or sent request
func (u *User) getFriends(c Collection, p populateType) ([]User, error) {
	var friends []User

	// get [model.User] friends of user
	var f []bson.ObjectId
	switch p {
	case Friend:
		f = u.Friends
	case FriendRequest:
		f = u.FriendRequests
	case SentRequest:
		f = u.SentRequests
	default:
		return nil, &yoda.ClientError{Code: yerr.Internal, Message: "Invalid populate type"}
	}

	err := c.Find(bson.M{"_id": bson.M{"$in": f}}).All(&friends)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

// SendFriendRequest updates the user's [sentRequests] with the id of the [to] user
func (u *User) SendFriendRequest(c Collection, to *User) error {
	// tell the to user a request was sent
	err := to.receiveFriendRequest(c, u)
	if err != nil {
		return nil
	}

	//TODO: check the err to see if the update was made or
	//if the req had previously been sent and send an error if that's the case
	// TODO: undo receiveFriendRequest if an error occurs from this
	return c.Update(bson.M{"username": u.Username}, bson.M{"$addToSet": bson.M{"sentRequests": to.ID}})
}

// ReceiveFriendRequest updates the user's [friendRequests] with the id of the [from] user
func (u *User) receiveFriendRequest(c Collection, from *User) error {

	//TODO: check the err to see if the update was made or
	// if the req had previously been sent and send an error if that's the case
	return c.Update(bson.M{"username": u.Username}, bson.M{"$addToSet": bson.M{"friendRequests": from.ID}})
}

// AddFriend removes a user from a user's friendRequests or sentRequets and
// adds that user to the friends
// TODO: this should handle moving the other user frm the other array as well
func (u *User) AddFriend(c Collection, other *User, received bool) error {
	err := u.modifyFriends(c, other, received, true)
	if err != nil {
		return err
	}

	err = other.modifyFriends(c, u, !received, true)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFriendRequest removes the [other] user from the [u]'s friendRequests
// and removes the [u] user from [other]'s sentRequests
//
// TODO: maybe do some logging of rejected requests,
// or add some limit as to how many times a user can send a request to another user and get rejected
func (u *User) DeleteFriendRequest(c Collection, other *User) error {
	err := u.modifyFriends(c, other, true, false)
	if err != nil {
		return err
	}

	err = other.modifyFriends(c, u, false, false)
	if err != nil {
		return err
	}

	return nil
}

// modifyFriends removes a user frim a user's friendRequests or sentRequets
// if [add] is true, it also adds that user to the friends and
// adds the user to the other friend's slice
func (u *User) modifyFriends(c Collection, other *User, received, add bool) error {
	var source string
	if received {
		source = "friendRequests"
	} else {
		source = "sentRequests"
	}

	err := c.Update(bson.M{"username": u.Username}, bson.M{"$pull": bson.M{source: other.ID}})
	if err != nil {
		return err
	}

	// add friend to friend list
	if add {
		err = c.Update(bson.M{"username": u.Username}, bson.M{"$addToSet": bson.M{"friends": other.ID}})
		if err != nil {
			return err
		}
	}

	return nil
}

// AddGameRequest updates the user's [gameRequests] with the id of the game and the user it's from
func (u *User) AddGameRequest(c, gameC, prevgameC Collection, from, g string) error {
	if !bson.IsObjectIdHex(g) {
		return &yoda.ClientError{Code: yerr.InvalidID, Message: "Invalid game ID"}
	}
	game, err := GetGameFromID(gameC, g)
	if err != nil {
		return err
	}

	//from this game request they can then get the user who invited them and the actual game by calling other functions with these items
	gameR := &GameRequest{
		From: from,
		Game: g,
	}
	err = c.Update(bson.M{"username": u.Username}, bson.M{"$addToSet": bson.M{"gameRequests": gameR}})
	if err != nil {
		return err
	}
	return prevgameC.Update(bson.M{"_id": game.ID}, bson.M{"$addToSet": bson.M{"memsInvited": u.ID}})
}

// RemoveGameRequest removes a game request from a user's gameRequests
func (u *User) RemoveGameRequest(c, prevgameC Collection, g string) error {
	if !bson.IsObjectIdHex(g) {
		return &yoda.ClientError{Code: yerr.InvalidID, Message: "Invalid game ID"}
	}
	update := bson.M{"$pull": bson.M{"gameRequests": bson.M{"game": g}}}
	length := len(u.GameRequests)
	_, err := c.FindAndModify(bson.M{"_id": u.ID}, update, &u, true)
	if err != nil {
		return err
	}
	if length == len(u.GameRequests) {
		return &yoda.ClientError{Code: yerr.InvalidRequest, Message: "Invalid game request"}
	}

	game, err := GetGameFromID(prevgameC, g)
	if err != nil {
		return err
	}
	return prevgameC.Update(bson.M{"_id": game.ID}, bson.M{"$pull": bson.M{"memsInvited": u.ID}})

}

//RemoveFriend will remove a users friend
func (u *User) RemoveFriend(c Collection, other *User) error {
	return c.Update(bson.M{"username": u.Username}, bson.M{"$pull": bson.M{"friends": other.ID}})
}

// ChangeAvatar changes the avatar(profile picture) of the user
func (u *User) ChangeAvatar(c Collection, image []byte) error {

	return c.Update(bson.M{"username": u.Username}, bson.M{"$set": bson.M{"avatar": image}})
}

// startGame unsafely updates the current game of the user
// no checks to see if the user is currently in another game
func (u *User) startGame(c Collection, g *Game) error {
	return c.Update(bson.M{"username": u.Username}, bson.M{"$set": bson.M{"currentGame": g.ID}})
}

// endGame removes the user's [CurrentGame]
func (u *User) endGame(c Collection, addPrev bool) error {
	// remove currentGame if user has no current game
	// or if we're not adding to prevGames
	if len(u.CurrentGame.Hex()) == 0 || !addPrev {
		return c.Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"currentGame": ""}})
	}

	if u.PrevGames == nil {
		u.PrevGames = make(map[string]float64, 1)
	}
	u.PrevGames[u.CurrentGame.Hex()] = 0 // TODO: might be an inefficient update if there are many prev games
	return c.Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"prevGames": u.PrevGames, "currentGame": ""}})

}

// GetPrevGames returns a slice of [Game]s representing the [u]'s prevGames
func (u *User) GetPrevGames(c Collection) ([]Game, error) {
	var games []Game
	err := c.Find(bson.M{"_id": bson.M{"$in": mapKeys(u.PrevGames)}}).All(&games)
	if err != nil {
		return nil, err
	}

	return games, nil
}

// Rate updates the rating of a previous game
func (u *User) Rate(c Collection, r float64, g *Game) error {
	if u.PrevGames == nil {
		u.PrevGames = make(map[string]float64)
	}
	u.PrevGames[g.ID.Hex()] = r

	return c.Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"prevGames": u.PrevGames}}) // TODO: might be inefficient
}

// Edit edits the properties of the user
func (u *User) Edit(c Collection, data map[string]string) error {
	defer analytics.Track("Edit", time.Now(), data)

	updates := 0
	errC := make(chan error)

	// username
	if uName, ok := data["username"]; ok {
		updates++
		go u.setUsername(c, uName, errC)
	}

	// email
	if email, ok := data["email"]; ok {
		updates++
		go u.setEmail(email, errC)
	}

	// password
	if pass, ok := data["password"]; ok {
		if newPass, ok := data["newPassword"]; ok {
			updates++
			go u.setPassword(pass, newPass, errC)
		}
	}

	// firstname
	if name, ok := data["firstname"]; ok {
		updates++
		go u.setName(name, true, errC)
	}

	// lastname
	if name, ok := data["lastname"]; ok {
		updates++
		go u.setName(name, false, errC)
	}

	for i := 0; i < updates; i++ {
		err := <-errC
		if err != nil {
			return err
		}
	}

	// persist changes to db
	return c.Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{
		"username":  u.Username,
		"password":  u.Password,
		"email":     u.Email,
		"firstname": u.Firstname,
		"lastname":  u.Lastname,
	}})
}

// setUsername changes the username of the user
func (u *User) setUsername(c Collection, name string, err chan<- error) {
	if !validation.IsValidUserName(name) {
		err <- yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username"}
		return
	}

	if u, _ := getUser(c, username(name)); u != nil {
		err <- yoda.ClientError{Code: yerr.UserAlreadyExists, Message: "username is already taken"}
		return
	}

	u.Username = name
	err <- nil
}

// setName changes the firstname or lastname of the user
//
// if [isFirstName] is true, then the first name is changed
func (u *User) setName(name string, isFirstName bool, err chan<- error) {
	if !validation.IsValidName(name) {
		err <- yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid name"}
		return
	}

	if isFirstName {
		u.Firstname = name
	} else {
		u.Lastname = name
	}

	err <- nil
}

// setEmail changes the email of the user
func (u *User) setEmail(email string, err chan<- error) {
	if !validation.IsValidEmail(email) {
		err <- yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid email"}
		return
	}

	u.Email = validation.CleanEmail(email)
	err <- nil
}

// setPassword changes the password of the user
func (u *User) setPassword(password, newPassword string, errC chan<- error) {
	// check the old password matches the current password
	if !validation.IsValidPassword([]rune(password)) || !validation.CheckPassword(u.Password, password) {
		errC <- yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "The password doesn't match our records"}
		return
	}

	// check that new password is valid
	if !validation.IsValidPassword([]rune(newPassword)) {
		errC <- yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "The new password is invalid"}
		return
	}

	//TODO: check that there's a limit to how many attempts can be made to failed password attempts

	// hash new password and set new password
	pass, err := validation.HashPassword(newPassword)
	if err != nil {
		errC <- err
	}
	u.Password = pass

	errC <- nil
}

func makeNewUserEntry(data map[string]string) (*newUser, error) {
	user := newUser{
		username:    data["username"],
		firstname:   data["firstname"],
		lastname:    data["lastname"],
		password:    data["password"],
		confirmPass: data["confirmPass"],
		email:       data["email"],
		city:        data["city"],
	}

	err := user.OK()
	if err != nil {
		return nil, err //TODO: put OK on pointer
	}
	return &user, nil //TODO: change user type to pointer
}

func (u newUser) createUser(c Collection) (*User, error) {
	pass, err := validation.HashPassword(u.password)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:             bson.NewObjectId(), //TODO: this might be a bug
		Username:       u.username,
		Email:          u.email,
		Firstname:      u.firstname,
		Lastname:       u.lastname,
		Password:       pass,
		City:           u.city,
		Friends:        []bson.ObjectId{},
		CurrentGame:    "",
		FriendRequests: []bson.ObjectId{},
		SentRequests:   []bson.ObjectId{},
		PrevGames:      map[string]float64{},

		Meta: map[string]interface{}{
			"email_verified_sent": time.Now(),                                        // the time at which the verification email was sent i.e. time of sign up
			"dictionary":          [constraints.MaxPasswordDictionaryLength]string{}, // hashes of previous password for making sure the user doesn't use an older password. (prolly max at 3)
			"locked":              nil,                                               // the time at which a users account was locked
			"attempts":            0,                                                 // number of failed login attempts before locked, or, and number of login attempts after locked
			"client_pref":         defaultPref,                                       // the user's client settings
		},
	}

	if err := c.Insert(user); err != nil {
		return nil, err
	}
	return user, nil
}

var defaultPref = map[string]interface{}{
	"teams_per_page": 10,
}

func (u *User) recalculateRating(c, gameC Collection) error {
	games, err := u.GetPrevGames(gameC)
	if err != nil {
		return err
	}

	var rating float64
	var count int
	for i := range games {
		r := games[i].Rating.Get()
		if r != 0 {
			rating += r
			count++
		}
	}

	if count != 0 {
		rating = rating / float64(count)
		return c.Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"rating": rating}})
	}
	return nil
}

// UpdateRatings updates the ratings of every player in a game
func UpdateRatings(s Store, g *Game) error {
	c, gameC := s.GetUsers(), s.GetCollection("prevgame")

	var members []User
	var allMembers = make([]bson.ObjectId, 0, len(g.Members)+1) //+1 for host
	// add the host
	allMembers = append(allMembers, g.Host)
	// add the members at the end of the game
	allMembers = append(allMembers, g.Members...)

	// No check for error because there will be at least the host user found
	c.Find(bson.M{"_id": bson.M{"$in": allMembers}}).All(&members)

	for i := range members {
		err := (&members[i]).recalculateRating(c, gameC)
		if err != nil {
			// This should be impossible since only updates are happening
			// If this ever happens, all the prev ratings will have to be rolled back
			return err
		}
	}

	return nil
}

// AreFriends returns true if two users are friends
func AreFriends(u1, u2 *User) bool {
	uWithLessFriends, uWithMoreFriends := u1, u2
	if len(u2.Friends) < len(u1.Friends) {
		uWithLessFriends, uWithMoreFriends = u2, u1
	}

	friends := uWithLessFriends.Friends
	for i := range friends {
		if friends[i] == uWithMoreFriends.ID {
			return true
		}
	}

	return false
}
