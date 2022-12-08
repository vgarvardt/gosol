package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized

import (
	"errors"
	"fmt"
	"sort"

	"oddstream.games/gosol/sound"
	"oddstream.games/gosol/util"
)

type ScriptBase struct {
	cells       []*Pile
	discards    []*Pile
	foundations []*Pile
	reserves    []*Pile
	stock       *Pile
	tableaux    []*Pile
	waste       *Pile
}

func (sb ScriptBase) Wikipedia() string {
	return "https://en.wikipedia.org/wiki/Patience_(game)"
}

func (sb ScriptBase) CardColors() int {
	return 2
}

func (sb ScriptBase) Cells() []*Pile {
	return sb.cells
}

func (sb ScriptBase) Foundations() []*Pile {
	return sb.foundations
}

func (sb ScriptBase) Discards() []*Pile {
	return sb.discards
}

func (sb ScriptBase) Reserves() []*Pile {
	return sb.reserves
}

func (sb ScriptBase) Stock() *Pile {
	return sb.stock
}

func (sb ScriptBase) Tableaux() []*Pile {
	return sb.tableaux
}

func (sb ScriptBase) Waste() *Pile {
	return sb.waste
}

// You can't use functions as keys in maps : the key type must be comparable
// so you can't do: var ExtendedColorMap = map[CardPairCompareFunc]bool{}
// type CardPairCompareFunc func(CardPair) (bool, error)

type Scripter interface {
	BuildPiles()
	StartGame()
	AfterMove()

	TailMoveError([]*Card) (bool, error)
	TailAppendError(*Pile, []*Card) (bool, error)
	UnsortedPairs(*Pile) int

	TailTapped([]*Card)
	PileTapped(*Pile)

	Cells() []*Pile
	Discards() []*Pile
	Foundations() []*Pile
	Reserves() []*Pile
	Stock() *Pile
	Tableaux() []*Pile
	Waste() *Pile

	Wikipedia() string
	CardColors() int
}

var Variants = map[string]Scripter{
	"Agnes Bernauer": &Agnes{
		wikipedia: "https://en.wikipedia.org/wiki/Agnes_(solitaire)",
	},
	"American Toad": &Toad{
		wikipedia: "https://en.wikipedia.org/wiki/American_Toad_(solitaire)",
	},
	"Australian": &Australian{
		wikipedia: "https://en.wikipedia.org/wiki/Australian_Patience",
	},
	"Baker's Dozen": &BakersDozen{
		wikipedia: "https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire)",
	},
	"Baker's Game": &Freecell{
		wikipedia:      "https://en.wikipedia.org/wiki/Baker%27s_Game",
		cardColors:     4,
		tabCompareFunc: CardPair.Compare_DownSuit,
	},
	"Blind Freecell": &Freecell{
		wikipedia:      "https://en.wikipedia.org/wiki/FreeCell",
		cardColors:     2,
		tabCompareFunc: CardPair.Compare_DownAltColor,
		blind:          true,
	},
	"Blockade": &Blockade{
		wikipedia: "https://en.wikipedia.org/wiki/Blockade_(solitaire)",
	},
	"Canfield": &Canfield{
		wikipedia:      "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		cardColors:     2,
		draw:           3,
		recycles:       32767,
		tabCompareFunc: CardPair.Compare_DownAltColorWrap,
	},
	"Storehouse": &Canfield{
		wikipedia:      "https://en.wikipedia.org/wiki/Canfield_(solitaire)",
		cardColors:     4,
		draw:           1,
		recycles:       2,
		tabCompareFunc: CardPair.Compare_DownSuitWrap,
		variant:        "storehouse",
	},
	"Duchess": &Duchess{
		wikipedia: "https://en.wikipedia.org/wiki/Duchess_(solitaire)",
	},
	"Klondike": &Klondike{
		wikipedia: "https://en.wikipedia.org/wiki/Solitaire",
		draw:      1,
		recycles:  2,
	},
	"Klondike Draw Three": &Klondike{
		wikipedia: "https://en.wikipedia.org/wiki/Solitaire",
		draw:      3,
		recycles:  9,
	},
	"Thoughtful": &Klondike{
		wikipedia:  "https://en.wikipedia.org/wiki/Solitaire",
		draw:       1,
		recycles:   32767,
		thoughtful: true,
	},
	"Easy": &Easy{},
	"Eight Off": &EightOff{
		wikipedia: "https://en.wikipedia.org/wiki/Eight_Off",
	},
	"Freecell": &Freecell{
		wikipedia:      "https://en.wikipedia.org/wiki/FreeCell",
		cardColors:     2,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Forty Thieves": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
	},
	"Josephine": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 4,
		moveType:    MOVE_ANY,
	},
	"Rank and File": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     2,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1, 2},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Indian": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     4,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    3,
		proneRows:      []int{0},
		tabCompareFunc: CardPair.Compare_DownOtherSuit,
	},
	"Streets": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     2,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Number Ten": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     2,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		proneRows:      []int{0, 1},
		tabCompareFunc: CardPair.Compare_DownAltColor,
		moveType:       MOVE_ANY,
	},
	"Limited": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 3,
	},
	"Forty and Eight": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab: 5,
		recycles:    1,
	},
	"Red and Black": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     2,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Lucas": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{5, 6, 7, 8, 9, 10, 11, 12},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		cardsPerTab: 3,
		dealAces:    true,
	},
	"Busy Aces": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		founds:      []int{4, 5, 6, 7, 8, 9, 10, 11},
		tabs:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		cardsPerTab: 1,
	},
	"Maria": &FortyThieves{
		wikipedia:      "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:     2,
		founds:         []int{3, 4, 5, 6, 7, 8, 9, 10},
		tabs:           []int{2, 3, 4, 5, 6, 7, 8, 9, 10},
		cardsPerTab:    4,
		tabCompareFunc: CardPair.Compare_DownAltColor,
	},
	"Sixty Thieves": &FortyThieves{
		wikipedia:   "https://en.wikipedia.org/wiki/Forty_Thieves_(solitaire)",
		cardColors:  4,
		packs:       3,
		founds:      []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		tabs:        []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		cardsPerTab: 5,
	},
	"Penguin": &Penguin{
		wikipedia: "https://www.parlettgames.uk/patience/penguin.html",
	},
	"Scorpion": &Scorpion{
		wikipedia: "https://en.wikipedia.org/wiki/Scorpion_(solitaire)",
	},
	"Simple Simon": &SimpleSimon{
		wikipedia: "https://en.wikipedia.org/wiki/Simple_Simon_(solitaire)",
	},
	"Spider One Suit": &Spider{
		wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		cardColors: 1,
		packs:      8,
		suits:      1,
	},
	"Spider Two Suits": &Spider{
		wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		cardColors: 2,
		packs:      4,
		suits:      2,
	},
	"Spider Four Suits": &Spider{
		wikipedia:  "https://en.wikipedia.org/wiki/Spider_(solitaire)",
		cardColors: 4,
		packs:      2,
		suits:      4,
	},
	"Classic Westcliff": &Westcliff{
		variant: "Classic",
	},
	"American Westcliff": &Westcliff{
		variant: "American",
	},
	"Easthaven": &Westcliff{
		variant: "Easthaven",
	},
	"Whitehead": &Whitehead{
		wikipedia: "https://en.wikipedia.org/wiki/Klondike_(solitaire)",
	},
	"Usk": &Usk{
		wikipedia: "https://politaire.com/help/usk",
	},
	"Yukon": &Yukon{
		wikipedia: "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
	},
	"Yukon Cells": &Yukon{
		wikipedia:  "https://en.wikipedia.org/wiki/Yukon_(solitaire)",
		extraCells: 2,
	},
}

