package client

import (
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"gopkg.in/mgo.v2/bson"
)

// Game represents a game object
type Game struct {
	ID        bson.ObjectId  `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	StartTime time.Time      `json:"startTime,omitempty"`
	Location  model.Location `json:"location,omitempty"`
	Host      User           `json:"host,omitempty"`
	Members   []User         `json:"members"`
	Duration  float64        `json:"duration,omitempty"`
	Sport     float64        `json:"sport"`
	Rating    float64        `json:"rating,omitempty"`
	AgeRange  [2]int8        `json:"agerange,omitempty"`
	JoinCode  string         `json:"joincode,omitempty"`
	Private   bool           `json:"private"`
}

func newGame(g *model.Game) *Game {
	return &Game{
		ID:        g.ID,
		Name:      g.Name,
		StartTime: g.StartTime,
		Location:  g.Location,
		Duration:  g.EndTime.Sub(g.StartTime).Minutes(),
		Sport:     float64(g.Sport),
		Rating:    g.Rating.Get(),
		AgeRange:  g.AgeRange,
		JoinCode:  g.JoinCode,
		Private:   g.Private,
	}
}

// NewGame returns a client version game populated with [client.User]s
// It uses [uc] as the user collection to populate the data
// It mutates [g] and makes it non-private if the user is in the game
func NewGame(uc model.Collection, g *model.Game, currUser string) (*Game, error) {
	game := newGame(g)
	members, err := g.GetMembers(uc)
	if err != nil {
		return nil, err
	}
	var isInGame bool
	game.Members = make([]User, len(g.Members))
	for i, m := range members {
		if currUser == m.Username {
			isInGame = true
		}
		game.Members[i] = *newUser(m)
	}

	u, err := g.GetHost(uc)
	if err != nil {
		return nil, err
	}
	if currUser == u.Username {
		isInGame = true
	}
	game.Host = *newUser(u)
	if isInGame {
		g.Private = false
	}
	return game, nil
}

// NewGames returns a client version of a slice of games
// populated with [client.User]s
// It uses [uc] as the user collection to populate the data
func NewGames(uc model.Collection, gs []*model.Game, currUser string) ([]*Game, error) {
	games := make([]*Game, 0, len(gs))
	for _, g := range gs {
		game, err := NewGame(uc, g, currUser)
		if err != nil {
			return nil, err
		}
		if !g.Private {
			games = append(games, game)
		}
	}
	return games, nil
}
