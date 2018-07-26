package game

import (
	"io"

	"github.com/gorilla/websocket"
)

type User struct {
	Id       int
	Username string
}

type Player struct {
	User       *User
	Score      int
	Connection websocket.Conn
	Hand       []AnswerCard
}

type Deck struct {
	Id            int
	Name          string
	QuestionCards []QuestionCard
	AnswerCards   []AnswerCard
}

type QuestionCard struct {
	Id         int
	Text       string
	NumAnswers int
}

type AnswerCard struct {
	Id   int
	Text string
}

type Game struct {
	Id              int
	Czar            *Player
	Players         []Player
	CurrentQuestion *QuestionCard
}

type Round struct {
	Winner   *Player
	Question *QuestionCard
}

type CardSubmission struct {
	Player *Player
	Cards  []AnswerCard
}
