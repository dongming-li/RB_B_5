package model

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"

	"gopkg.in/mgo.v2/bson"
)

// Game represents a game played by [User]s
// A game is also a temporary team
type Game struct {
	ID          bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string          `json:"name,omitempty" bson:"name,omitempty"`
	StartTime   time.Time       `json:"startTime,omitempty" bson:"startTime,omitempty"`
	Location    Location        `json:"location,omitempty" bson:"location,omitempty"`
	Host        bson.ObjectId   `json:"host,omitempty" bson:"host,omitempty"`
	Members     []bson.ObjectId `json:"members,omitempty" bson:"members,omitempty"`
	MemsInvited []bson.ObjectId `json:"memsInvited,omitempty" bson:"memsInvited,omitempty"`
	EndTime     time.Time       `json:"endTime,omitempty" bson:"endTime,omitempty"`
	Sport       sport           `json:"sport,omitempty" bson:"sport,omitempty"`
	Rating      Rating          `json:"rating,omitempty" bson:"rating,omitempty"`
	AgeRange    [2]int8         `json:"agerange,omitempty" bson:"agerange,omitempty"`
	JoinCode    string          `json:"joincode,omitempty" bson:"joincode,omitempty"`
	Private     bool            `json:"private,omitempty" bson:"private,omitempty"`
}

// Location represents the geopgrahic address as a place through
// longitude and latitide
type Location struct {
	Lng float64 `json:"lng,omitempty" bson:"lng,omitempty"`
	Lat float64 `json:"lat,omitempty" bson:"lat,omitempty"`
}

// sport is an integer representation of each supported sport
type sport float64

// supported sports
const (
	soccer sport = iota
	basketball
	volleyball
	baseball
	frisbee
	discgolf
)

var sports = []sport{
	soccer,
	basketball,
	volleyball,
	baseball,
	frisbee,
	discgolf,
}

// Rating represents the a pair of rating to count
// count is the number of raters
// [0] -> rating
// [1] -> count
type Rating [2]float64

// newGame represents a data structure to hold new game data before it is processed
type newGame struct {
	name      interface{} //string
	startTime interface{} //time.Time
	endTime   interface{} //time.Time
	host      interface{} //string
	sport     interface{} //sport
	minAge    interface{} //int8
	maxAge    interface{} //int8
	lat       interface{} //float64
	lng       interface{} //float64
	private   interface{} //bool
}

type iD string

type joincode string

func (id iD) Identify() map[string]interface{} {
	return map[string]interface{}{"_id": bson.ObjectIdHex(string(id))}
}

func (j joincode) Identify() map[string]interface{} {
	return map[string]interface{}{"joincode": string(j)}
}

// Get returns the rating value
func (r *Rating) Get() float64 {
	return r[0]
}

// Add updates a rating by including an additional rating to the rolling average
func (r *Rating) add(ra float64) {
	total := r[0]*r[1] + ra
	r[1]++
	r[0] = total / r[1]
}

