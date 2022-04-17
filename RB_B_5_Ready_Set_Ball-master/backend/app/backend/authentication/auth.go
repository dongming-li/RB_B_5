package authentication

import (
	tr "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/transaction"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
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

// Login returns an [error] if user cannot be given a session
func Login(c model.Collection, cred map[string]string) (tr.Result, error) { //TODO: return a session?
	//TODO: check csrf
	user, err := model.GetUserByUsername(c, cred["username"])
	if err == nil && validation.CheckPassword(user.Password, cred["password"]) {
		return &result{data: map[string]interface{}{
			"userid":    user.ID.Hex(), //TODO:remove
			"avatar":    user.Avatar,
			"email":     user.Email,
			"username":  user.Username,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"rating":    user.TeamRating,
		}}, nil
	}

	return nil, &yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username or password"}
}


