package game

import (
	"encoding/json"
	"errors"
)

type Deck struct {
	Id            int            `json:"id"`
	AnswerCards   []AnswerCard   `json:"answerCards"`
	Name          string         `json:"name"`
	QuestionCards []QuestionCard `json:"questionCards"`
}

type QuestionCard struct {
	Id         int    `json:"id"`
	NumAnswers int    `json:"numAnswers"`
	Text       string `json:"text"`
}

type AnswerCard struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type PlayDeck struct {
	AnswerDeck      *AnswerDeck
	QuestionDeck    []QuestionCard
	QuestionDiscard []QuestionCard
}

type AnswerDeck struct {
	Deck        []AnswerCard
	DiscardPile []AnswerCard
}

// Init sets up the PlayDeck by loading decks and emptying discard piles.
func (p *PlayDeck) Init(decks []Deck) (err error) {
	p.AnswerDeck = new(AnswerDeck)
	p.QuestionDeck = new([]QuestionCard)
	p.QuestionDiscard = new([]QuestionCard)

	p.AnswerDeck.Init(decks)

	for _, deck := range decks {
		p.QuestionDeck = append(p.QuestionDeck, deck.QuestionCards)
	}

	shuffle(p.QuestionDeck)
}

func (d *AnswerDeck) Init(decks []Deck) {
	d.Deck = new([]AnswerCard)
	d.DiscardPile = new([]AnswerCard)
	for _, deck := range decks {
		p.Deck = append(p.Deck, deck.AnswerCards)
	}
	shuffle(d.Deck)
}

type AnswerSource interface {
	Draw() (card AnswerCard, err error)
}

type AnswerDiscard interface {
	Discard(card AnswerCard) (err error)
}

// DrawAnswer removes an answer card from the deck and returns it.
//
// If there are no more cards in the deck, the discard pile is shuffled
// and the card drawn from there.
//
// If both piles are empty, an error is returned.
func (p *AnswerDeck) Draw() (card AnswerCard, err error) {
	if len(p.AnswerDeck) > 0 {
		card = p.AnswerDeck[len(p.AnswerDeck)]
		p.AnswerDeck = p.AnswerDeck[:len(p.AnswerDeck)-1]
		return
	} else if len(p.AnswerDiscard) > 0 {
		p.AnswerDeck, p.AnswerDiscard = p.AnswerDiscard[:], new([]AnswerCard)
		shuffle(p.AnswerDeck)
		card = p.AnswerDeck[len(p.AnswerDeck)]
		p.AnswerDeck = p.AnswerDeck[:len(p.AnswerDeck)-1]
		return
	}

	return nil, errors.New("answers deck and discard empty")
}

// DiscardAnswer takes the given card and puts it into the discard
// pile.
func (p *AnswerDeck) Discard(card AnswerCard) (err error) {
	p.AnswerDiscard = append(p.AnswerDiscard, card)
}

// DrawQuestion removes an answer card from the deck and returns it.
//
// If there are no more cards in the deck, the discard pile is shuffled
// and the card drawn from there.
//
// If both piles are empty, an error is returned.
func (p *PlayDeck) DrawQuestion() (card QuestionCard, err error) {
	if len(p.QuestionDeck) > 0 {
		card = p.QuestionDeck[len(p.QuestionDeck)]
		p.QuestionDeck = p.QuestionDeck[:len(p.QuestionDeck)-1]
		return
	} else if len(p.QuestionDiscard) > 0 {
		p.QuestionDeck, p.QuestionDiscard = p.QuestionDiscard[:], new([]QuestionCard)
		shuffle(p.QuestionDeck)
		card = p.QuestionDeck[len(p.QuestionDeck)]
		p.QuestionDeck = p.QuestionDeck[:len(p.QuestionDeck)-1]
		return
	}

	return nil, errors.New("answers deck and discard empty")
}

// DiscardQuestion takes the given card and puts it into the discard
// pile.
func (p *PlayDeck) DiscardQuestion(card QuestionCard) (err error) {
	p.QuestionDiscard = append(p.QuestionDiscard, card)
}

func shuffle(vals []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}