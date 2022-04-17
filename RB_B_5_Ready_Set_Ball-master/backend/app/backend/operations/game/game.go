package game

import (
	"strconv"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/socket_handler"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model/client"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"

	tr "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/transaction"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"gopkg.in/mgo.v2/bson"
)

type result struct {
	meta map[string]string
	data interface{}
}

func (r result) GetMeta() map[string]string {
	return r.meta
}

func (r result) GetData() interface{} {
	return r.data
}

// CreateGame attempts to create a game hosted by [host] and returns a [transaction.Result] and an [error]
func CreateGame(store model.Store, rData interface{}, host string) (tr.Result, error) {
	data, _ := rData.(map[string]interface{})
	data["host"] = host
	joincode, err := model.CreateGame(store, data)
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{"error": "false"}, data: joincode}, nil //TODO: meta
}

// GetGamesNearLocation returns a list of games around a location
// Currently in a 10x10 radius of the location
func GetGamesNearLocation(s model.Store, data map[string]string, currUser string) ([]*client.Game, error) {
	lng, err := strconv.ParseFloat(data["lng"], 64)
	if err != nil {
		return nil, yoda.ErrMissingData("lng", "float")
	}

	lat, err := strconv.ParseFloat(data["lat"], 64)
	if err != nil {
		return nil, yoda.ErrMissingData("lat", "float")
	}

	const radius float64 = 10
	games, err := model.GetGamesNearLocation(s.GetCollection("game"), model.Location{Lng: lng, Lat: lat}, radius)
	if err != nil {
		return nil, err
	}

	return client.NewGames(s.GetUsers(), games, currUser)
}

// GetGame returns a game found from the id or join code
func GetGame(s model.Store, data map[string]string) (*client.Game, error) {
	g, err := getGame(s, data, false)
	if err != nil {
		return nil, err
	}

	game, err := client.NewGame(s.GetUsers(), g, data["username"])
	if err != nil {
		return nil, err
	}
	return game, nil
}

// getGame makes a call to [fetchGame]
// if the game is found, it returns it if
// it is not private or if it is private but the requester is
// a member or a host, else it returns a NotFound error
// if [force] is true, the game is returned even if it is private and the requester
// is not a member
func getGame(s model.Store, data map[string]string, force bool) (*model.Game, error) {
	g, err := fetchGame(s, data)
	if err != nil {
		return nil, err
	}

	// if the game is private, only members can get it
	if !force && g.Private && !g.HasMember(bson.ObjectIdHex(data["userid"])) {
		return nil, yoda.ErrGameNotFound
	}

	// remove the join code if the game is private
	// and the requester is not the host
	if g.Private && g.Host.Hex() != data["userid"] {
		g.JoinCode = ""
	}

	return g, nil
}

// fetchGame returns a game found from the id or join code
// or returns the users current game if the user has one
func fetchGame(s model.Store, data map[string]string) (*model.Game, error) {
	c := s.GetCollection("game")

	id, _ := data["code"]
	if id == "" {
		u, err := model.GetUserByUsername(s.GetUsers(), data["username"])
		if err != nil {
			return nil, err
		}

		if u.CurrentGame == "" {
			return nil, yoda.ErrGameNotFound
		}
		game, err := model.GetGameFromID(c, u.CurrentGame.Hex())
		if err != nil {
			// TODO: account for other error? log?
			model.ClearGame(s, u.CurrentGame.Hex())
			return nil, yoda.ErrGameNotFound
		}
		return game, err
	}

	if ok := bson.IsObjectIdHex(id); ok {
		game, err := model.GetGameFromID(c, id)
		if err != nil {
			// TODO: account for other error? log?
			model.ClearGame(s, id)
			return nil, yoda.ErrGameNotFound
		}
		return game, err
	}

	if validation.IsValidJoinCode(id) {
		g, err := model.GetGameFromJoinCode(c, id)
		if err != nil {
			// TODO: account for other error? log?
			model.ClearGame(s, id)
			return nil, yoda.ErrGameNotFound
		}

		// if you have the joincode, then the game is as good as
		// public to you
		g.Private = false
		return g, nil
	}

	return nil, yoda.ErrInvalidGame
}

