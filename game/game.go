package game

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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
	Id        int
	Decks     []Deck
	GamePhase GamePhase
	MaxPoints int
	Name      string
	PlayDeck  PlayDeck
	Players   []Player
	Round     *Round
}

func (g *Game) Start() (err error) {
	if g.MaxPoints < 1 {
		return errors.New("max points not set")
	}

	if err = (&g.PlayDeck).Init(g.PlayDeck); err != nil {
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

func (g *Game) DealAll(upTo int) {
	for _, player := range g.Players {
		g.Deal(&player)
	}
}

func (g *Game) Deal(player *Player, upTo int) {
	numNew := len(player.Hand) - upTo
	for i := 0; i < numNew; i++ {
		card, err := g.PlayDeck.AnswerDeck.Draw()
		if err != nil {
			log.Println(err)
			continue
		}
		player.Gain(card)
	}
}

func (p *Player) Gain(card AnswerCard) (err error) {
	player.Hand = append(player.Hand, card)
}

func (p *Player) Draw(source AnswerSource) {
	card, _ := source.Draw()
	player.Hand = append(player.Hand, card)
}

func (p *Player) Discard(discard AnswerDiscard, id int) (err error) {
	for i := 0; i < len(p.Hand); i++ {
		card := p.Hand[i]
		if card.Id == id {
			p.Hand = p.Hand[:i] + p.Hand[:i+1]
			discard.Discard(card)
			return
		}
	}
	return errors.New("card with id not in hand")
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
