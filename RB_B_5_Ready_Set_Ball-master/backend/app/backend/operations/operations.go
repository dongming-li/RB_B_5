package operations

import (
	"log"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/account"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/upload"

	"fmt"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/authentication"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/game"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/invites"
	tr "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/transaction"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model/client"
)

// Do attempts to perform a server operation and returns a result and an error
func Do(op string, store model.Store, req yoda.Request) (tr.Result, error) {
	var res tr.Result
	var err error
	switch op {
	case "create":
		res, err = createOp(store, req)
	case "get_users":
		res, err = getUsers(store, req)
	case "get_games_from_location":
		res, err = getGamesNearLocation(store, req)
	case "get_game":
		res, err = getGame(store, req)
	case "join_game":
		res, err = joinGame(store, req)
	case "rate_game":
		res, err = rateGame(store, req)
	case "exit_game":
		res, err = exitGame(store, req)
	case "login":
		res, err = login(store, req)
	case "get_user":
		res, err = getUser(store, req)
	case "invite":
		res, err = invite(store, req)
	case "friend_user":
		res, err = friendUserOp(store, req)
	case "remove":
		res, err = remove(store, req)
	case "upload":
		res, err = uploadFile(store, req)
	case "edit":
		res, err = edit(store, req)
	default:
		msg := fmt.Sprintf("Unknown Operation: %s\n", op)
		res, err = nil, yoda.ClientError{Message: msg, Code: yerr.InvalidOperation}
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return res, nil
}

func getUsers(store model.Store, req yoda.Request) (tr.Result, error) {
	users, err := model.GetUsers(store)
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{}, data: users}, nil
}

func getGamesNearLocation(store model.Store, req yoda.Request) (tr.Result, error) {
	data := map[string]string{
		"lng":    req.Params["lng"],
		"lat":    req.Params["lat"],
		"userid": req.Session["userid"].(string),
	}
	games, err := game.GetGamesNearLocation(store, data, req.Session["username"].(string))
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{"error": "false"}, data: games}, nil
}

func getGame(store model.Store, req yoda.Request) (tr.Result, error) {
	data := map[string]string{
		"code":     req.Params["id"],
		"username": req.Session["username"].(string),
		"userid":   req.Session["userid"].(string),
	}
	game, err := game.GetGame(store, data)
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{"error": "false"}, data: game}, nil
}

func joinGame(store model.Store, req yoda.Request) (tr.Result, error) {
	data, err := tr.ReadableMap(req.GetData())
	data["userid"] = req.Session["userid"].(string)
	data["username"] = req.Session["username"].(string)

	if err != nil {
		return nil, err
	}
	g := new(client.Game)
	switch req.Params["mode"] {
	case "i":
		fallthrough
	case "j":
		g, err = game.JoinGame(store, data, false, req.SocketManager)
	default:
		err = yoda.ErrInvalidArgument
	}

	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"successful": "true"}, data: g}, nil
}

// rateGame handles rating a game
func rateGame(store model.Store, req yoda.Request) (tr.Result, error) {
	data, ok := req.GetData().(map[string]interface{})
	if !ok {
		return nil, yoda.ErrInvalidData
	}

	data["username"] = req.Session["username"].(string)
	err := game.Rate(store, data)
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{"error": "false"}, data: map[string]interface{}{"successful": true}}, nil
}

func exitGame(store model.Store, req yoda.Request) (tr.Result, error) {
	user := &owner{
		id:   req.Session["userid"].(string),
		name: req.Session["username"].(string),
	}
	err := game.ExitGame(store, user, req.SocketManager)
	if err != nil {
		return nil, err
	}
	return result{meta: map[string]string{"error": "false"}, data: map[string]interface{}{"successful": true}}, nil
}

//getUserInfo will get friends
func getUserInfo(store model.Store, req yoda.Request, wantFriends bool) (tr.Result, error) {
	data, err := tr.ReadableMap(req.GetData())
	if err != nil {
		return nil, err
	}
	username, ok := data["username"]
	if !ok {
		return nil, yoda.ErrMissingData("username", "string")
	}

	ownerName := req.Session["username"].(string)
	ownerID := req.Session["userid"].(string)

	var info interface{}
	user, err := client.GetUserAndPopulate(store, username, ownerName, ownerID)
	if err != nil {
		return nil, err
	}

	if wantFriends {
		info = user.Friends
	} else {
		info = user.PrevGames
	}

	return result{meta: nil, data: info}, err
}

