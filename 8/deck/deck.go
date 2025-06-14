package deck

import (
	"math/rand"
	"sort"
	"time"
)

type Card struct {
	Rank int
	Suit string
}

const (
	Spades   = "spades"
	Hearts   = "hearts"
	Diamonds = "diamonds"
	Clubs    = "clubs"
	Joker    = "joker"
)

const (
	Ace = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var SuitOrder = map[string]int{
	Spades:   0,
	Hearts:   1,
	Diamonds: 2,
	Clubs:    3,
	Joker:    4,
}

const suitLength = 13

func (c Card) String() string {
	if c.Suit == Joker {
		return "Joker"
	}
	RankNames := []string{
		"Ace", "Two", "Three", "Four", "Five", "Six", "Seven",
		"Eight", "Nine", "Ten", "Jack", "Queen", "King",
	}
	return RankNames[c.Rank] + " of " + c.Suit
}

var Suits = []string{Spades, Hearts, Diamonds, Clubs}

type Option func([]Card) []Card

func New(options ...Option) []Card {
	var cards []Card

	for i := 0; i < 52; i++ {
		cards = append(cards, Card{
			Rank: i % suitLength,
			Suit: Suits[i/suitLength],
		})
	}

	for _, option := range options {
		cards = option(cards)
	}

	return cards
}

func WithJokers(count int) Option {
	return func(cards []Card) []Card {
		for i := 0; i < count; i++ {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: -1,
			})
		}
		return cards
	}
}

func WithDecks(n int) Option {
	return func(cards []Card) []Card {
		var newCards []Card
		for i := 0; i < n; i++ {
			newCards = append(newCards, cards...)
		}
		return newCards
	}
}

func Filter(predicate func(Card) bool) Option {
	return func(cards []Card) []Card {
		var filtered []Card
		for _, c := range cards {
			if !predicate(c) {
				filtered = append(filtered, c)
			}
		}
		return filtered
	}
}

func Shuffle(cards []Card) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func Sort(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Suit != cards[j].Suit {
			return SuitOrder[cards[i].Suit] < SuitOrder[cards[j].Suit]
		}
		return cards[i].Rank < cards[j].Rank
	})
}

func SortWith(cards []Card, less func(i, j int) bool) {
	sort.Slice(cards, less)
}
