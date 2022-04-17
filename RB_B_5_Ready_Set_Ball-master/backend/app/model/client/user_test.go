package client

import (
	"testing"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
)

func TestUser_populate(t *testing.T) {
	type fields struct {
		Username       string
		Firstname      string
		Lastname       string
		PersonalRating float32
		TeamRating     float64
		Friends        []User
		Preferences    map[string]interface{}
	}
	type args struct {
		friends []model.User
	}
	type want []User
	tests := []struct {
		name    string
		fields  fields
		args    args
		friends want
	}{
		{
			name:   "Adds all friends",
			fields: fields{},
			args: args{
				friends: []model.User{
					model.User{Username: "one"}, model.User{Username: "two"},
					model.User{Username: "three"}, model.User{Username: "four"},
				},
			},
			friends: []User{
				{Username: "one"}, {Username: "two"},
				{Username: "three"}, {Username: "four"},
			},
		},
		{
			name:   "Handles nil friends",
			fields: fields{},
			args: args{
				friends: nil,
			},
			friends: want{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Username:       tt.fields.Username,
				Firstname:      tt.fields.Firstname,
				Lastname:       tt.fields.Lastname,
				PersonalRating: tt.fields.PersonalRating,
				TeamRating:     tt.fields.TeamRating,
				Friends:        tt.fields.Friends,
			}
			u.populate(tt.args.friends)
			if len(u.Friends) != len(tt.friends) {
				t.Errorf("user.populate(%v), got len(u.Friends) as %v, expected len(u.Friends) to be %v", tt.args, len(u.Friends), len(tt.friends))
			}
			for i := range u.Friends {
				if u.Friends[i].Username != tt.friends[i].Username {
					t.Errorf("All friends should be added,  u.Friends[%d].Username = %v, tt.friends[%d].Username %v", i, u.Friends[i], i, tt.friends[i])
				}
			}
		})
	}
}