// getUser returns a client-based user
func getUser(store model.Store, req yoda.Request) (tr.Result, error) {
	data, err := tr.ReadableMap(req.GetData())
	if err != nil {
		return nil, err
	}

	username, _ := data["username"]
	ownerName := req.Session["username"].(string)
	ownerID := req.Session["userid"].(string)

	if username == "" {
		username = ownerName
	}

	var user *client.User

	switch req.Params["argument"] {
	case "1":
		user, err = client.GetUserAndPopulate(store, username, ownerName, ownerID)
	default:
		user, err = client.GetUser(store, username, ownerName, ownerID)
	}

	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"error": "false"}, data: user}, nil
}

// createOp starts a creation operation
// e.g creating a user
func createOp(store model.Store, req yoda.Request) (tr.Result, error) {
	var res tr.Result
	var err error

	switch req.Params["entity"] {
	case "user":
		res, err = account.CreateUser(store, req.Result)
	case "game":
		res, err = game.CreateGame(store, req.Result, req.Session["username"].(string))
	default:
		err = fmt.Errorf("Unknown Create Entity: %s", req.Params["entity"])
	}

	return res, err
}

// edit starts an edit operation to modify an entity
func edit(store model.Store, req yoda.Request) (tr.Result, error) {
	var err error
	var res tr.Result

	user := &owner{
		id:   req.Session["userid"].(string),
		name: req.Session["username"].(string),
	}

	switch req.Params["entity"] {
	case "user":
		res, err = account.EditUser(store, user, req.Result)
	case "game":
		data, ok := req.GetData().(map[string]interface{})
		if !ok {
			return nil, yoda.ErrInvalidData
		}
		res, err = game.EditGame(store, req.Session["username"].(string), data)
	default:
		err = fmt.Errorf("Unknown Entity: %s", req.Params["entity"])
	}

	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"error": "false"}, data: res.GetData()}, nil
}

//friendUserOp decides whether we are wanting friends or user
func friendUserOp(store model.Store, req yoda.Request) (tr.Result, error) {
	var res tr.Result
	var err error

	switch req.Params["command"] {
	case "p":
		res, err = getUser(store, req)
	case "f":
		res, err = getUserInfo(store, req, true)
	case "pg":
		res, err = getUserInfo(store, req, false)
	default:
		err = fmt.Errorf("Unknown Command: %s", req.Params["command"])
	}

	return res, err
}

func login(store model.Store, req yoda.Request) (tr.Result, error) {
	reqCred, err := tr.ReadableMap(req.GetData())
	if err != nil {
		return result{meta: map[string]string{"code": "204"}, data: "user found"}, err
	}
	cred := map[string]string{
		"username": reqCred["username"],
		"password": reqCred["password"],
		"csrf":     reqCred["csrf"],
	}
	res, err := authentication.Login(store.GetUsers(), cred)
	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"code": "204"}, data: res.GetData()}, nil
}

func remove(store model.Store, req yoda.Request) (tr.Result, error) {
	data, err := tr.ReadableMap(req.GetData())
	if err != nil {
		return nil, err
	}

	fromUser := req.Session["username"].(string)
	removeUser := data["unfriend"]

	err = invites.RemoveFriend(store, fromUser, removeUser, req.SocketManager)
	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"error": "false"}, data: map[string]interface{}{"successful": true}}, err
}

func invite(store model.Store, req yoda.Request) (tr.Result, error) {
	multidata, ok := req.Result.(map[string]interface{})
	if !ok {
		return nil, yoda.ErrInvalidData
	}

	data, err := tr.ReadableMap(multidata)
	if err != nil {
		return nil, yoda.ErrInvalidData
	}

	switch req.Params["mode"] {
	case "send":

		//No error check becsuse the session safely stores the username as a string
		from := req.Session["username"].(string)
		to := data["to"]

		err = invites.SendInvite(store, req.Params["type"], from, to, req.SocketManager)
	case "cancel":
		from := req.Session["username"].(string)
		to := data["to"]

		err = invites.Cancel(store, req.Params["type"], from, to, req.SocketManager)
	case "review":
		to := req.Session["username"].(string)
		from := data["from"]

		accepted, ok := multidata["accept"].(bool)
		if !ok {
			return nil, yoda.ErrInvalidData
		}
		gameid, ok := multidata["game"].(string)
		if !ok {
			gameid = ""
		}

		err = invites.ReviewInvite(store, req.Params["type"], from, to, gameid, accepted, req.SocketManager)
	default:
		err = fmt.Errorf("Unknown mode of invite: %s", req.Params["mode"])
	}

	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"error": "false"}, data: map[string]interface{}{"successful": true}}, err
}

// uploadFile handles uploading a file to the db
func uploadFile(store model.Store, req yoda.Request) (tr.Result, error) {
	data, _ := req.GetData().(map[string]interface{})

	data["username"] = req.Session["username"]
	err := upload.NewFile(store, data)
	if err != nil {
		return nil, err
	}

	return result{meta: map[string]string{"error": "false"}, data: map[string]bool{"successful": true}}, nil
}
