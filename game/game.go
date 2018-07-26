package game

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gorilla/websocket"
)

type GamePhase int

const (
	Lobby GamePhase = iota
	RoundInProgress
	EndOfRound
	EndOfGame
)

func (p GamePhase) MarshalJSON() (result []byte, err error) {
	options := [...]string{"lobby", "roundInProgress", "endOfRound", "endOfGame"}
	result = []byte(options[p])
	return
}

func (p *GamePhase) UnmarshalJSON(input []byte) (err error) {
	switch input {
	case "lobby":
		p = Lobby
	case "roundInProgress":
		p = RoundInProgress
	case "endOfRound":
		p = EndOfRound
	case "endOfGame":
		p = EndOfGame
	default:
		err = errors.Error("Invalid GamePhase value")
	}
	return
}

type User struct {
	Id       int
	Username string
}

type Player struct {
	User       *User
	Connection websocket.Conn
	Hand       []AnswerCard
	Score      int
}

type Deck struct {
	Id            int
	AnswerCards   []AnswerCard
	Name          string
	QuestionCards []QuestionCard
}

type QuestionCard struct {
	Id         int
	NumAnswers int
	Text       string
}

type AnswerCard struct {
	Id   int
	Text string
}

type Game struct {
	Id              int
	AnswerDeck      []AnswerCard
	AnswerDiscard   []AnswerCard
	Decks           []Deck
	GamePhase       GamePhase
	MaxPoints       int
	Name            string
	Players         []Player
	QuestionDeck    []QuestionCard
	QuestionDiscard []QuestionCard
	Round           *Round
}

type Round struct {
	CardSubmissions []CardSubmission
	Czar            Player
	Question        QuestionCard
	Winner          *Player
}

type CardSubmission struct {
	Cards  []AnswerCard
	Player *Player
}
