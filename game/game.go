package game

import (
	// Imported for the JSON marshal/unmarshal interfaces.
	_ "encoding/json"
	"errors"
	"log"
)

// DefaultHandSize is the default number of cards
// maintained in players' hands.
const DefaultHandSize int = 10

// Phase represents which phase the game is currently in.
type Phase int

const (
	// Lobby is when play has yet to start.
	Lobby Phase = iota

	// RoundInProgress is when awaiting player card submissions.
	RoundInProgress

	// WinnerSelection is when the Czar is meant to select a winning card.
	WinnerSelection

	// EndOfRound is when the round is over, but the next has yet to begin.
	EndOfRound

	// EndOfGame is when the game is over and further action is required.
	EndOfGame
)

// MarshalJSON attempts to serialise the phase as a JSON string.
func (p Phase) MarshalJSON() (result []byte, err error) {
	options := [...]string{"lobby", "roundInProgress", "winnerSelection", "endOfRound", "endOfGame"}
	if int(p) < 0 || int(p) >= len(options) {
		return nil, errors.New("invalid game phase")
	}
	result = []byte(options[p])
	return
}

// UnmarshalJSON attempts to deserialise phase from a JSON string.
func (p *Phase) UnmarshalJSON(input []byte) (err error) {
	switch string(input) {
	case "lobby":
		*p = Lobby
	case "roundInProgress":
		*p = RoundInProgress
	case "winnerSelection":
		*p = WinnerSelection
	case "endOfRound":
		*p = EndOfRound
	case "endOfGame":
		*p = EndOfGame
	default:
		err = errors.New("invalid Phase value")
	}
	return
}

// Player represents a user who has joined a game.
type Player struct {
	ID       int
	Username string
	Hand     []*AnswerCard
	Score    int
}

// Game represents the state of a single game.
type Game struct {
	ID        int
	Decks     []*Deck
	Phase     Phase
	MaxPoints int
	Name      string
	Owner     Player
	Password  string
	PlayDeck  PlayDeck
	Players   []Player
	Round     *Round
}

// Create attempts to create a new game.
func Create(id int, name, password string) (game *Game, err error) {
	if len(name) < 4 {
		return nil, errors.New("game name must be at least 4 characters")
	}

	game = &Game{
		ID:       id,
		Name:     name,
		Password: password,
	}

	return game, nil
}

// Start moves the game state to `InProgress` and deals
// cards to joined players.
func (g *Game) Start() (err error) {
	if g.MaxPoints < 1 {
		return errors.New("max points not set")
	}

	g.PlayDeck.Init(g.Decks[0])

	g.DealAll(DefaultHandSize)
	g.Phase = RoundInProgress
	card, err := g.PlayDeck.DrawQuestion()
	g.Round = &Round{
		Czar:     &g.Players[0],
		Question: card,
	}
	return
}

// DealAll deals cards to all joined players.
func (g *Game) DealAll(upTo int) {
	for _, player := range g.Players {
		g.Deal(&player, 1)
	}
}

// Deal deals cards to a single player.
func (g *Game) Deal(player *Player, upTo int) {
	numNew := len(player.Hand) - upTo
	for i := 0; i < numNew; i++ {
		card, err := g.PlayDeck.DrawAnswer()
		if err != nil {
			log.Println(err)
			continue
		}
		player.Hand = append(player.Hand, card)
	}
}

// Round represents the state of the current game round.
type Round struct {
	CardSubmissions []CardSubmission
	Czar            *Player
	Question        *QuestionCard
	Winner          *Player
}

// CardSubmission represents a player's submission for their
// answer to the Czar's question.
type CardSubmission struct {
	Cards  []AnswerCard
	Player *Player
}

// SetMaxPoints changes the number of points required to win
// the game.
//
// Will fail if the game is outside the lobby phase or if
// the value isn't in the range [3, 10].
func (g *Game) SetMaxPoints(maxPoints int) (err error) {
	if maxPoints < 3 {
		return errors.New("cannot have max points under 3")
	} else if maxPoints >= 10 {
		return errors.New("cannot have max points above 10")
	}

	g.MaxPoints = maxPoints
	return
}

// SetName changes the name of the game.
//
// Will fail if the game is outside the lobby phase.
func (g *Game) SetName(name string) (err error) {
	if g.Phase != Lobby {
		return errors.New("can't change game name outside lobby phase")
	}

	g.Name = name
	return
}
