package sol

//lint:file-ignore ST1005 Error messages are toasted, so need to be capitalized
//lint:file-ignore ST1006 Receiver name will be anything I like, thank you

import (
	"errors"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gosol/schriftbank"
)

type Foundation struct {
	parent *Pile
}

func NewFoundation(slot image.Point) *Pile {
	foundation := NewPile("Foundation", slot, FAN_NONE, MOVE_NONE)
	foundation.vtable = &Foundation{parent: &foundation}
	TheBaize.AddPile(&foundation)
	return &foundation
}

// CanAcceptTail does some obvious check on the tail before passing it to the script
func (self *Foundation) CanAcceptTail(tail []*Card) (bool, error) {
	if len(tail) > 1 {
		return false, errors.New("Cannot move more than one card to a Foundation")
	}
	if AnyCardsProne(tail) {
		return false, errors.New("Cannot add a face down card to a Foundation")
	}
	if self.Complete() {
		return false, errors.New("The Foundation is full")
	}
	return TheBaize.script.TailAppendError(self.parent, tail)
}

func (*Foundation) TailTapped([]*Card) {}

func (*Foundation) Conformant() bool {
	return true
}

// Complete - a foundation pile is complete when it contains a complete run of cards
func (self *Foundation) Complete() bool {
	return self.parent.Len() == len(CardLibrary)/len(TheBaize.script.Foundations())
}

func (*Foundation) UnsortedPairs() int {
	// you can only put a sorted sequence into a Foundation, so this will always be zero
	return 0
}

func (self *Foundation) MovableTails() []*MovableTail {
	return nil
}

func (self *Foundation) Placeholder() *ebiten.Image {
	dc := gg.NewContext(CardWidth, CardHeight)
	dc.SetColor(color.NRGBA{255, 255, 255, 31})
	dc.SetLineWidth(2)
	// draw the RoundedRect entirely INSIDE the context
	dc.DrawRoundedRectangle(1, 1, float64(CardWidth-2), float64(CardHeight-2), CardCornerRadius)
	if self.parent.label != "" {
		dc.SetFontFace(schriftbank.CardOrdinalLarge)
		dc.DrawStringAnchored(self.parent.label, float64(CardWidth)*0.5, float64(CardHeight)*0.4, 0.5, 0.5)
	}
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}
