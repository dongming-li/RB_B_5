package invites

import (
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/game"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/socket_handler"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"
)

// Represent the type of a request
const (
	Friend = "0"
	Team   = "1"
)

// SendInvite sends an invite from
// with username [from] to user with username [to]
// typeOfRequest:
// 		0 - friend request
// 		1 - team request
func SendInvite(store model.Store, typeOfRequest, from, to string, sm *sockets.Manager) error {
	c := store.GetUsers()

	fromUser, toUser, err := getFromAndToUsers(c, from, to)
	if typeOfRequest == Friend {
		if model.AreFriends(fromUser, toUser) {
			return &yoda.ClientError{Code: yerr.InvalidInvite, Message: "The user cannot send a request to its friend"}
		}

		err = fromUser.SendFriendRequest(c, toUser)
		if err != nil {
			return err
		}

		sm.SendMessage(map[string]interface{}{
			"name": sockethandler.ReceivedFriendInvite,
			"from": from,
			"to":   to,
		}, []string{toUser.ID.Hex()})

	} else if typeOfRequest == Team {
		err = toUser.AddGameRequest(c, store.GetCollection("game"), store.GetCollection("prevgame"), fromUser.Username, fromUser.CurrentGame.Hex())
		if err != nil {
			return err
		}

		go sm.SendMessage(map[string]interface{}{
			"name": sockethandler.ReceivedGameInvite,
			"to":   to,
		}, []string{toUser.ID.Hex()})

		return err
	}

	return nil
}

// ReviewInvite adds or rejects an invite/ game invite from another user
// typeOfRequest:
// 		0 - friend request
// 		1 - team request
func ReviewInvite(store model.Store, typeOfInvite, from, to, g string, accept bool, sm *sockets.Manager) error {
	c := store.GetUsers()

	fromUser, toUser, err := getFromAndToUsers(c, from, to)
	if err != nil {
		return err
	}

	switch typeOfInvite {
	case Friend:
		if accept {
			err = toUser.AddFriend(c, fromUser, true)
		} else {
			err = toUser.DeleteFriendRequest(c, fromUser)
		}

		if err != nil {
			return err
		}

		go sm.SendMessage(map[string]interface{}{
			"name":   sockethandler.ResponseFriendInvite,
			"to":     to,
			"accept": accept,
		}, []string{fromUser.ID.Hex()})

		return nil
	case Team:
		if accept {
			data := map[string]string{"code": g, "username": to, "userid": toUser.ID.Hex()}
			_, err := game.JoinGame(store, data, true, sm)
			if err != nil {
				return err
			}
		}

		err = toUser.RemoveGameRequest(c, store.GetCollection("prevgame"), g)
		if err != nil {
			return err
		}
	default:
		return yoda.ErrInvalidArgument
	}

	return nil
}

// Cancel cancels a send invite to a user
// with username [from] to user with username [to]
// typeOfRequest:
// 		0 - friend request
// 		1 - team request
func Cancel(store model.Store, typeOfRequest, from, to string, sm *sockets.Manager) error {
	c := store.GetUsers()

	fromUser, toUser, err := getFromAndToUsers(c, from, to)
	if err != nil {
		return err
	}

	switch typeOfRequest {
	case Friend:
		err = toUser.DeleteFriendRequest(c, fromUser)
	default:
		err = yoda.ErrInvalidArgument
	}

	if err != nil {
		return err
	}

	sm.SendMessage(map[string]interface{}{
		"name": sockethandler.CancelFriendInvite,
		"from": from,
	}, []string{toUser.ID.Hex()})

	return nil
}

// getFromAndToUsers returns the from and to users based on the [from] and [to] usernames
func getFromAndToUsers(c model.Collection, from, to string) (*model.User, *model.User, error) {
	//validate input
	// This error should never happen because this the 'from' user is gotten from the session
	if !validation.IsValidUserName(from) {
		return nil, nil, &yoda.ClientError{Code: yerr.InvalidInvite, Message: "The 'from' username is invalid"}
	}
	if !validation.IsValidUserName(to) {
		return nil, nil, &yoda.ClientError{Code: yerr.InvalidInvite, Message: "The 'to' username is invalid"}
	}

	if from == to {
		return nil, nil, yoda.ErrSameUserInvite
	}

	//There is no check to see if the 'to' or 'from' user is present
	//because [getUser] already checks this

	fromUser, err := model.GetUserByUsername(c, from)
	if err != nil {
		return nil, nil, err
	}

	toUser, err := model.GetUserByUsername(c, to)
	if err != nil {
		return nil, nil, err
	}

	return fromUser, toUser, nil
}
