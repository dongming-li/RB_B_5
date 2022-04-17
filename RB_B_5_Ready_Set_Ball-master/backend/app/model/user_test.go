package model

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
)

func Test_getUser(t *testing.T) {
	type args struct {
		c          Collection
		identifier Identity
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		// Test cases.
		{
			name:    "Empty user",
			args:    args{c: mockCollection{}, identifier: mockIdentifier{}},
			want:    &User{},
			wantErr: false,
		},
		{
			name:    "Identify by username",
			args:    args{c: mockCollection{}, identifier: mockIdentifier{key: "username", value: "mich8787"}},
			want:    &User{Username: "mich8787"},
			wantErr: false,
		},
		{
			name:    "Identify by id",
			args:    args{c: mockCollection{}, identifier: mockIdentifier{key: "id", value: "thisisarandomid"}},
			want:    &User{ID: "thisisarandomid"},
			wantErr: false,
		},
		{
			name:    "Invalid user",
			args:    args{c: mockCollection{}, identifier: mockIdentifier{key: "username", value: "invaliduser"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUser(tt.args.c, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	type args struct {
		c     Collection
		uName string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr error
	}{
		//Test cases.
		{
			name:    "username with consecutive underscores",
			args:    args{c: mockCollection{}, uName: "invalid__username"},
			want:    nil,
			wantErr: &yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username/password"},
		},
		{
			name:    "valid username",
			args:    args{c: mockCollection{}, uName: "mich8787"},
			want:    &User{Username: "mich8787"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserByUsername(tt.args.c, tt.args.uName)

			if !errorEquals(err, tt.wantErr) {
				t.Errorf("GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeNewUserEntry(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *newUser
		wantErr error
	}{
		// Test cases.
		{
			name: "valid data",
			args: args{map[string]string{"username": "mich8787", "password": "TH15P@$sWord _", "firstname": "John", "lastname": "Doe",
				"confirmPass": "TH15P@$sWord _", "email": "mich@smart.com", "city": "Jos"}},
			want: &newUser{
				username:    "mich8787",
				password:    "TH15P@$sWord _",
				firstname:   "John",
				lastname:    "Doe",
				confirmPass: "TH15P@$sWord _",
				email:       "mich@smart.com",
				city:        "Jos",
			},
			wantErr: nil,
		},
		{
			name: "email absent",
			args: args{map[string]string{"username": "mich8787", "password": "TH15P@$sWord _", "firstname": "John", "lastname": "Doe",
				"confirmPass": "TH15P@$sWord _"}},
			want:    nil,
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid email"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeNewUserEntry(tt.args.data)
			if !errorEquals(err, tt.wantErr) {
				t.Errorf("makeNewUserEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeNewUserEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUser_OK(t *testing.T) {
	type fields struct {
		username    string
		password    string
		confirmPass string
		firstname   string
		lastname    string
		email       string
		city        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		// Test cases.
		{
			name: "valid new user",
			fields: fields{
				username:    "mich8787",
				password:    "TH15P@$sWord _",
				firstname:   "John",
				lastname:    "Doe",
				confirmPass: "TH15P@$sWord _",
				email:       "mich@smart.com",
				city:        "Jos",
			},
			wantErr: nil,
		},
		{
			name: "invalid firstname",
			fields: fields{
				firstname: "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid firstname"},
		},
		{
			name: "invalid lastname",
			fields: fields{
				firstname: "John",
				lastname:  "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid lastname"},
		},
		{
			name: "invalid username",
			fields: fields{
				username:  "",
				firstname: "John",
				lastname:  "Doe",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid username"},
		},
		{
			name: "invalid password",
			fields: fields{
				username:  "mich8787",
				firstname: "John",
				lastname:  "Doe",
				password:  "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid password"},
		},
		{
			name: "invalid email",
			fields: fields{
				username:    "mich8787",
				firstname:   "John",
				lastname:    "Doe",
				password:    "TH15P@$sWord _",
				confirmPass: "TH15P@$sWord _",
				email:       "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Invalid email"},
		},
		{
			name: "non matching passwords",
			fields: fields{
				username:    "mich8787",
				firstname:   "John",
				lastname:    "Doe",
				password:    "TH15P@$sWord _",
				confirmPass: "TH1SP@$sWord _",
				email:       "mich@smart.com",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidUsernamePassword, Message: "Passwords do not match"},
		},
		{
			name: "incorrect city with number",
			fields: fields{
				username:    "mich8787",
				firstname:   "John",
				lastname:    "Doe",
				password:    "TH15P@$sWord _",
				confirmPass: "TH15P@$sWord _",
				email:       "mich@smart.com",
				city:        "j4tw",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidParameter, Message: "Invalid city"},
		},
		{
			name: "incorrect city short",
			fields: fields{
				username:    "mich8787",
				firstname:   "John",
				lastname:    "Doe",
				password:    "TH15P@$sWord _",
				confirmPass: "TH15P@$sWord _",
				email:       "mich@smart.com",
				city:        "j",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidParameter, Message: "Invalid city"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &newUser{
				username:    tt.fields.username,
				password:    tt.fields.password,
				confirmPass: tt.fields.confirmPass,
				firstname:   tt.fields.firstname,
				lastname:    tt.fields.lastname,
				email:       tt.fields.email,
				city:        tt.fields.city,
			}
			if err := u.OK(); !errorEquals(err, tt.wantErr) {
				t.Errorf("newUser.OK() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	//usernames are changed to lowercase
	user := &newUser{firstname: "John", lastname: "Doe", username: "LuCasKerN"}
	user.OK()
	if user.username != "lucaskern" {
		t.Errorf("Username is to be converted to lowercase. Expected lucaskern, Got %s", user.username)
	}

	//email domains are changed to lowercase
	user = &newUser{firstname: "John", lastname: "Doe", username: "mich8787", password: "TH15P@$sWord _", email: "DaNNy@SMArt.COm"}
	user.OK()
	if user.email != "DaNNy@smart.com" {
		t.Errorf("Email domain is to be converted to lowercase. Expected DaNNy@smart.com, Got %s", user.email)
	}
}

//Mocks
type mockIdentifier struct {
	key   string
	value interface{}
}

func (i mockIdentifier) Identify() map[string]interface{} {
	return map[string]interface{}{i.key: i.value}
}

type mockQueryResult struct {
	all interface{}
	one interface{}
	err error
}

func (r mockQueryResult) All(q interface{}) error {
	q = r.all
	return r.err
}

func (r mockQueryResult) One(q interface{}) error {
	if r.err != nil {
		return r.err
	}
	usr := *r.one.(*User)
	*(q.(*User)) = usr
	return nil
}

func (r mockQueryResult) Apply(change Change, result interface{}) (info ChangeInfo, err error) {
	return nil, nil
}

type mockCollection struct{}

func (c mockCollection) Find(q interface{}) QueryResult {
	key := fmt.Sprint(q)
	res, ok := map[string]QueryResult{
		"map[username:mich8787]":    mockQueryResult{all: []*User{&User{Username: "mich8787"}}, one: &User{Username: "mich8787"}, err: nil},
		"map[id:thisisarandomid]":   mockQueryResult{all: []*User{&User{ID: "thisisarandomid"}}, one: &User{ID: "thisisarandomid"}, err: nil},
		"map[username:invaliduser]": mockQueryResult{all: nil, one: nil, err: errors.New("user not found")},
	}[key]
	if !ok {
		return mockQueryResult{all: []*User{&User{}}, one: &User{}, err: nil}
	}
	return res
}

func (c mockCollection) Insert(docs ...interface{}) error {
	return nil
}

func (c mockCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

func (c mockCollection) Remove(selector interface{}) error {
	return nil
}

func (c mockCollection) FindAndModify(selector interface{}, update interface{}, result interface{}, returnNew bool) (ChangeInfo, error) {
	return nil, nil
}

func (c mockCollection) EnsureIndex(key []string, time time.Duration) error {
	return nil
}

func errorEquals(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return err1 == err2
	}
	return err1.Error() == err2.Error()
}
