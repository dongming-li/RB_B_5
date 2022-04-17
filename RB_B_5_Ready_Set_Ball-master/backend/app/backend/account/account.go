package account

import (
	tr "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/transaction"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
)

// account Ooerations
const (
	create operation = iota
	edit
)

func do(store model.Store, rData interface{}, user *model.User, op operation) (*result, error) {
	res := new(result)
	data, err := tr.ReadableMap(rData.(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	switch op {
	case create:
		err = model.CreateUser(store, data)
	case edit:
		err = user.Edit(store.GetUsers(), data)
		res.data = map[string]interface{}{
			"successful": true,
			"username":   user.Username,
			"userid":     user.ID.Hex(),
			"email":      user.Email,
		}
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateUser attempts to create a user from [data] and returns a [transaction.Result] and an [error]
func CreateUser(store model.Store, data interface{}) (tr.Result, error) {
	return do(store, data, nil, create)
}

// EditUser attempts to edit a user with [data] and returns a [transaction.Result] and an [error]
func EditUser(store model.Store, owner tr.Owner, data interface{}) (tr.Result, error) {
	u, err := model.GetUserByUsername(store.GetUsers(), owner.Username())
	if err != nil {
		return nil, err
	}

	return do(store, data, u, edit)
}