//JoinGame joins a game and then returns a game found from the join code or id
func JoinGame(s model.Store, data map[string]string, fromInvite bool, sm *sockets.Manager) (*client.Game, error) {
	g, err := getGame(s, data, fromInvite)
	if err != nil {
		return nil, err
	}

	username := data["username"]
	err = g.JoinGame(s, username)
	if err != nil {
		return nil, err
	}

	game, err := client.NewGame(s.GetUsers(), g, username)
	if err != nil {
		return nil, err
	}

	members := model.IDsToHex(g.Members, 1)
	members[len(members)-1] = g.Host.Hex()

	go sm.SendMessage(map[string]interface{}{
		"name":   sockethandler.NewGameMember,
		"member": game.Members[len(game.Members)-1],
	}, members)
	return game, nil
}

// Rate checks if the user participated in a game and rates the
// game if the user did
func Rate(s model.Store, data map[string]interface{}) error {
	id := data["code"].(string)
	if !bson.IsObjectIdHex(id) {
		return yoda.ErrMissingData("code", "string")
	}

	rating, ok := data["rating"].(float64)
	if !ok || !validation.IsValidGameRating(rating) {
		return &yoda.ClientError{Code: yerr.InvalidRating, Message: "the game rating was invalid"}
	}

	c := s.GetUsers()
	user, err := model.GetUserByUsername(c, data["username"].(string))
	if err != nil {
		return err
	}

	// check that user was in that game
	var r float64
	if r, ok = user.PrevGames[bson.ObjectIdHex(id).Hex()]; !ok {
		return yoda.ErrInvalidGame
	}

	// check that the user hasn't rated the game
	if r != 0 {
		return &yoda.ClientError{Code: yerr.InvalidRating, Message: "you have already rated the game"}
	}

	// rate the game
	g, err := rate(s.GetCollection("prevgame"), rating, id)
	if err != nil {
		return err
	}

	// update the user's rating of the game
	err = user.Rate(c, rating, g)
	if err != nil {
		//TODO: maybe rollback prevgame
		return err
	}

	// update the ratings of every member of the game including the host
	err = model.UpdateRatings(s, g)
	if err != nil {
		//TODO: maybe rollback prevgame
		return err
	}

	return nil
}

// rate updates the aggregate rating of a game with id, [id] by including the [r] rating
func rate(c model.Collection, r float64, id string) (*model.Game, error) {
	g, err := model.GetGameFromID(c, id)
	if err != nil {
		return nil, err
	}

	err = g.UpdateRating(c, r)
	if err != nil {
		return nil, err
	}

	return g, nil
}

//ExitGame has a user exit a game
func ExitGame(store model.Store, user tr.Owner, sm *sockets.Manager) error {
	userC := store.GetUsers()
	u, err := model.GetUserByUsername(userC, user.Username())
	if err != nil {
		return err
	}
	data := map[string]string{"username": user.Username(), "userid": user.ID()}
	g, err := getGame(store, data, true)
	if err != nil {
		return err
	}

	members := model.IDsToHex(g.Members, 1)
	members[len(members)-1] = g.Host.Hex()

	if u.ID == g.Host {
		err = g.HostLeaveGame(store, u)
		if err != nil {
			return err
		}

		go sm.SendMessage(map[string]interface{}{
			"name": sockethandler.GameHostLeave,
		}, members)
		return nil
	}

	err = g.LeaveGame(store, u)
	if err != nil {
		return err
	}

	go sm.SendMessage(map[string]interface{}{
		"name":     sockethandler.GameMemberLeave,
		"username": u.Username,
	}, members)
	return nil
}

//EditGame will find what information needs to be edited and call the correct function
func EditGame(s model.Store, host string, info map[string]interface{}) (tr.Result, error) {
	res := new(result)
	user, _ := model.GetUserByUsername(s.GetUsers(), host) //No error because we got the username from session
	gameC := s.GetCollection("game")
	g, err := model.GetGameFromID(gameC, user.CurrentGame.Hex())
	if err != nil {
		return nil, err
	}
	if user.ID != g.Host {
		return nil, yoda.ClientError{Code: yerr.InvalidHost, Message: "Invalid Host"}
	}
	if time.Now().After(g.StartTime) {
		return nil, yoda.ClientError{Code: yerr.GameHasStarted, Message: "Game has started"}
	}
	err = g.ChangeGame(gameC, info)
	if err != nil {
		return nil, err
	}
	res.data = map[string]interface{}{"successful": true}
	return res, nil
}

// SocketSubscibe is a handler for incoming socket messages pertaining to game
func SocketSubscibe(from, to string, message map[string]interface{}) ([]string, interface{}) {
	// TODO: this is here for example purposes
	// and should be removed after the first real use of
	// a socket subscription
	return []string{from}, "Acknowledge"
}