var VariantGroups = map[string][]string{
	// "All" added dynamically by func init()
	// don't have Agnes here (as a group) because it would come before All
	// and Agnes Sorel is retired because it's just too hard
	"> Canfield":      {"Canfield", "Storehouse", "Duchess", "American Toad"},
	"> Easier":        {"American Toad", "American Westcliff", "Blockade", "Classic Westcliff", "Lucas", "Spider One Suit"},
	"> Harder":        {"Baker's Dozen", "Easthaven", "Forty Thieves", "Spider Four Suits", "Usk"},
	"> Forty Thieves": {"Forty Thieves", "Number Ten", "Red and Black", "Indian", "Rank and File", "Sixty Thieves", "Josephine", "Limited", "Forty and Eight", "Lucas", "Busy Aces", "Maria", "Streets"},
	"> Freecell":      {"Baker's Game", "Blind Freecell", "Freecell", "Eight Off"},
	"> Klondike":      {"Klondike", "Klondike Draw Three", "Thoughtful", "Whitehead"},
	"> People":        {"Agnes Bernauer", "Duchess", "Josephine", "Maria", "Simple Simon", "Baker's Game"},
	"> Places":        {"Australian", "Yukon", "Klondike", "Usk"},
	"> Puzzlers":      {"Penguin", "Simple Simon", "Baker's Dozen", "Freecell"},
	"> Spider":        {"Spider One Suit", "Spider Two Suits", "Spider Four Suits", "Scorpion"},
	"> Yukon":         {"Yukon", "Yukon Cells"},
}

