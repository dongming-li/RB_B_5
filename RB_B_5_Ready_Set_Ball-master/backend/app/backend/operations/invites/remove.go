package invites

import (
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/socket_handler"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"
)

//RemoveFriend will make two users not be friends any longer
func RemoveFriend(store model.Store, from, toBeRemoved string, sm *sockets.Manager) error {
	//validate input

	if !validation.IsValidUserName(toBeRemoved) {
		return &yoda.ClientError{Code: yerr.InvalidInvite, Message: "The username you are trying to remove is invalid"}
	}

	if from == toBeRemoved {
		return &yoda.ClientError{Code: yerr.InvalidInvite, Message: "The user cannot send request to itself"}
	}

	c := store.GetUsers()

	fromUser, err := model.GetUserByUsername(c, from)
	if err != nil {
		return err
	}

	removeUser, err := model.GetUserByUsername(c, toBeRemoved)
	if err != nil {
		return err
	}

	err = fromUser.RemoveFriend(c, removeUser)
	if err != nil {
		return err
	}

	err = removeUser.RemoveFriend(c, fromUser)
	if err != nil {
		return err
	}

	go sm.SendMessage(map[string]interface{}{
		"name": sockethandler.UnfriendUser,
		"from": fromUser.Username,
	}, []string{removeUser.ID.Hex()})

	return err
}
