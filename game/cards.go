package game

import (
	// Imported for the JSON marshal/unmarshal interfaces.
	_ "encoding/json"
	"errors"
	"math/rand"
	"time"
)

// Deck represents a deck of answer cards and question cards
// before being put into play.
type Deck struct {
	ID            int            `json:"id"`
	AnswerCards   []AnswerCard   `json:"answerCards"`
	Name          string         `json:"name"`
	QuestionCards []QuestionCard `json:"questionCards"`
}

// QuestionCard represents a black question card.
type QuestionCard struct {
	ID         int    `json:"id"`
	NumAnswers int    `json:"numAnswers"`
	Text       string `json:"text"`
}

// AnswerCard represents a white answer card.
type AnswerCard struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// PlayDeck contains both the answer deck and the question deck.
type PlayDeck struct {
	AnswerDeck   CardDeck
	QuestionDeck CardDeck
}

// CardDeck represents a deck of cards, with a deck portion
// to draw from and a discard pile to discard to.
type CardDeck struct {
	Deck        []interface{}
	DiscardPile []interface{}
}

// Init initialises a card deck with the appropriate type.
func (d *CardDeck) Init(cards []interface{}) (err error) {
	d.Deck = cards
	return
}

// Draw retrieves the top card in the deck and removes it.
func (d *CardDeck) Draw() (card interface{}, err error) {
	if len(d.Deck) < 1 {
		return nil, errors.New("card deck empty")
	}
	deckLength := len(d.Deck) - 1
	card = d.Deck[deckLength]
	d.Deck = d.Deck[:deckLength]
	return
}

// Discard adds a card to the discard pile.
func (d *CardDeck) Discard(card interface{}) (err error) {
	d.Deck = append(d.Deck, card)
	return
}

// Shuffle randomises the order of the non-discard deck.
func (d *CardDeck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(d.Deck); n > 0; n-- {
		randIndex := r.Intn(n)
		d.Deck[n-1], d.Deck[randIndex] = d.Deck[randIndex], d.Deck[n-1]
	}
}

// Reshuffle puts the discard pile back in the deck and shuffles.
func (d *CardDeck) Reshuffle() {
	d.Deck, d.DiscardPile = append(d.Deck, d.Discard), d.DiscardPile[:0]
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(d.Deck); n > 0; n-- {
		randIndex := r.Intn(n)
		d.Deck[n-1], d.Deck[randIndex] = d.Deck[randIndex], d.Deck[n-1]
	}
}

// Init sets up the PlayDeck by loading decks and emptying discard piles.
func (p *PlayDeck) Init(deck Deck) {
	questionCards := make([]interface{}, 0, len(deck.QuestionCards))
	for card := range deck.QuestionCards {
		questionCards = append(questionCards, card)
	}

	answerCards := make([]interface{}, 0, len(deck.AnswerCards))
	for card := range deck.AnswerCards {
		answerCards = append(answerCards, card)
	}

	p.QuestionDeck.Init([]interface{}(questionCards))
	p.AnswerDeck.Init([]interface{}(answerCards))

	p.QuestionDeck.Shuffle()
	p.AnswerDeck.Shuffle()
}

// DrawQuestion removes a question card from the deck and returns it.
//
// If there are no more cards in the deck, the discard pile is shuffled
// and the card drawn from there.
//
// If both piles are empty, an error is returned.
func (p *PlayDeck) DrawQuestion() (card *QuestionCard, err error) {
	if len(p.QuestionDeck.Deck) < 1 {
		p.QuestionDeck.Reshuffle()
	}

	drawnCard, err := p.QuestionDeck.Draw()
	if err != nil {
		return nil, err
	}

	switch v := drawnCard.(type) {
	case *QuestionCard:
		return v, nil
	default:
		return nil, errors.New("non-question card returned")
	}
}

// DiscardQuestion takes the given card and puts it into the discard
// pile.
func (p *PlayDeck) DiscardQuestion(card QuestionCard) (err error) {
	p.QuestionDeck.Discard(card)
	return nil
}

// DrawAnswer removes an answer card from the deck and returns it.
//
// If there are no more cards in the deck, the discard pile is shuffled
// and the card drawn from there.
//
// If both piles are empty, an error is returned.
func (p *PlayDeck) DrawAnswer() (card *AnswerCard, err error) {
	if len(p.AnswerDeck.Deck) < 1 {
		p.AnswerDeck.Reshuffle()
	}

	drawnCard, err := p.AnswerDeck.Draw()
	if err != nil {
		return nil, err
	}

	switch v := drawnCard.(type) {
	case *AnswerCard:
		return v, nil
	default:
		return nil, errors.New("non-answer card returned")
	}
}

// DiscardAnswer takes the given card and puts it into the discard
// pile.
func (p *PlayDeck) DiscardAnswer(card AnswerCard) (err error) {
	p.AnswerDeck.Discard(card)
	return nil
}

func shuffle(vals []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}
