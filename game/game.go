package game

import (
	"encoding/json"
	"errors"
	"io"
	"rand"

	"github.com/gorilla/websocket"
)

const DefaultHandSize int = 10

type GamePhase int

const (
	Lobby GamePhase = iota
	RoundInProgress
	WinnerSelection
	EndOfRound
	EndOfGame
)

func (p GamePhase) MarshalJSON() (result []byte, err error) {
	options := [...]string{"lobby", "roundInProgress", "winnerSelection", "endOfRound", "endOfGame"}
	result = []byte(options[p])
	return
}

func (p *GamePhase) UnmarshalJSON(input []byte) (err error) {
	switch input {
	case "lobby":
		p = Lobby
	case "roundInProgress":
		p = RoundInProgress
	case "winnerSelection":
		p = WinnerSelection
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
	User  *User
	Hand  []AnswerCard
	Score int
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

func (g *Game) Start() (err error) {
	if g.MaxPoints < 1 {
		return errors.New("max points not set")
	}

	if err = g.populateDecks(); err != nil {
		return
	}

	g.DealAll(DefaultHandSize)
	g.GamePhase = RoundInProgress
	g.Round = &Round{
		Czar:     g.Players[0],
		Question: g.QuestionDeck[len(g.QuestionDeck)-1],
	}
	g.QuestionDeck = g.QuestionDeck[len(g.QuestionDeck)-1]

}

func shuffle(vals []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}

func (g *Game) populateDecks() (err error) {
	if len(g.Decks) < 1 {
		return errors.New("empty decks")
	}

	g.AnswerDeck = make([]AnswerCard, 0)
	g.QuestionDeck = make([]QuestionCard, 0)

	for _, deck := range g.Decks {
		for _, card := range deck.AnswerCards {
			g.AnswerDeck = append(g.AnswerDeck, card)
		}
		for _, card := range deck.QuestionCards {
			g.QuestionDeck = append(g.QuestionDeck, card)
		}
	}

	shuffle(g.AnswerDeck)
	shuffle(g.QuestionDeck)

	g.AnswerDiscard = make([]AnswerCard, 0)
	g.QuestionDiscard = make([]QuestionCard, 0)
}

func (g *Game) DealAll(upTo int) {
	for _, player := range g.Players {
		g.Deal(&player)
	}
}

func (g *Game) Deal(player *Player, upTo int) {
	if len(player.Hand) >= upTo {
		return
	}

	numNew := len(player.Hand) - upTo
	player.Hand = append(player.Hand, g.QuestionDeck[len(g.QuestionDeck)-numNew:])
	g.QuestionDeck = g.QuestionDeck[:len(g.QuestionDeck)-numNew]
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