// init is used to assemble the "> All" alpha-sorted group of variants for the picker menu
func init() {
	var vnames []string
	for k := range Variants {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	VariantGroups["> All"] = vnames
}

// VariantGroupNames returns an alpha-sorted []string of the variant group names
func VariantGroupNames() []string {
	var vnames []string = make([]string, 0, len(VariantGroups))
	for k := range VariantGroups {
		vnames = append(vnames, k)
	}
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

// VariantNames returns an alpha-sorted []string of the varaints in a group
func VariantNames(group string) []string {
	var vnames []string = make([]string, 0, len(VariantGroups[group]))
	vnames = append(vnames, VariantGroups[group]...)
	sort.Slice(vnames, func(i, j int) bool { return vnames[i] < vnames[j] })
	return vnames
}

// useful generic game library of functions

func Compare_Empty(p *Pile, c *Card) (bool, error) {

	if p.Label() != "" {
		if p.Label() == "x" {
			return false, errors.New("Cannot move cards there")
		}
		ord := util.OrdinalToShortString(c.Ordinal())
		if ord != p.Label() {
			return false, fmt.Errorf("Can only accept %s, not %s", util.ShortOrdinalToLongOrdinal(p.Label()), util.ShortOrdinalToLongOrdinal(ord))
		}
	}
	return true, nil
}

func RecycleWasteToStock(waste *Pile, stock *Pile) {
	if TheBaize.Recycles() > 0 {
		for waste.Len() > 0 {
			MoveCard(waste, stock)
		}
		TheBaize.SetRecycles(TheBaize.Recycles() - 1)
		switch {
		case TheBaize.recycles == 0:
			sound.Play("Error")
			TheUI.Toast("No more recycles")
		case TheBaize.recycles == 1:
			sound.Play("Bong")
			TheUI.Toast(fmt.Sprintf("%d recycle remaining", TheBaize.Recycles()))
		case TheBaize.recycles < 10:
			sound.Play("Bong")
			TheUI.Toast(fmt.Sprintf("%d recycles remaining", TheBaize.Recycles()))
		}
	} else {
		sound.Play("Error")
		TheUI.Toast("No more recycles")
	}
}

func UnsortedPairs(pile *Pile, fn func(CardPair) (bool, error)) int {
	if pile.Len() < 2 {
		return 0
	}
	var unsorted int
	for _, pair := range NewCardPairs(pile.cards) {
		if pair.EitherProne() {
			unsorted++
		} else {
			if ok, _ := fn(pair); !ok {
				unsorted++
			}
		}
	}
	return unsorted
}

type CardPair struct {
	c1, c2 *Card
}

func (cp CardPair) EitherProne() bool {
	return cp.c1.Prone() || cp.c2.Prone()
}

type CardPairs []CardPair

func NewCardPairs(cards []*Card) []CardPair {
	if len(cards) < 2 {
		return []CardPair{}
	}
	var cpairs []CardPair
	c1 := cards[0]
	for i := 1; i < len(cards); i++ {
		c2 := cards[i]
		cpairs = append(cpairs, CardPair{c1, c2})
		c1 = c2
	}
	return cpairs
}

func (cpairs CardPairs) Print() {
	for _, pair := range cpairs {
		println(pair.c1.String(), pair.c2.String())
	}
}

func (cp CardPair) Compare_Up() (bool, error) {
	if cp.c1.Ordinal()+1 != cp.c2.Ordinal() {
		return false, errors.New("Cards must be in ascending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_Down() (bool, error) {
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence")
	}
	return true, nil
}

func (cp CardPair) Compare_DownColor() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownAltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownColorWrap() (bool, error) {
	if cp.c1.Black() != cp.c2.Black() {
		return false, errors.New("Cards must be the same color")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
	}
	return true, nil
}

func (cp CardPair) Compare_DownAltColorWrap() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal() != cp.c2.Ordinal()+1 {
		return false, errors.New("Cards must be in descending sequence (Kings on Aces allowed)")
	}
	return true, nil
}

func (cp CardPair) Compare_UpAltColor() (bool, error) {
	if cp.c1.Black() == cp.c2.Black() {
		return false, errors.New("Cards must be in alternating colors")
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_UpSuit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return cp.Compare_Up()
}

func (cp CardPair) Compare_DownSuit() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_DownOtherSuit() (bool, error) {
	if cp.c1.Suit() == cp.c2.Suit() {
		return false, errors.New("Cards must not be the same suit")
	}
	return cp.Compare_Down()
}

func (cp CardPair) Compare_UpSuitWrap() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if cp.c1.Ordinal() == 13 && cp.c2.Ordinal() == 1 {
		return true, nil // Ace on King
	}
	if cp.c1.Ordinal() == cp.c2.Ordinal()-1 {
		return true, nil
	}
	return false, errors.New("Cards must go up in rank (Aces on Kings allowed)")
}

func (cp CardPair) Compare_DownSuitWrap() (bool, error) {
	if cp.c1.Suit() != cp.c2.Suit() {
		return false, errors.New("Cards must be the same suit")
	}
	if cp.c1.Ordinal() == 1 && cp.c2.Ordinal() == 13 {
		return true, nil // King on Ace
	}
	if cp.c1.Ordinal()-1 == cp.c2.Ordinal() {
		return true, nil
	}
	return false, errors.New("Cards must go down in rank (Kings on Aces allowed)")
}
