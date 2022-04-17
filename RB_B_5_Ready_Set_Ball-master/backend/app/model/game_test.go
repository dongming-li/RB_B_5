package model

import (
	"reflect"
	"testing"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"gopkg.in/mgo.v2/bson"
)

func Test_getGame(t *testing.T) {
	type args struct {
		c          Collection
		identifier Identity
	}
	tests := []struct {
		name    string
		args    args
		want    *Game
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGame(tt.args.c, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createGameEntry(t *testing.T) {
	st := time.Now()
	et := time.Now().Local().Add(time.Hour * time.Duration(2))
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *newGame
		wantErr error
	}{
		// Test cases.
		{
			name: "valid game",
			args: args{map[string]interface{}{"name": "VictorBall", "startTime": st.Format(time.RFC3339), "endTime": et.Format(time.RFC3339), "host": "kerno",
				"sport": 1.0, "minAge": 14.0, "maxAge": 58.0, "lat": 16.6, "lng": 24.6, "private": true}},
			want: &newGame{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "kerno",
				sport:     1.0,
				minAge:    14.0,
				maxAge:    58.0,
				lat:       16.6,
				lng:       24.6,
				private:   true,
			},
			wantErr: nil,
		},
		{
			name: "invalid game",
			args: args{map[string]interface{}{"name": "VictorBall", "startTime": st.Format(time.RFC3339), "endTime": et.Format(time.RFC3339), "host": "kerno",
				"sport": -3.0, "minAge": 14.0, "maxAge": 58.0, "lat": 16.6, "lng": 24.6}},
			want:    nil,
			wantErr: yoda.ClientError{Code: yerr.InvalidSport, Message: "Invalid Sport"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createGameEntry(tt.args.data)
			if !errorEquals(err, tt.wantErr) {
				t.Errorf("createGameEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createGameEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGame_OK(t *testing.T) {
	st := time.Now()
	et := time.Now().Local().Add(time.Hour * time.Duration(2))
	type fields struct {
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
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		// Test cases.
		{
			name: "invalid game name",
			fields: fields{
				name: "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidGameName, Message: "Invalid Game Name"},
		},
		{
			name: "invalid time range",
			fields: fields{
				name:      "VictorBall",
				startTime: et.Format(time.RFC3339),
				endTime:   st.Format(time.RFC3339),
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidTimeRange, Message: "Invalid Time Range"},
		},
		{
			name: "invalid time",
			fields: fields{
				name:      "VictorBall",
				startTime: 0,
				endTime:   st.Format(time.RFC3339),
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidTimeRange, Message: "Invalid Time Range"},
		},
		{
			name: "invalid host",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "",
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidHost, Message: "Invalid Host"},
		},
		{
			name: "invalid sport",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     9.0,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidSport, Message: "Invalid Sport"},
		},
		{
			name: "invalid age old",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    16.0,
				maxAge:    123.0,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidAgeRange, Message: "Invalid age range"},
		},
		{
			name: "invalid age young",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    11.0,
				maxAge:    42.0,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidAgeRange, Message: "Invalid age range"},
		},
		{
			name: "invalid age range",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    84.0,
				maxAge:    32.0,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidAgeRange, Message: "Invalid age range"},
		},
		{
			name: "invalid longitude",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    16.0,
				maxAge:    32.0,
				lat:       54,
				lng:       183,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidLocation, Message: "Invalid location"},
		},
		{
			name: "invalid latitude",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    16.0,
				maxAge:    32.0,
				lat:       183,
				lng:       -64,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidLocation, Message: "Invalid location"},
		},
		{
			name: "invalid private game value",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      "victor",
				sport:     3.0,
				minAge:    16.0,
				maxAge:    32.0,
				lat:       50.0,
				lng:       64.0,
				private:   2,
			},
			wantErr: yoda.ClientError{Code: yerr.InvalidPrivateGame, Message: "Invalid private game value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &newGame{
				name:      tt.fields.name,
				startTime: tt.fields.startTime,
				endTime:   tt.fields.endTime,
				host:      tt.fields.host,
				sport:     tt.fields.sport,
				minAge:    tt.fields.minAge,
				maxAge:    tt.fields.maxAge,
				lat:       tt.fields.lat,
				lng:       tt.fields.lng,
				private:   tt.fields.private,
			}
			if err := g.OK(); !errorEquals(err, tt.wantErr) {
				t.Errorf("newGame.OK() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newGame_createGame(t *testing.T) {
	st := time.Now()
	et := time.Now().Local().Add(time.Hour * time.Duration(2))
	type fields struct {
		name      interface{}
		startTime interface{}
		endTime   interface{}
		host      interface{}
		sport     interface{}
		minAge    interface{}
		maxAge    interface{}
		lat       interface{}
		lng       interface{}
		private   interface{}
	}
	type args struct {
		c Collection
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Game
		wantErr error
	}{
		{
			name: "createGame valid",
			fields: fields{
				name:      "VictorBall",
				startTime: st.Format(time.RFC3339),
				endTime:   et.Format(time.RFC3339),
				host:      bson.ObjectId("59cec0369ae3352c3507d5c1"),
				sport:     3.0,
				minAge:    16.0,
				maxAge:    32.0,
				lat:       34.0,
				lng:       -64.0,
				private:   true,
			},
			args: args{mockCollection{}},
			want: &Game{
				Name:      "VictorBall",
				StartTime: st,
				Host:      bson.ObjectIdHex("59cec0369ae3352c3507d5c1"),
				Sport:     3.0,
				AgeRange:  [2]int8{16, 32},
				Location:  Location{Lat: 34.0, Lng: -64.0},
				EndTime:   et,
				Rating:    Rating{0, 0},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGame{
				name:      tt.fields.name,
				startTime: tt.fields.startTime,
				endTime:   tt.fields.endTime,
				host:      tt.fields.host,
				sport:     tt.fields.sport,
				minAge:    tt.fields.minAge,
				maxAge:    tt.fields.maxAge,
				lat:       tt.fields.lat,
				lng:       tt.fields.lng,
				private:   tt.fields.private,
			}
			got, err := g.createGame(tt.args.c, tt.args.c)
			if err != tt.wantErr {
				t.Errorf("newGame.createGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.equals(tt.want) {
				t.Errorf("newGame.createGame() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_sport_OK(t *testing.T) {
	tests := []struct {
		name string
		s    sport
		want bool
	}{
		{
			name: "Valid sport",
			s:    sport(0),
			want: true,
		},
		{
			name: "Valid sport",
			s:    sport(len(sports) - 1),
			want: true,
		},
		{
			name: "Sport index too low",
			s:    sport(-1),
			want: false,
		},
		{
			name: "Sport index too high",
			s:    sport(len(sports)),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.OK(); got != tt.want {
				t.Errorf("sport.OK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_ChangeAgeRange(t *testing.T) {
	type fields struct {
		ID          bson.ObjectId
		Name        string
		StartTime   time.Time
		Location    Location
		Host        bson.ObjectId
		Members     []bson.ObjectId
		MemsInvited []bson.ObjectId
		EndTime     time.Time
		Sport       sport
		Rating      Rating
		AgeRange    [2]int8
		JoinCode    string
	}
	type args struct {
		startAge interface{}
		endAge   interface{}
		errC     chan error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Invalid age range",
			args: args{
				startAge: 28.0,
				endAge:   24.0,
				errC:     make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Not okay",
			args: args{
				startAge: 28.0,
				endAge:   24,
				errC:     make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Not ok",
			args: args{
				startAge: 28,
				endAge:   24.0,
				errC:     make(chan error),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				StartTime:   tt.fields.StartTime,
				Location:    tt.fields.Location,
				Host:        tt.fields.Host,
				Members:     tt.fields.Members,
				MemsInvited: tt.fields.MemsInvited,
				EndTime:     tt.fields.EndTime,
				Sport:       tt.fields.Sport,
				Rating:      tt.fields.Rating,
				AgeRange:    tt.fields.AgeRange,
				JoinCode:    tt.fields.JoinCode,
			}
			go g.changeAgeRange(tt.args.startAge, tt.args.endAge, tt.args.errC)
			if err := <-tt.args.errC; (err != nil) != tt.wantErr {
				t.Errorf("Game.ChangeAgeRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ChangeName(t *testing.T) {
	type fields struct {
		ID          bson.ObjectId
		Name        string
		StartTime   time.Time
		Location    Location
		Host        bson.ObjectId
		Members     []bson.ObjectId
		MemsInvited []bson.ObjectId
		EndTime     time.Time
		Sport       sport
		Rating      Rating
		AgeRange    [2]int8
		JoinCode    string
	}
	type args struct {
		n    interface{}
		errC chan error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Not ok",
			args: args{
				n:    6,
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Not valid name",
			args: args{
				n:    "hi!!",
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Valid name",
			args: args{
				n:    "Cyclones",
				errC: make(chan error),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				StartTime:   tt.fields.StartTime,
				Location:    tt.fields.Location,
				Host:        tt.fields.Host,
				Members:     tt.fields.Members,
				MemsInvited: tt.fields.MemsInvited,
				EndTime:     tt.fields.EndTime,
				Sport:       tt.fields.Sport,
				Rating:      tt.fields.Rating,
				AgeRange:    tt.fields.AgeRange,
				JoinCode:    tt.fields.JoinCode,
			}
			go g.changeName(tt.args.n, tt.args.errC)
			if err := <-tt.args.errC; (err != nil) != tt.wantErr {
				t.Errorf("Game.ChangeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ChangeTime(t *testing.T) {
	st := time.Now()
	et := time.Now().Local().Add(time.Hour * time.Duration(2))
	type fields struct {
		ID          bson.ObjectId
		Name        string
		StartTime   time.Time
		Location    Location
		Host        bson.ObjectId
		Members     []bson.ObjectId
		MemsInvited []bson.ObjectId
		EndTime     time.Time
		Sport       sport
		Rating      Rating
		AgeRange    [2]int8
		JoinCode    string
	}
	type args struct {
		start interface{}
		end   interface{}
		errC  chan error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Not ok",
			args: args{
				start: et,
				end:   st,
				errC:  make(chan error),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				StartTime:   tt.fields.StartTime,
				Location:    tt.fields.Location,
				Host:        tt.fields.Host,
				Members:     tt.fields.Members,
				MemsInvited: tt.fields.MemsInvited,
				EndTime:     tt.fields.EndTime,
				Sport:       tt.fields.Sport,
				Rating:      tt.fields.Rating,
				AgeRange:    tt.fields.AgeRange,
				JoinCode:    tt.fields.JoinCode,
			}
			go g.changeTime(tt.args.start, tt.args.end, tt.args.errC)
			if err := <-tt.args.errC; (err != nil) != tt.wantErr {
				t.Errorf("Game.ChangeTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ChangeLocation(t *testing.T) {
	type fields struct {
		ID          bson.ObjectId
		Name        string
		StartTime   time.Time
		Location    Location
		Host        bson.ObjectId
		Members     []bson.ObjectId
		MemsInvited []bson.ObjectId
		EndTime     time.Time
		Sport       sport
		Rating      Rating
		AgeRange    [2]int8
		JoinCode    string
	}
	type args struct {
		lat  interface{}
		lng  interface{}
		errC chan error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Not ok",
			args: args{
				lat:  6,
				lng:  8.0,
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Not okay",
			args: args{
				lat:  6.0,
				lng:  8,
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Invalid latitude",
			args: args{
				lat:  183.0,
				lng:  8.0,
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Invalid longitude",
			args: args{
				lat:  12.0,
				lng:  883.0,
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Valid location",
			args: args{
				lat:  50.0,
				lng:  50.0,
				errC: make(chan error),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				StartTime:   tt.fields.StartTime,
				Location:    tt.fields.Location,
				Host:        tt.fields.Host,
				Members:     tt.fields.Members,
				MemsInvited: tt.fields.MemsInvited,
				EndTime:     tt.fields.EndTime,
				Sport:       tt.fields.Sport,
				Rating:      tt.fields.Rating,
				AgeRange:    tt.fields.AgeRange,
				JoinCode:    tt.fields.JoinCode,
			}
			go g.changeLocation(tt.args.lat, tt.args.lng, tt.args.errC)
			if err := <-tt.args.errC; (err != nil) != tt.wantErr {
				t.Errorf("Game.ChangeLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ChangeSport(t *testing.T) {
	type fields struct {
		ID          bson.ObjectId
		Name        string
		StartTime   time.Time
		Location    Location
		Host        bson.ObjectId
		Members     []bson.ObjectId
		MemsInvited []bson.ObjectId
		EndTime     time.Time
		Sport       sport
		Rating      Rating
		AgeRange    [2]int8
		JoinCode    string
	}
	type args struct {
		errC chan error
		s    interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid sport (soccer)",
			args: args{
				s:    float64(0),
				errC: make(chan error),
			},
			wantErr: false,
		},
		{
			name: "Invalid sport, too large index",
			args: args{
				s:    float64(8),
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Invalid sport, too small index",
			args: args{
				s:    float64(-1),
				errC: make(chan error),
			},
			wantErr: true,
		},
		{
			name: "Invalid type, sport has to be float64",
			args: args{
				s:    int(0),
				errC: make(chan error),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				StartTime:   tt.fields.StartTime,
				Location:    tt.fields.Location,
				Host:        tt.fields.Host,
				Members:     tt.fields.Members,
				MemsInvited: tt.fields.MemsInvited,
				EndTime:     tt.fields.EndTime,
				Sport:       tt.fields.Sport,
				Rating:      tt.fields.Rating,
				AgeRange:    tt.fields.AgeRange,
				JoinCode:    tt.fields.JoinCode,
			}
			go g.changeSport(tt.args.s, tt.args.errC)
			if err := <-tt.args.errC; (err != nil) != tt.wantErr {
				t.Errorf("Game.ChangeSport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (g *Game) equals(o *Game) bool {
	return g.AgeRange[0] == o.AgeRange[0] &&
		g.AgeRange[1] == o.AgeRange[1] &&
		g.Name == o.Name &&
		g.Sport == o.Sport &&
		g.Location == o.Location
	// g.Host == o.Host
	// g.StartTime == o.StartTime
}
