package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 I'll call the receiver anything I like, thank you

import (
	"errors"
	"image"
)

type SimpleSimon struct {
	ScriptBase
	wikipedia string
}

func (self *SimpleSimon) BuildPiles() {

	self.stock = NewStock(image.Point{-5, -5}, FAN_NONE, 1, 4, nil, 0)

	self.discards = nil
	for x := 3; x < 7; x++ {
		d := NewDiscard(image.Point{x, 0}, FAN_NONE)
		self.discards = append(self.discards, d)
	}

	self.tableaux = nil
	for x := 0; x < 10; x++ {
		t := NewTableau(image.Point{x, 1}, FAN_DOWN, MOVE_ANY)
		self.tableaux = append(self.tableaux, t)
	}
}

func (self *SimpleSimon) StartGame() {
	// 3 piles of 8 cards each
	for i := 0; i < 3; i++ {
		pile := self.tableaux[i]
		for j := 0; j < 8; j++ {
			MoveCard(self.stock, pile)
		}
	}
	var deal int = 7
	for i := 3; i < 10; i++ {
		pile := self.tableaux[i]
		for j := 0; j < deal; j++ {
			MoveCard(self.stock, pile)
		}
		deal--
	}
	if DebugMode && self.stock.Len() > 0 {
		println("*** still", self.stock.Len(), "cards in Stock ***")
	}
}

func (*SimpleSimon) AfterMove() {
}

func (*SimpleSimon) TailMoveError(tail []*Card) (bool, error) {
	var pile *Pile = tail[0].Owner()
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (pile).category {
	case "Tableau":
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	}
	return true, nil
}

func (*SimpleSimon) TailAppendError(dst *Pile, tail []*Card) (bool, error) {
	// why the pretty asterisks? google method pointer receivers in interfaces; *Tableau is a different type to Tableau
	switch (dst).category {
	case "Discard":
		if tail[0].Ordinal() != 13 {
			return false, errors.New("Can only discard starting from a King")
		}
		for _, pair := range NewCardPairs(tail) {
			if ok, err := pair.Compare_DownSuit(); !ok {
				return false, err
			}
		}
	case "Tableau":
		if dst.Empty() {
		} else {
			return CardPair{dst.Peek(), tail[0]}.Compare_Down()
		}
	}
	return true, nil
}

func (*SimpleSimon) UnsortedPairs(pile *Pile) int {
	return UnsortedPairs(pile, CardPair.Compare_DownSuit)
}

func (*SimpleSimon) TailTapped(tail []*Card) {
	tail[0].Owner().vtable.TailTapped(tail)
}

func (*SimpleSimon) PileTapped(*Pile) {}

func (self *SimpleSimon) Wikipedia() string {
	return self.wikipedia
}

func (*SimpleSimon) CardColors() int {
	return 4
}

func (self *SimpleSimon) Complete() bool {
	return self.SpiderComplete()
}
