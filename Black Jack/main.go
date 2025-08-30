package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	value int
	suit  int // 0 - spades, 1 - hearts, 2 - diamonds, 3 - clubs
}

func (card Card) getString() string {
	var suit string
	var value string

	switch card.suit {
	case 0:
		suit = "♠"
	case 1:
		suit = "♥"
	case 2:
		suit = "♦"
	case 3:
		suit = "♣"
	}

	switch card.value {
	case 11:
		value = "J"
	case 12:
		value = "Q"
	case 13:
		value = "K"
	case 1:
		value = "A"
	default:
		value = fmt.Sprintf("%d", card.value)
	}

	return value + suit
}

type Deck struct {
	cards []Card
}

func (d *Deck) deal(num uint) []Card {
	if int(num) > len(d.cards) {
		num = uint(len(d.cards))
	}
	hand := d.cards[:num]
	d.cards = d.cards[num:]
	return hand
}

func (d *Deck) create() {
	d.cards = make([]Card, 0, 52)
	for suit := 0; suit < 4; suit++ {
		for value := 1; value <= 13; value++ {
			d.cards = append(d.cards, Card{value: value, suit: suit})
		}
	}
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

func getHandValue(hand []Card) int {
	value := 0
	aceCount := 0
	for _, card := range hand {
		if card.value > 10 {
			value += 10
		} else {
			value += card.value
		}
		if card.value == 1 {
			aceCount++
		}
	}
	for aceCount > 0 && value <= 11 {
		value += 10
		aceCount--
	}
	return value
}

type Game struct {
	deck        Deck
	playerCards []Card
	dealerCards []Card
}

func (game *Game) dealStartingCards() {
	game.deck.create()
	game.deck.shuffle()
	game.playerCards = game.deck.deal(2)
	game.dealerCards = game.deck.deal(2)
}

func (game *Game) play(bet float64) float64 {
	game.dealStartingCards()
	if !game.playerTurn() {
		return -bet
	}
	game.dealerTurn()
	if getHandValue(game.dealerCards) > 21 {
		fmt.Println("Dealer busts! You win.")
		return bet
	} else if getHandValue(game.dealerCards) > getHandValue(game.playerCards) {
		fmt.Println("Dealer wins.")
		return -bet
	} else if getHandValue(game.dealerCards) == getHandValue(game.playerCards) {
		fmt.Println("It's a tie!")
		return 0
	} else {
		fmt.Println("You win!")
		return bet
	}
}

func (game *Game) playerTurn() bool {
	for {
		fmt.Println("Your hand:", game.playerCards)
		fmt.Println("Dealer's hand:", game.dealerCards[0], "??")
		fmt.Print("Do you want to (h)it or (s)tand? ")
		choice := enterString()
		if choice == "h" {
			game.playerCards = append(game.playerCards, game.deck.deal(1)[0])
			if getHandValue(game.playerCards) > 21 {
				fmt.Println("You bust! Dealer wins.")
				return false
			}
		} else {
			break
		}
	}
	return true
}

func (game *Game) dealerTurn() {
	fmt.Println("Dealer's hand:", game.dealerCards)
	for getHandValue(game.dealerCards) < 17 {
		fmt.Println("Dealer hits.")
		game.dealerCards = append(game.dealerCards, game.deck.deal(1)[0])
	}
}

func enterString() string {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return ""
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	return input
}

func main() {
	balance := float64(100)

	for balance > 0 {
		fmt.Printf("Your balance is: $%.2f\n", balance)
		fmt.Printf("Enter your bet (q to quit): ")
		bet, err := strconv.ParseFloat(enterString(), 64)
		if err != nil {
			break
		}
		if bet > balance || bet <= 0 {
			fmt.Println("Invalid bet.")
			continue
		}

		game := Game{}
		balance += game.play(bet)
	}

	fmt.Printf("You left with: $%2.f\n", balance)
}