// getGame returns the game based on the identifier
func getGame(c Collection, identifier Identity) (*Game, error) {
	game := new(Game)
	query := c.Find(identifier.Identify())

	err := query.One(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s sport) OK() bool {
	if s < sport(0) || s > sport(len(sports)-1) {
		return false
	}

	return true
}

func timeRangeOK(start, end interface{}) (string, string, bool) {
	startStr, ok := start.(string)
	endStr, okay := end.(string)
	if !ok || !okay {
		return "", "", false
	}
	startT, err := time.Parse(time.RFC3339, startStr)
	endT, erro := time.Parse(time.RFC3339, endStr)
	if err != nil || erro != nil || !validation.IsValidTimeRange(startT, endT) {
		return "", "", false
	}
	return startStr, endStr, true
}

func (g *newGame) OK() error {
	if n, ok := g.name.(string); !ok || !validation.IsValidGameName(n) {
		return yoda.ClientError{Code: yerr.InvalidGameName, Message: "Invalid Game Name"}
	}

	if _, _, ok := timeRangeOK(g.startTime, g.endTime); !ok {
		return yoda.ClientError{Code: yerr.InvalidTimeRange, Message: "Invalid Time Range"}
	}
	if h, ok := g.host.(string); !ok || !validation.IsValidUserName(h) {
		return yoda.ClientError{Code: yerr.InvalidHost, Message: "Invalid Host"}
	}
	if i, ok := g.sport.(float64); !(ok && sport(i).OK()) {
		return yoda.ClientError{Code: yerr.InvalidSport, Message: "Invalid Sport"}
	}
	minAge, ok := g.minAge.(float64)
	maxAge, okay := g.maxAge.(float64)
	if !ok || !okay || !validation.IsValidAgeRange(minAge, maxAge) {
		return yoda.ClientError{Code: yerr.InvalidAgeRange, Message: "Invalid age range"}
	}
	lat, ok := g.lat.(float64)
	lng, okay := g.lng.(float64)
	if !ok || !okay || !validation.IsValidLocation(lat, lng) {
		return yoda.ClientError{Code: yerr.InvalidLocation, Message: "Invalid location"}
	}
	if _, ok := g.private.(bool); !ok {
		return yoda.ClientError{Code: yerr.InvalidPrivateGame, Message: "Invalid private game value"}
	}
	return nil
}

//createGameEntry assigns data to a variable in a game
func createGameEntry(data map[string]interface{}) (*newGame, error) {
	game := &newGame{
		name:      data["name"],
		startTime: data["startTime"],
		lat:       data["lat"],
		lng:       data["lng"],
		host:      data["host"],
		endTime:   data["endTime"],
		sport:     data["sport"],
		minAge:    data["minAge"],
		maxAge:    data["maxAge"],
		private:   data["private"],
	}
	// if game.endTime == nil {
	//make game end time be two hours later than start time
	// }
	err := game.OK()
	if err != nil {
		return nil, err //TODO: put OK on pointer
	}

	return game, nil
}

func (g newGame) createGame(c, prevGameC Collection) (*Game, error) {
	var src = rand.NewSource(time.Now().UnixNano())
	g.startTime, _ = time.Parse(time.RFC3339, g.startTime.(string))
	g.endTime, _ = time.Parse(time.RFC3339, g.endTime.(string))

	game := &Game{
		ID:        bson.NewObjectId(),
		Name:      g.name.(string),
		StartTime: g.startTime.(time.Time),
		Location:  Location{Lat: g.lat.(float64), Lng: g.lng.(float64)},
		Host:      g.host.(bson.ObjectId),
		Members:   []bson.ObjectId{},
		EndTime:   g.endTime.(time.Time),
		Sport:     sport(g.sport.(float64)),
		Rating:    Rating{0, 0},
		AgeRange:  [2]int8{int8(g.minAge.(float64)), int8(g.maxAge.(float64))},
		JoinCode:  generateJoinCode(constraints.JoinCodeLen, src),
		Private:   g.private.(bool),
	}

	if err := c.Insert(game); err != nil {
		return nil, err
	}
	if err := prevGameC.Insert(game); err != nil {
		return nil, err
	}
	return game, nil
}

//CreateGame attempts to create a game and returns an error if unsuccessful
func CreateGame(store Store, data map[string]interface{}) (string, error) {
	//validate game
	g, err := createGameEntry(data)
	if err != nil {
		return "", err
	}

	// present, _ = store.UserIsPresent(g.host.(string))
	// if !present {
	// 	return "", err
	// }
	userC := store.GetCollection("user")
	h, err := GetUserByUsername(userC, g.host.(string))
	if err != nil {
		return "", err
	}
	if h.CurrentGame != "" {
		existingG, _ := GetGameFromID(store.GetCollection("game"), h.CurrentGame.Hex())
		if existingG != nil {
			if !existingG.hasEnded() {
				return "", yoda.ClientError{Code: yerr.AlreadyInGame, Message: "User already in a game"}
			}
			// If there was an existing game and it wasn't live, it means an expired game was still
			// in the game collection
			log.Println("MONGOD did not remove this game when it was supposed to")
		}
	}

	//TODO make sure no game the user is in has time during this game
	g.host = h.ID

	c := store.GetCollection("game")
	game, err := g.createGame(c, store.GetCollection("prevgame"))

	// give the host the current game
	// There should be no error from this because the user was just gotten from the db
	// Worst case there's a race condition where the user just got deleted, in that case,
	// undo the created game
	h.startGame(userC, game)

	return game.JoinCode, err
}

//JoinGame will have a user ID be added to the game list and return the game information
func (g *Game) JoinGame(store Store, username string) error {
	userC := store.GetUsers()
	u, err := GetUserByUsername(userC, username)
	if err != nil {
		return err
	}

	if u.CurrentGame != "" {
		return yoda.ClientError{Code: yerr.AlreadyInGame, Message: "User Already In a Game"}
	}
	prevgameC := store.GetCollection("prevgame")
	gameC := store.GetCollection("game")
	update := bson.M{"$addToSet": bson.M{"members": u.ID}}

	_, err = gameC.FindAndModify(bson.M{"_id": g.ID}, update, &g, true)
	if err != nil {
		return err
	}
	_, err = prevgameC.FindAndModify(bson.M{"_id": g.ID}, update, nil, false)
	if err != nil {
		return err
	}

	return u.startGame(userC, g)
}

//LeaveGame will have a user ID be removed from the game list and set the users current game to ""
func (g *Game) LeaveGame(store Store, u *User) error {
	c := store.GetUsers()
	halftime := g.EndTime.Sub(g.StartTime) / time.Duration(2)
	if time.Now().Before(g.StartTime.Add(halftime)) {
		gameC := store.GetCollection("game")
		prevgameC := store.GetCollection("prevgame")
		err := gameC.Update(bson.M{"_id": g.ID}, bson.M{"$pull": bson.M{"members": u.ID}})
		if err != nil {
			return err
		}
		err = prevgameC.Update(bson.M{"_id": g.ID}, bson.M{"$pull": bson.M{"members": u.ID}})
		if err != nil {
			return err
		}
		return u.endGame(c, false)
	}

	return u.endGame(c, true)
}

//HostLeaveGame will replace host and remove new host from the game list and set the old hosts current game to ""
func (g *Game) HostLeaveGame(store Store, host *User) error {
	gameC := store.GetCollection("game")
	prevgameC := store.GetCollection("prevgame")
	c := store.GetUsers()

	if len(g.Members) == 0 {
		err := endGame(gameC, g)
		if err != nil {
			return err
		}
		return host.endGame(c, false)
	}

	g.Host = g.Members[0]

	halftime := g.EndTime.Sub(g.StartTime) / time.Duration(2)
	err := gameC.Update(bson.M{"_id": g.ID}, bson.M{"$pull": bson.M{"members": g.Host}, "$set": bson.M{"host": g.Host}})
	if err != nil {
		return err
	}
	if time.Now().Before(g.StartTime.Add(halftime)) {
		err = prevgameC.Update(bson.M{"_id": g.ID}, bson.M{"$pull": bson.M{"members": g.Host}, "$set": bson.M{"host": g.Host}})
		if err != nil {
			return err
		}
		return host.endGame(c, false)
	}
	return host.endGame(c, true)
}

// hasEnded returns true if a game has ended
func (g *Game) hasEnded() bool {
	return time.Now().After(g.EndTime)
}

//endGame will take the game and remove the it from the collection of games
func endGame(c Collection, g *Game) error {
	err := c.Remove(bson.M{"_id": g.ID})
	return err
}

//ClearGame will clear all of the user's in the game current game to ""
func ClearGame(s Store, g string) error {
	uc, prevgameC := s.GetUsers(), s.GetCollection("prevgame")
	var game *Game
	var err error

	if !bson.IsObjectIdHex(g) {
		return yoda.ErrMissingData("code", "string")
	}

	game, err = GetGameFromID(prevgameC, g)
	if err != nil {
		return err
	}

	// remove game from the host's current game
	h, err := game.GetHost(uc)
	if err != nil {
		return err
	}
	h.endGame(uc, true)

	// remove game from the members' current game
	users, err := game.GetMembers(uc)
	if err != nil {
		return err
	}

	var errOccurred error
	for _, user := range users {

		// if an error occurs when ending a user's games,
		// keep doing it but save the last error so the
		// client is notified and also for debugging
		err = user.endGame(uc, true)
		if err != nil {
			errOccurred = err
		}
	}
	if errOccurred != nil {
		return errOccurred
	}
	// remove game from the members' current game
	invitedMems, err := game.GetInvitedMems(uc)
	if err != nil {
		return err
	}
	for _, user := range invitedMems {

		// if an error occurs when removing a user's game request,
		// keep doing it but save the last error so the
		// client is notified and also for debugging
		err = user.RemoveGameRequest(uc, prevgameC, g)
		if err != nil {
			errOccurred = err
		}
	}
	return errOccurred
}

// GetMembers returns the members of the game as [model.User]s
func (g *Game) GetMembers(uc Collection) ([]*User, error) {
	return g.getMembersList(uc, g.Members)
}

// GetInvitedMems returns the invited members of the game as [model.User]s
func (g *Game) GetInvitedMems(uc Collection) ([]*User, error) {
	return g.getMembersList(uc, g.MemsInvited)
}

func (g *Game) getMembersList(uc Collection, memIds []bson.ObjectId) ([]*User, error) {
	users := make([]*User, len(memIds))
	err := uc.Find(bson.M{"_id": bson.M{"$in": memIds}}).All(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetHost returns the host of the game as a [model.User]
func (g *Game) GetHost(uc Collection) (*User, error) {
	var h *User
	err := uc.Find(bson.M{"_id": g.Host}).One(&h)
	if err != nil {
		return nil, err
	}

	return h, nil
}

// UpdateRating updates the rating of a game
func (g *Game) UpdateRating(c Collection, r float64) error {
	g.Rating.add(r)
	return c.Update(bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"rating": g.Rating}})
}

// HasMember returns true if a user with id [id] is a member
// or the host of the game
func (g *Game) HasMember(id bson.ObjectId) bool {
	if g.Host == id {
		return true
	}

	for _, m := range g.Members {
		if m == id {
			return true
		}
	}

	return false
}

// generateJoinCode generates a join code for a game
func generateJoinCode(n int, source rand.Source) string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const (
		letterBits  = 6                 // 6 bits to represent a letter index
		letterBMask = 1<<letterBits - 1 // All 1-bits, as many as letterIdxBits
		letterBMax  = 63 / letterBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterBMax characters!
	for i, cache, remain := n-1, source.Int63(), letterBMax; i >= 0; {
		if remain == 0 {
			cache, remain = source.Int63(), letterBMax
		}
		if idx := int(cache & letterBMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterBits
		remain--
	}
	//TODO make sure game code hasn't been used before
	return string(b)
}

//GetGamesNearLocation returns games near the given location
func GetGamesNearLocation(c Collection, l Location, radius float64) ([]*Game, error) {
	games := make([]*Game, 10)

	// TODO: use buffer.Write to optimize this
	query := fmt.Sprintf("this.location.lng > %f && this.location.lng < %f && this.location.lat > %f && this.location.lat < %f",
		l.Lng-radius, l.Lng+radius, l.Lat-radius, l.Lat+radius)

	err := c.Find(bson.M{"$where": bson.JavaScript{Code: query}}).All(&games)
	return games, err
}

//GetGameFromID returns game a game with the given ID or join code
func GetGameFromID(c Collection, code string) (*Game, error) {
	return getGame(c, iD(code))
}

//GetGameFromJoinCode returns game a game with the given ID or join code
func GetGameFromJoinCode(c Collection, code string) (*Game, error) {
	return getGame(c, joincode(code))
}

//ChangeGame will change the game
func (g *Game) ChangeGame(gameC Collection, info map[string]interface{}) error {
	errC := make(chan error)
	updates := 0
	if name, ok := info["name"]; ok {
		updates++
		go g.changeName(name, errC)
	}

	minAge, ok := info["minAge"]
	if !ok {
		minAge = float64(g.AgeRange[0])
	}
	maxAge, okay := info["maxAge"]
	if !okay {
		maxAge = float64(g.AgeRange[1])
	}
	if okay || ok {
		updates++
		go g.changeAgeRange(minAge, maxAge, errC)
	}
	lat, ok := info["lat"]
	lng, okay := info["lng"]
	if ok && okay {
		updates++
		go g.changeLocation(lat, lng, errC)
	}

	if sport, ok := info["sport"]; ok {
		updates++
		go g.changeSport(sport, errC)
	}

	if private, ok := info["private"].(bool); ok {
		g.Private = private
	}

	startTime, ok := info["startTime"]
	if !ok {
		startTime = g.StartTime.String()
	}
	endTime, okay := info["endTime"]
	if !okay {
		endTime = g.EndTime.String()
	}
	if okay || ok {
		updates++
		go g.changeTime(startTime, endTime, errC)
	}
	for i := 0; i < updates; i++ {
		err := <-errC
		if err != nil {
			return err
		}
	}
	return gameC.Update(bson.M{"_id": g.ID}, bson.M{"$set": bson.M{"name": g.Name, "agerange": g.AgeRange, "location": g.Location, "sport": g.Sport, "startTime": g.StartTime, "endTime": g.EndTime, "private": g.Private}})
}

//changeName changes the name of the game
func (g *Game) changeName(n interface{}, errC chan<- error) {
	name, ok := n.(string)
	if !ok || !validation.IsValidGameName(name) {
		errC <- yoda.ClientError{Code: yerr.InvalidGameName, Message: "Invalid Game Name"}
		return
	}
	g.Name = name
	errC <- nil
}

//changeTime changes the TimeRange of the game
func (g *Game) changeTime(start, end interface{}, errC chan<- error) {
	startStr, endStr, ok := timeRangeOK(start, end)
	if !ok {
		errC <- yoda.ClientError{Code: yerr.InvalidTimeRange, Message: "Invalid Time Range"}
		return
	}
	g.StartTime, _ = time.Parse(time.RFC3339, startStr)
	g.EndTime, _ = time.Parse(time.RFC3339, endStr)
	errC <- nil
}

//changeAgeRange changes the agerange of the game
func (g *Game) changeAgeRange(startAge, endAge interface{}, errC chan<- error) {
	minAge, ok := startAge.(float64)
	maxAge, okay := endAge.(float64)
	if !ok || !okay || !validation.IsValidAgeRange(minAge, maxAge) {
		errC <- yoda.ClientError{Code: yerr.InvalidAgeRange, Message: "Invalid age range"}
		return
	}
	g.AgeRange = [2]int8{int8(minAge), int8(maxAge)}
	errC <- nil
}

//changeLocation changes the Location of the game
func (g *Game) changeLocation(lat, lng interface{}, errC chan<- error) {
	latitude, ok := lat.(float64)
	long, okay := lng.(float64)
	if !ok || !okay || !validation.IsValidLocation(latitude, long) {
		errC <- yoda.ClientError{Code: yerr.InvalidLocation, Message: "Invalid location"}
		return
	}
	g.Location = Location{Lat: latitude, Lng: long}
	errC <- nil
}

//changeSport changes the sport of the game
func (g *Game) changeSport(s interface{}, errC chan<- error) {
	sp, ok := s.(float64)
	if !(ok && sport(sp).OK()) {
		errC <- yoda.ClientError{Code: yerr.InvalidSport, Message: "Invalid Sport"}
		return
	}
	g.Sport = sport(sp)
	errC <- nil
}
